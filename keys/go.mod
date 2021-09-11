module github.com/regen-network/keystone/keys

go 1.16

require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/cosmos/cosmos-sdk v0.43.0
	github.com/frumioj/crypto11 v1.2.5-0.20210823151709-946ce662cc0e
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.12 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
