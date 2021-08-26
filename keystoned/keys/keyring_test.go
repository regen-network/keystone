package keys

import (
	"log"
	"testing"
	"crypto/x509"
	"encoding/pem"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

)

func TestCreateKeySecp256k1(t *testing.T) {

	// hardcoded path for now - might change this API to actually take a JSON string and let caller decide how to get that JSON
	kr, err := NewPkcs11FromConfig("./pkcs11-config")
	require.NoError(t, err)
	label, err := randomBytes(16)
	require.NoError(t, err)
	
	key, err := kr.NewKey( KEYGEN_SECP256K1, string(label) )
	require.NoError(t, err)
	require.NotNil(t, key)

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
	log.Printf("Address: %s", string(key.PubKey().Address()))
	err = key.Delete()
	require.NoError(t, err)
}
