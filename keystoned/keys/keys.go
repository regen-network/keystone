package keys

import (
	"log"
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/rand"
	"golang.org/x/crypto/ripemd160"

	"github.com/frumioj/crypto11"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

const (
	KEYGEN_SECP256K1 KeygenAlgorithm = iota
	KEYGEN_SECP256R1
	KEYGEN_ED25519
	KEYGEN_RSA
)

const PUBLIC_KEY_SIZE = 33

type KeygenAlgorithm int

type CryptoKey struct {
	Label  string
	Algo   KeygenAlgorithm
	signer crypto11.Signer
}

// CryptoPrivKey looks exactly the same as the LedgerPrivKey
// interface from cosmos-sdk/crypto/types. There is no ability
// to retrieve the private key bytes because these are stored
// and used only within the HSM.
type CryptoPrivKey interface {
	Bytes() []byte
	Sign(msg []byte) ([]byte, error)
	PubKey() PubKey
	Equals(CryptoPrivKey) bool
	Type() string
}

// CryptoPubKey looks a lot like a tmcrypto-inherited
// PubKey, but is not defined in a protobuf message
// See cosmos-sdk/crypto/keys/internal/ecdsa/pubkey.go for inspiration
type CryptoPubKey struct {
	crypto.PublicKey

	address tmcrypto.Address
}

// PubKey is exactly the same as the cosmos-sdk version
// except without the proto.Message dependency
type PubKey interface {
	Address() tmcrypto.Address
	Bytes() []byte
	VerifySignature(msg []byte, sig []byte) bool
	Equals(PubKey) bool
	Type() string
}

// Bytes will return only an empty byte array
// because this key does not have access to
// the actual key bytes
func (pk CryptoKey) Bytes() []byte {
	return []byte{}
}

// Sign a plaintext with this private key, first hashing the
// plaintext with the "optional" hash function (pass nil for
// a 'null' hash).
func (pk CryptoKey) Sign(plaintext []byte, hashFun *crypto.Hash ) ([]byte, error) {

	if hashFun != nil {

		h := hashFun.New()
		_, err := h.Write(plaintext)
		digest := h.Sum(nil)
		
		if err != nil {
			log.Printf("Error creating digest %s", err.Error())
			return []byte{}, nil
		}
		
		return pk.signer.Sign(rand.Reader, digest[:], nil)
	} else {
		return pk.signer.Sign(rand.Reader, plaintext, nil)
	}
}

// Equals checks whether two CryptoKeys are equal -
// because there are no pk bytes to compare, this
// comparison is done by signing a plaintext with both
// keys, and if the signature bytes are equal, then
// the keys are considered equal.
func (pk CryptoKey) Equals(other CryptoKey) bool {

	plaintext := "are these keys equal?"

	this, err := pk.Sign([]byte(plaintext), nil)

	if err != nil {
		// should this actually return an error though
		// if signing fails?
		return false
	}

	that, err := other.Sign([]byte(plaintext), nil)

	if err != nil {
		return false
	}

	if bytes.Equal( this, that ) {
		return true
	} else {
		return false
	}
	
	
}

func (pk CryptoKey) PubKey() *CryptoPubKey { return &CryptoPubKey{pk.signer.Public(), nil }}

func (pk CryptoKey) Type() string { return "CryptoKey" }

func (pk CryptoKey) Delete() error { return pk.signer.Delete() }

func (pk CryptoKey) Public() crypto.PublicKey { return pk.signer.Public() }

func (pubKey *CryptoPubKey) Bytes() []byte {
	switch pub := pubKey.PublicKey.(type) {
	case *ecdsa.PublicKey:
		return elliptic.MarshalCompressed(pub.Curve, pub.X, pub.Y)
	default:
		panic("Unsupported public key type!")
	}
}

// Address takes a CryptoPubKey, expecting that it has
// a crypto.PublicKey base struct, marshals the struct into bytes using
// ANSI X.
func (pubKey *CryptoPubKey) Address() tmcrypto.Address {

	if pubKey.address == nil {
		switch pub := pubKey.PublicKey.(type) {
		case *ecdsa.PublicKey:
			// @@ TODO: currently does the btc secp256k1 transform
			// but should also support r1, by looking first at
			// curve params - switch inside a switch
			publicKeyBytes := pubKey.Bytes()
			sha := sha256.Sum256(publicKeyBytes)
			hasherRIPEMD160 := ripemd160.New()
			hasherRIPEMD160.Write(sha[:]) // does not error
			pubKey.address = tmcrypto.Address(hasherRIPEMD160.Sum(nil))
			return pubKey.address
		default:
			log.Printf("Type: %T", pub)
			panic("Unsupported public key!")
		}
	} else {
		return pubKey.address
	}

}
