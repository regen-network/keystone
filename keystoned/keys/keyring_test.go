package keys

import (
	"log"
	"testing"
	"crypto/x509"
	"encoding/pem"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	btcsecp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

func TestCreateKeySecp256k1(t *testing.T) {

	// hardcoded path for now - might change this API to actually
	// take a JSON string and let caller decide how to get that
	// JSON
	kr, err := NewPkcs11FromConfig("./pkcs11-config")
	require.NoError(t, err)
	label, err := randomBytes(16)
	require.NoError(t, err)
	
	key, err := kr.NewKey( KEYGEN_SECP256K1, string(label) )
	require.NoError(t, err)
	require.NotNil(t, key)

	msg := []byte("Signing this plaintext tells me what exactly?")
	signed, err := key.Sign(msg, nil)

	require.NoError(t, err)
	log.Printf("Signed byes: %v", signed)
	pubKey := getPubKey( &key )
	log.Printf("Key: %v", pubKey.(*secp256k1.PubKey))
	secp256k1key := pubKey.(*secp256k1.PubKey)
	pub, err := btcsecp256k1.ParsePubKey(secp256k1key.Key, btcsecp256k1.S256())
	
	if err != nil {
		log.Printf("Not a secp256k1 key?")
	}

	log.Printf("Pub: %v", pub) 

	// @@TODO fails because sign doesn't sign in the right way
	// including the possibility of malleable sig
	// need to unDER the signature, and make sure it's low-s
	// normalized
	valid := secp256k1key.VerifySignature( msg, signed )
	
	log.Printf("Did the signature verify? True = yes: %v", valid)
	
	err = key.Delete()
	require.NoError(t, err)
}

func TestCreateKeySecp256r1(t *testing.T) {

	// hardcoded path for now - might change this API to actually take a JSON string and let caller decide how to get that JSON
	kr, err := NewPkcs11FromConfig("./pkcs11-config")
	require.NoError(t, err)
	label, err := randomBytes(16)
	require.NoError(t, err)
	
	key, err := kr.NewKey( KEYGEN_SECP256R1, string(label) )
	require.NoError(t, err)
	require.NotNil(t, key)

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	log.Printf("Public: %s", pemEncodedPub)
	//pub := key.PubKey()
	
	//log.Printf("Address: %s", string(pub.Address()))

	key2 := key

	//log.Printf("Keys should be equal: %v", key.Equals(key2))

	//key3, err := kr.NewKey( KEYGEN_SECP256K1, string(label) )
	//require.NoError(t, err)
	//log.Printf("Keys should NOT be equal: %v", key.Equals(key3))
	
	err = key.Delete()
	require.NoError(t, err)

	// This delete should fail since key2 is a pointer to key
	// which was already deleted
	err = key2.Delete()

	// Yes, there SHOULD be an error on this delete!
	require.Error(t, err)

	// key3 delete should pass
	//err = key3.Delete()
	//require.NoError(t, err)
	
}
