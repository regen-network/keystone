package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	//"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	grptypes "github.com/regen-network/regen-ledger/x/group"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// MakeEncodingConfig creates an EncodingConfig
// registering all of the types that are needed
// for the caller to call the APIs it does
func makeEncodingConfig() params.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	authtypes.RegisterInterfaces(interfaceRegistry)
	cryptotypes.RegisterInterfaces(interfaceRegistry)
	grptypes.RegisterTypes(interfaceRegistry)

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	
	return encodingConfig
}
