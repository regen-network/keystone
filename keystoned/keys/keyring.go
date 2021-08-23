package keys

import (
	"os"
	"log"
	"encoding/json"
	"crypto/rand"
	"crypto/elliptic"
	
	"github.com/frumioj/crypto11"
)

type Pkcs11Keyring struct {
	ModulePath string
	TokenLabel string
	pin string
	ctx *crypto11.Context
}

// Keyring interface provides the methods for keyring
// implementations.
// NewKey generates a new key, using the given keygen algorithm
// supported algorithms are in keys.go
// Key returns a filled out key with the given label, retrieved from the
// keyring
// ListKeys lists all of the keys on the keyring
type Keyring interface {
	NewKey( algorithm KeygenAlgorithm, label string ) (CryptoKey, error)
	Key( label string ) (CryptoKey, error)
	ListKeys() ([]CryptoKey, error)
}

// NewKey creates a new ECC key on a Pkcs11 token
// using the given algorithm from the keygen algos supported. A label
// can be passed in. This is used as a way of uniquely identifying the key
// and typically is a large (unguessable) random number
func (ring Pkcs11Keyring) NewKey(algorithm KeygenAlgorithm, label string) (CryptoKey, error) {

	id, err := randomBytes(16)

	if err != nil {
		log.Printf("Error making key ID: %s", err.Error())
		return CryptoKey{}, err
	}

	var key crypto11.Signer
	
	switch algorithm {
	case KEYGEN_SECP256K1:
		key, err = ring.ctx.GenerateECDSAKeyPairWithLabel(id, []byte(label), crypto11.P256K1())
	case KEYGEN_SECP256R1:
		key, err = ring.ctx.GenerateECDSAKeyPairWithLabel(id, []byte(label), elliptic.P256())
	case KEYGEN_RSA:
		key, err = ring.ctx.GenerateRSAKeyPairWithLabel(id, []byte(label), 2048)
	default:
		return CryptoKey{}, err
	}

	if err != nil {
		log.Printf("Error generating key: %s", err.Error())
		return CryptoKey{}, err
	} else {
		log.Printf("Key made: %v", key)
	}

	
	return CryptoKey{Label: label, Algo: algorithm, signer: key}, nil
}

func NewPkcs11FromConfig(configPath string) (Pkcs11Keyring, error) {

	kr := Pkcs11Keyring{}
	cfg, err := getConfig( configPath )
	
	if err != nil {
		log.Printf("Could not create new Pkcs11 keyring: %s", err.Error())
		return Pkcs11Keyring{}, err
	}

	kr.ctx, err = crypto11.Configure( cfg )

	if err != nil {
		log.Printf("Slot configuration failed with %s", err.Error())
		return Pkcs11Keyring{}, err
	}
	
	kr.ModulePath = cfg.Path
	kr.TokenLabel = cfg.TokenLabel

	return kr, nil
}

func getConfig(configLocation string) (ctx *crypto11.Config, err error) {
	file, err := os.Open(configLocation)
	
	if err != nil {
		log.Printf("Could not open config file: %s", configLocation)
		return nil, err
	}
	
	defer func() {
		err = file.Close()
	}()

	configDecoder := json.NewDecoder(file)
	config := &crypto11.Config{}
	err = configDecoder.Decode(config)
	
	if err != nil {
		log.Printf("Could not decode config file: %s", err.Error())
		return nil, err
	}
	
	return config, nil
}

func randomBytes( n int ) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		log.Printf("Error reading random bytes: %s", err.Error())
		return nil, err
	}

	return b, nil
}
