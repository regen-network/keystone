The `keys` package allows the creation of cryptographic keys on a
PKCS11 accessible hardware token (HSM), such as AWS CloudHSM, YubiHSM
and Thales Luna devices. 

The package provides two APIs:

1. Keyring

A keyring in this case, is a device that contains one or more keys,
but is unwilling to share the private key bytes with the API
caller. Such devices include HSMs, secure enclaves, and the WebCrypto
feature of modern web browsers. This makes it different than the
Cosmos SDK keyring, which requires the key bytes to be shared with the
API caller.

2. Key

The `key` API itself combines functionality from three different
APIs - the PKCS11-based API provided by the Thales go-lang crypto11
package, the go-lang crypto.PublicKey API, and the Cosmos PubKey API
used to link public keys to their Cosmos account addresses.

## Building the package

`go build .` from within this directory, should be sufficient.

