package keys

import (
	"crypto"
	"crypto/sha256"
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
type CryptoPubKey struct {
	Key []byte
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

func (pk CryptoKey) Sign(msg []byte) ([]byte, error) {
	return []byte{}, nil
}

func (pk CryptoKey) Equals(other CryptoKey) bool {
	return true
}

func (pk CryptoKey) PubKey() PubKey {
	return pk.pubKey
}

func (pk CryptoKey) Type() string { return "CryptoKey" }

func (pk CryptoKey) Public() crypto.PublicKey { return pk.signer.Public() }

func (pub *CryptoPubKey) Address() tmcrypto.Address {
	if len(pub.Key) != PUBLIC_KEY_SIZE {
		panic("length of pubkey is incorrect")
	}

	sha := sha256.Sum256(pub.Key)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha[:]) // does not error
	return tmcrypto.Address(hasherRIPEMD160.Sum(nil))
}
