package keys

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"testing"

	//"github.com/stretchr/testify/assert"
	btcsecp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCreateKeySecp256k1(t *testing.T) {

	// hardcoded path for now - might change this API to actually
	// take a JSON string and let caller decide how to get that
	// JSON
	kr, err := NewPkcs11FromConfig("./pkcs11-config")
	require.NoError(t, err)
	label, err := randomBytes(16)
	require.NoError(t, err)

	key, err := kr.NewKey(KEYGEN_SECP256K1, string(label))
	require.NoError(t, err)
	require.NotNil(t, key)

	msg := []byte("Signing this plaintext tells me what exactly?")
	signed, err := key.Sign(msg, nil)

	require.NoError(t, err)
	log.Printf("Signed byes: %v", signed)

	pubkey := key.PubKey()

	log.Printf("Pubkey: %v", pubkey)

	require.Equal(t, key.KeyType(), KEYGEN_SECP256K1)
	
	if key.KeyType() == KEYGEN_SECP256K1 {
		secp256k1key := pubkey.(*secp256k1.PubKey)
		pub, err := btcsecp256k1.ParsePubKey(secp256k1key.Key, btcsecp256k1.S256())

		if err != nil {
			log.Printf("Not a secp256k1 key?")
		}
		
		log.Printf("Pub: %v", pub)

		// Validate the signature made by the HSM key, but using the
		// BTC secp256k1 public key
		valid := secp256k1key.VerifySignature(msg, signed)
		log.Printf("Did the signature verify? True = yes: %v", valid)

		log.Printf("TM blockchain address from pubkey: %v", secp256k1key.Address())
	}
	
	err = key.Delete()
	require.NoError(t, err)
}

func TestCreateKeySecp256r1(t *testing.T) {

	// hardcoded path for now - might change this API to actually take a JSON string and let caller decide how to get that JSON
	kr, err := NewPkcs11FromConfig("./pkcs11-config")
	require.NoError(t, err)
	label, err := randomBytes(16)
	require.NoError(t, err)

	key, err := kr.NewKey(KEYGEN_SECP256R1, string(label))
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, key.KeyType(), KEYGEN_SECP256R1)

	// Verify that I can also use this key in the X509 universe
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	log.Printf("Public: %s", pemEncodedPub)

	// retrieve the cosmos crypto pubkey
	pub := key.PubKey()

	// Tendermint address
	log.Printf("AccAddress: %v", sdk.AccAddress( pub.Address()))

	// point a second key object to the same key to test equality
	key2 := key

	log.Printf("Keys are equal (should be true)?: %v", key.Equals(*key2))

	label2, err := randomBytes(16)
	require.NoError(t, err)

	// Generate a completely different key
	key3, err := kr.NewKey( KEYGEN_SECP256R1, string(label2) )
	require.NoError(t, err)
	log.Printf("Keys are equal (should be false)?: %v", key.Equals(*key3))

	key4, err := kr.Key( string(label) )
	require.NoError(t, err)

	log.Printf("Keys are equal (should be true)?: %v", key4.Equals(*key))
	err = key.Delete()
	require.NoError(t, err)

	// This delete should fail since key2 is a pointer to key
	// which was already deleted
	err = key2.Delete()

	// Yes, there SHOULD be an error on this delete
	require.Error(t, err)

	// key3 delete should pass
	err = key3.Delete()
	require.NoError(t, err)

	// key4 should again be the same as a key already deleted, so
	// should fail
	err = key4.Delete()
	require.Error(t, err)
}
