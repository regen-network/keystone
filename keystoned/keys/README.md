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

## Getting started

0. git clone this repository.

1. Install SoftHSM (https://github.com/opendnssec/SoftHSMv2)

    Follow the instructions in the SoftHSM repository to build SoftHSM. Personally my configure line was:

    `./configure --enable-ecc --enable-eddsa --disable-gost`

    I run `make install` to install SoftHSM in `/usr/local/lib/softhsm`

    Once you've made the software, you need to create a token (an HSM is represented to the client as a "token"):

    This is done using the following command (as shown in the SoftHSM README):

    `softhsm2-util --init-token --slot 0 --label "My token 1"`

    The token label will then be used in your configuration file (see step 3. below) - make sure the value in your configuration file matches the one you passed into the `init-token` step.

2. Install pkcs11-tool as an independent way of verifying you have access to the token, and can see keys

    pkcs11-tool is part of the opensc package, which may already be on your machine if you are using a Linux flavour. If not, install OpenSC from https://github.com/OpenSC/OpenSC/wiki

3. Create the PKCS11 configuration file with appropriate values

    This repository contains a file called pkcs11-config.template. Configure the file according to your own SoftHSM setup. My example is shown here below:

```
    {
      "Path": "/usr/local/lib/softhsm/libsofthsm2.so",
      "TokenLabel": "The Cosmos",
      "Pin": "5565455835367496544668"
    }
```

## Building the package

`go build .` from within this directory, should be sufficient.

## Running the tests

The tests for the package are in `keyring_test.go`, and may be run via `go test`, if you have previously configured a `pkcs11-config` file, present in the local directory, as described above.
