// The keys package is collection of structures and functions used to
// manipulate a "keyring" -- a collection of keys, and the keys
// themselves. A keyring could consist of an OS directory containing
// keys stored in files, or other implementations could be
// imagined. In this particular case though, a concrete implementation
// is given of a PKCS11-based keyring, which is a set of keys stored
// on a cryptographic token, such as an HSM, which offers the PKCS11
// API to its keys.

package keys

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"log"
	"os"

	"github.com/frumioj/crypto11"
)

type Pkcs11Keyring struct {
	ModulePath string
	TokenLabel string
	pin        string
	ctx        *crypto11.Context
}

// Keyring interface provides the methods for keyring
// implementations.
// NewKey generates a new key, using the given keygen algorithm
// supported algorithms are in keys.go
// Key returns a filled out key with the given label, retrieved from the
// keyring
// ListKeys lists all of the keys on the keyring
type Keyring interface {
	NewKey(algorithm KeygenAlgorithm, label string) (*CryptoKey, error)
	Key(label string) (*CryptoKey, error)
	// @@TODO - not implemented for PKCS11 keyring 9/9/2021
	//ListKeys() ([]CryptoKey, error)
}

// NewKey creates a new ECC key on a Pkcs11 token
// using the given algorithm from the keygen algos supported. A label
// can be passed in. This is used as a way of uniquely identifying the key
// and typically is a large (unguessable) random number
func (ring Pkcs11Keyring) NewKey(algorithm KeygenAlgorithm, label string) (*CryptoKey, error) {

	// Crypto-secure random bytes
	id, err := CryptoRandomBytes(16)

	if err != nil {
		log.Printf("Error making key ID: %s", err.Error())
		return nil, err
	}

	var key crypto11.Signer

	switch algorithm {
	case KEYGEN_SECP256K1:
		key, err = ring.ctx.GenerateECDSAKeyPairWithLabel(id, []byte(label), crypto11.P256K1())
	case KEYGEN_SECP256R1:
		key, err = ring.ctx.GenerateECDSAKeyPairWithLabel(id, []byte(label), elliptic.P256())
	default:
		return nil, err
	}

	if err != nil {
		log.Printf("Error generating key: %s", err.Error())
		return nil, err
	} else {
		log.Printf("Key made: %v", key)
	}

	newkey := CryptoKey{Label: label, Algo: algorithm, signer: key}
	pubkey := getPubKey(&newkey)
	newkey.pubk = pubkey
	
	return &newkey, nil
}

// Key retrieves a keypair from the PKCS11 token and populates a
// CryptoKey object, based on finding the key[air based on the label
// that is supplied in the API call.
func (ring Pkcs11Keyring) Key(label string) (*CryptoKey, error) {
	
	// Note: this API retrieves key PAIRS, so only asymmetric key
	// algorithms
	keys, err := ring.ctx.FindKeyPairs(nil, []byte(label))

	if err != nil {
		log.Printf("Key could not be found, with error: %s", err.Error())
		return nil, err
	}

	// @@TODO fill out the Algo by retrieving the key type and
	// thus the curve name - requires some testing though
	// to determine exactly how
	newkey := CryptoKey{Label: label, signer: keys[0]}
	pubkey := getPubKey(&newkey)
	newkey.pubk = pubkey

	return &newkey, nil
}

// NewPkcs11FromConfig returns a new Pkcs11Keyring structure when
// given the path to a configuration file that describes the Pkcs11
// token which holds the actual cryptographic keys.
func NewPkcs11FromConfig(configPath string) (*Pkcs11Keyring, error) {

	kr := Pkcs11Keyring{}
	cfg, err := getConfig(configPath)

	if err != nil {
		log.Printf("Could not create new Pkcs11 keyring: %s", err.Error())
		return nil, err
	}

	kr.ctx, err = crypto11.Configure(cfg)

	if err != nil {
		log.Printf("Slot configuration failed with %s", err.Error())
		return nil, err
	}

	kr.ModulePath = cfg.Path
	kr.TokenLabel = cfg.TokenLabel

	return &kr, nil
}

// getConfig returns a crypto11 Config struct representing the Pkcs11
// token, when given the location of a JSON configuration file.
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

// CryptoRandomBytes returns n bytes obtained from a local source of
// crypto-secure randomness. This can be used for generating
// hard-to-guess key labels, for example.
func CryptoRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		log.Printf("Error reading random bytes: %s", err.Error())
		return nil, err
	}

	return b, nil
}
