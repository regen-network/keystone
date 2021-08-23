package keys

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/frumioj/crypto11"
)

const (
	KEYGEN_SECP256K1 KeygenAlgorithm = iota
	KEYGEN_SECP256R1
	KEYGEN_ED25519
	KEYGEN_RSA
)

type KeygenAlgorithm int

type CryptoKey struct {
	Label      string
	Algo       KeygenAlgorithm
	pubKey     PubKey
	signer     crypto11.Signer
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

// PubKey is exactly the same as the cosmos-sdk version
// except without the proto.Message dependency
type PubKey interface {
	Address() tmcrypto.Address
	Bytes() []byte
	VerifySignature(msg []byte, sig []byte) bool
	Equals(PubKey) bool
	Type() string
}

func (pk CryptoKey) Bytes() []byte{
	return []byte{}
}

func (pk CryptoKey) Sign( msg []byte) ([]byte, error){
	return []byte{}, nil
}

func (pk CryptoKey) Equals(other CryptoKey) bool {
	return true
}

func (pk CryptoKey) PubKey() PubKey {
	return pk.pubKey
}

func (pk CryptoKey) Type() string { return "CryptoKey" }
