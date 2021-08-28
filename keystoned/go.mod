module github.com/regen-network/keystone/keystoned

go 1.16

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/regen-network/regen-ledger/orm => github.com/regen-network/regen-ledger/orm v0.0.0-20210804173213-3265a868bf83

replace github.com/regen-network/regen-ledger/types => github.com/regen-network/regen-ledger/types v0.0.0-20210804173213-3265a868bf83

require (
	github.com/ThalesIgnite/crypto11 v1.2.4 // indirect
	github.com/cosmos/cosmos-sdk v0.43.0-rc0
	github.com/regen-network/keystone/keystoned/keys v0.0.0-20210823173722-5de519e0203f
	github.com/regen-network/regen-ledger/x/group v0.0.0-20210804173213-3265a868bf83
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210804223703-f1db76f3300d // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)
