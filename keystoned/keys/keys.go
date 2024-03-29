package keys

import (
	"bytes"
	"errors"
	"log"
	"math/big"

	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/frumioj/crypto11"
)

const (
	KEYGEN_SECP256K1 KeygenAlgorithm = iota
	KEYGEN_SECP256R1
	KEYGEN_ED25519
)

// SigningProfile is a combination of cryptographic signing mechanism,
// prior hashing of the plaintext, transformations such as
// s-normalization of signature, and post-encoding of the signature.
const (
	// SIGNING_OPTS_BC_ECDSA_SHAXXX means
	//   i) SHAXXX hash prior to signing
	//  ii) Raw signature (R||S, no DER)
	// iii) low-s normalized
	SIGNING_OPTS_BC_ECDSA_SHA256 SigningProfile = iota

	// Could also have:
	//SIGNING_OPTS_BC_ECDSA_SHA384
	//SIGNING_OPTS_BC_ECDSA_SHA512
	// just need to add the appropriate hashing into the Sign API

	// SIGNING_OPTS_ECDSA means
	//   i) No hash in the signing process
	//  ii) DER signature as in usual ECDSA
	// iii) No low-s normalization
	SIGNING_OPTS_ECDSA
)

const PUBLIC_KEY_SIZE = 33

type KeygenAlgorithm int
type SigningProfile int

type CryptoKey struct {
	Label  string
	Algo   KeygenAlgorithm
	signer crypto11.Signer
	pubk   types.PubKey
}

// CryptoPrivKey looks almost exactly the same as the LedgerPrivKey
// interface from cosmos-sdk/crypto/types. There is no ability to
// retrieve the private key bytes because these are stored and used
// only within the HSM.
type CryptoPrivKey interface {
	Bytes() []byte
	Sign(msg []byte, opts *SigningProfile) ([]byte, error)
	PubKey() types.PubKey
	Equals(CryptoPrivKey) bool
	Type() string
}

// Bytes will return only an empty byte array
// because this key does not have access to
// the actual key bytes
func (pk *CryptoKey) Bytes() []byte {
	return []byte{}
}

// Sign a plaintext with this private key. The SigningProfile tells
// the function which way of pre- and post-encoding of the actual
// cryptographic signature, which includes prior hashing, and whether
// or not the signature should be DER-encoded or "raw" (two
// concatenated big Ints)
func (pk *CryptoKey) Sign(plaintext []byte, opts *SigningProfile) ([]byte, error) {

	var digested []byte

	// Blockchain-flavoured ECDSA (as of 9/2021) means required
	// sha256 hashing of plaintext prior to signing.
	if opts == nil || *opts == SIGNING_OPTS_BC_ECDSA_SHA256 {
		digest := sha256.Sum256(plaintext)
		digested = digest[:]
	} else {
		digested = plaintext
	}

	sigbytes, err := pk.signer.Sign(rand.Reader, digested, nil)

	if err != nil {
		log.Printf("Signature failed: %s", err.Error())
		return nil, err
	}

	// Default signature mechanism is the blockchain flavour of
	// ECDSA (see const definitions above) which means now getting
	// the raw signature and low-s normalizing the s component of
	// the signature
	if opts == nil || *opts == SIGNING_OPTS_BC_ECDSA_SHA256 {
		// un-DER the sig
		var rawsig *dsaSignature
		rawsig, err := unmarshalDER(sigbytes)

		if err != nil {
			log.Printf("Error getting ints from DER: %s", err.Error())
			return nil, err
		}

		return signatureRaw(rawsig.R, NormalizeS(rawsig.S, crypto11.P256K1())), nil

	} else {
		return sigbytes, nil
	}
}

// Equals checks whether two CryptoKeys are equal -
// because there are no pk bytes to compare, this
// comparison is done by signing a plaintext with both
// keys, and if the signature bytes are equal, then
// the keys are considered equal.
func (pk *CryptoKey) Equals(other CryptoKey) bool {

	this := pk.PubKey()
	that := other.PubKey()

	return bytes.Equal(this.Bytes(), that.Bytes())
}

func (pk *CryptoKey) PubKey() types.PubKey { return pk.pubk }

func (pk *CryptoKey) Type() string { return "CryptoKey" }

func (pk *CryptoKey) KeyType() KeygenAlgorithm { return pk.Algo }

func (pk *CryptoKey) Delete() error { return pk.signer.Delete() }

func (pk *CryptoKey) Public() crypto.PublicKey { return pk.signer.Public() }

func (pk *CryptoKey) PubKeyBytes() []byte {
	switch pub := pk.Public().(type) {
	case *ecdsa.PublicKey:
		// is this OK for a *btcec* secp256k1 key?
		return elliptic.MarshalCompressed(pub.Curve, pub.X, pub.Y)
	default:
		panic("Unsupported public key type!")
	}
}

// Address takes a PubKey, expecting that it has
// a crypto.PublicKey base struct, marshals the struct into bytes using
// ANSI X.
// func (pubKey *PubKey) Address() tmcrypto.Address {

// 	if pubKey.address == nil {
// 		switch pub := pubKey.PublicKey.(type) {
// 		case *ecdsa.PublicKey:
// 			// @@ TODO: currently does the btc secp256k1 transform
// 			// but should also support r1, by looking first at
// 			// curve params - switch inside a switch
// 			publicKeyBytes := pubKey.Bytes()
// 			sha := sha256.Sum256(publicKeyBytes)
// 			hasherRIPEMD160 := ripemd160.New()
// 			hasherRIPEMD160.Write(sha[:]) // does not error
// 			pubKey.address = tmcrypto.Address(hasherRIPEMD160.Sum(nil))
// 			return pubKey.address
// 		default:
// 			log.Printf("Type: %T", pub)
// 			panic("Unsupported public key!")
// 		}
// 	} else {
// 		return pubKey.address
// 	}

// }

// Equals checks whether two PubKeys are equal -
// by checking their marshalled byte values
// func (pubk *PubKey) Equals(other cryptotypes.PubKey) bool {

// 	this := pubk.Bytes()
// 	that := other.Bytes()

// 	return bytes.Equal(this, that)
// }

// dsaSignature contains the two integers needed for
// an ECDSA signature value. They must be put in a struct
// to allow the asn1 unmarshalling which uses an interface{}
// type to return the values, instead of just returning the
// two integers.
type dsaSignature struct {
	R, S *big.Int
}

// unmarshalDER takes a DER-encoded byte array, and dumps
// it into a (hopefully-appropriate) struct. If the struct
// given, is not appropriate for the data, then unmarshalling
// will fail.
func unmarshalDER(sigDER []byte) (*dsaSignature, error) {
	var sig dsaSignature

	if rest, err := asn1.Unmarshal(sigDER, &sig); err != nil {
		return nil, err
	} else if len(rest) > 0 {
		return nil, errors.New("unexpected data found after DSA signature")
	}

	return &sig, nil
}

// isSNormalized returns true for the integer sigS if sigS falls in
// lower half of the curve order
// It is expected that the caller passes the curve order as a big Int along
// with the s portion of the signature.
func isSNormalized(sigS *big.Int, order *big.Int) bool {
	// return the result of comparing the given s signature
	// component with half the value of the curve order. If the s
	// component is less than or equal to half the curve order,
	// then returns true (!= 1), if > than, will return false
	// (==1)
	return sigS.Cmp(new(big.Int).Rsh(order, 1)) != 1
}

// NormalizeS will invert the s value if not already in the lower half
// of curve order value by subtracting it from the curve order (N)
func NormalizeS(sigS *big.Int, curve elliptic.Curve) *big.Int {
	if isSNormalized(sigS, curve.Params().N) {
		return sigS
	} else {
		order := curve.Params().N
		return new(big.Int).Sub(order, sigS)
	}
}

// signatureRaw takes two big integers and returns a byte value that
// is the result of concatenating the byte values of each of the given
// integers. The byte values are left-padded with zeroes
func signatureRaw(r *big.Int, s *big.Int) []byte {

	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}

// VerifySignature takes a plaintext msg and the signed plaintext
// which should be a DER-encoded byte array which can be marshalled
// into two big Ints - r and s, which represent an ECDSA signature.
// @@TODO: what if the signature is EDDSA or some other non-ECDSA
// option that doesn't marshal to r and s?
// func (pubk *PubKey) VerifySignature(msg []byte, sig []byte) bool {

// 	rawsig, err := unmarshalDER(sig)

// 	if err != nil {
// 		log.Printf("Signature verification failed DER decode with: %s", err.Error())
// 		return false
// 	}

// 	return ecdsa.Verify(pubk.PublicKey.(*ecdsa.PublicKey), msg, rawsig.R, rawsig.S)
// }

func getPubKey(pk *CryptoKey) types.PubKey {
	switch pub := pk.Public().(type) {
	case *ecdsa.PublicKey:
		log.Printf("Curve: %s", pub.Curve.Params().Name)

		log.Printf("Curve = p256? %v", pub.Curve == elliptic.P256())
		// is this OK for a *btcec* secp256k1 key?
		return &secp256k1.PubKey{Key: elliptic.MarshalCompressed(pub.Curve, pub.X, pub.Y)}
	default:
		panic("Unsupported public key type!")
	}
}
