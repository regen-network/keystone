package main

import (
	"fmt"
	b64 "encoding/base64"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/regen-network/regen-ledger/x/group"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MY_ADDRESS = "regen1w6j8wsmwzc0s474dqpe6kqtdmeaqz8rwnuzfrp"
const MY_CHAIN = "test"
const MY_KEYRING_BACKEND = "test"
const MY_KEYRING_DIR = "/Users/frumiousj/.regen/keyring-test"

type Member struct {
	Address    string
	Weight     int
	Metadata   string
}

// group.MsgCreateGroup{
// 		Admin:    s.addr1.String(),
// 		Members:  members,
// 		Metadata: nil,
// })

// how to retrieve node context beyond this one transaction?
func getLocalContext() (*client.Context, error) {

	addr, err := sdk.AccAddressFromBech32(MY_ADDRESS)
	encodingConfig := MakeEncodingConfig()
	
	if err != nil {
		return nil, err
	}
	
	c := client.Context{FromAddress: addr, ChainID: MY_CHAIN,}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig)

	k, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, MY_KEYRING_DIR, nil)

	if err != nil {
		fmt.Printf("error opening keyring: ", err)
		return nil, err
	} else {
		fmt.Printf("no error with keyring")
	}
	
	c = c.WithKeyring(k)
	fmt.Printf("%+v\n", c)	
	return &c, nil
}

func CreateGroup(creatorAddress []byte, memberList []group.Member, metadata string, localContext *client.Context ) ([]byte, error){

	adminAddr, err := sdk.AccAddressFromBech32(MY_ADDRESS)

	if err != nil {
		fmt.Printf("Error converting address string: ", err)
		return nil, err
	}
	
	txBuilder := localContext.TxConfig.NewTxBuilder()
	txBuilder.SetMsgs(&group.MsgCreateGroup{
		Admin: adminAddr.String(),
		Members: memberList,
		Metadata: nil,
	})
	
	return []byte{}, nil
}

func main() {

	addr1, err := sdk.AccAddressFromBech32(MY_ADDRESS)

	var metadata []byte
	b64.StdEncoding.Encode(metadata, []byte("some metadata"))
	
	member1 := group.Member{
		Address: addr1.String(),
		Weight: strconv.Itoa(1),
		Metadata: metadata,
	}

	members := []group.Member{ member1 }

	localContext, err := getLocalContext()

	if err != nil {
		fmt.Println("Error getting local node context: ", err)
		return
	}

	// may need to first create a key/address if one not sent in request
	// as "from" address (creator) must already exist
	groupAddress, err := CreateGroup( []byte(MY_ADDRESS), members, "", localContext )

	if err != nil {
		fmt.Println("Error creating group: ", err)
		return
	}

	fmt.Println("Group: ", groupAddress)

	// Generate the signing payload
	// signerData := authsign.SignerData{
	// 	ChainID:       "regen-network-devnet-5",
	// 	AccountNumber: 2,
	// 	Sequence:      2,
	// }
	// signBytes, err := suite.clientCtx.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signerData, tx.GetTx())
	// suite.txBuilder.GetTx()
	// Sign the signBytes
	// setup txFactory
	// txFactory := tx.Factory{}.
	// 	WithChainID(suite.ClientCtx.ChainID).
	// 	WithKeybase(suite.ClientCtx.Keyring).
	// 	WithTxConfig(suite.ClientCtx.TxConfig).
	// 	WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	// Sign Tx.
	//err := client.SignTx(txFactory, val.ClientCtx, val.Moniker, suite.txBuilder, false, true)
	// Build the TxRaw bytes.
	//bz, err := suite.clientCtx.TxConfig.TxEncoder()(suite.txBuilder)
	// Broadcast Tx.
	//res, err := suite.clientCtx.BroadcastTx(bz)
	return

}
