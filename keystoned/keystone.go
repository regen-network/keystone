package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/regen-network/regen-ledger/x/group"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	acc "github.com/cosmos/cosmos-sdk/x/auth/types"

	keystonepb "github.com/regen-network/keystone/keystoned/proto"
)

const (
	//MY_ADDRESS = "regen1m0nxnyeq6ywy7hzlrvncytr4c5x5css3mhhxme" // delegator - signing fails
	MY_ADDRESS         = "regen19m2337xhcdd9ylwsxklcdeyanf25p6h266dd9m" // validator
	MY_CHAIN           = "test"
	MY_KEYRING_BACKEND = "test"
	MY_KEYRING_DIR     = "/Users/frumiousj/.regen/"
	MY_RPC_URI         = "tcp://localhost:26657"
)

type server struct{}

func createAddress() (address string){
	return MY_ADDRESS
}

func (s *server) Register(ctx context.Context, in *keystonepb.RegisterRequest) (*keystonepb.RegisterResponse, error) {
	log.Printf("Receive message body from client: %s %v", in.Address, in.EncryptedKey)

	var addr1 sdk.AccAddress = nil
	var err error
	
	if len(in.Address) > 0 {
		log.Printf("Address passed in request")
		addr1, err = sdk.AccAddressFromBech32(in.Address)

		if err != nil {
			log.Println("Address conversion from bech32 failed")
			return nil, err
		}
	} else {
		
		addr1, err = sdk.AccAddressFromBech32(createAddress())
		
		if err != nil {
			log.Println("Address conversion from bech32 failed")
			return nil, err
		}
	}	
	
	metadata := "some metadata"

	b := make([]byte, b64.StdEncoding.EncodedLen(len(metadata)))
	b64.StdEncoding.Encode(b, []byte(metadata))

	member1 := group.Member{
		Address:  addr1.String(),
		Weight:   strconv.Itoa(1),
		Metadata: b,
	}

	members := []group.Member{member1}

	localContext, err := getLocalContext()

	if err != nil {
		fmt.Println("Error getting local node context: ", err)
		return nil, err
	}

	// @@todo: create key/address and use that as creator address
	// may need to first create a key/address if one not sent in request
	// as "from" address (creator) must already exist
	groupAddress, err := CreateGroup([]byte(MY_ADDRESS), members, "", localContext)

	if err != nil {
		fmt.Println("Error creating group: ", err)
		return nil, err
	}

	log.Println("Group: ", groupAddress)

	return &keystonepb.RegisterResponse{Greeting: "Hello From the Server!", Status: 0}, nil
}

func (s *server) Sign(ctx context.Context, in *keystonepb.SignRequest) (*keystonepb.SignResponse, error) {
	return &keystonepb.SignResponse{Status: -1, SignedBytes: []byte{}}, nil
}

// go relayer/block explorer examples?

// how to retrieve node context beyond this one transaction?
func getLocalContext() (*client.Context, error) {

	addr, err := sdk.AccAddressFromBech32(MY_ADDRESS)
	encodingConfig := makeEncodingConfig()

	if err != nil {
		return nil, err
	}

	rpcclient, err := client.NewClientFromNode(MY_RPC_URI)

	if err != nil {
		return nil, err
	}

	k, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, MY_KEYRING_DIR, nil)

	// l, err := k.List()

	// fmt.Printf("%v", l)

	if err != nil {
		fmt.Printf("error opening keyring: ", err)
		return nil, err
	}

	c := client.Context{FromAddress: addr, ChainID: MY_CHAIN}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithBroadcastMode(flags.BroadcastSync).
		WithNodeURI(MY_RPC_URI).
		WithAccountRetriever(acc.AccountRetriever{}).
		WithClient(rpcclient).
		WithKeyringDir(MY_KEYRING_DIR).
		WithKeyring(k)

	return &c, nil
}

// something like this to abstract the tx building for multiple messages -- func createTx( txcfg params.EncodingConfig,

// CreateGroup creates a Cosmos Group using the MsgCreateGroup, filling the message with the input fields
func CreateGroup(creatorAddress []byte, memberList []group.Member, metadata string, localContext *client.Context) ([]byte, error) {

	encCfg := makeEncodingConfig()
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	// @@todo, how to get the private key from the keyring
	// associated with this address?
	adminAddr, err := sdk.AccAddressFromBech32(MY_ADDRESS)

	if err != nil {
		fmt.Printf("Error converting address string: ", err)
		return nil, err
	}

	err = localContext.AccountRetriever.EnsureExists(*localContext, adminAddr)

	if err != nil {
		fmt.Println("Account does not exist because: ", err)
		return nil, err
	}

	num, seq, err := localContext.AccountRetriever.GetAccountNumberSequence(*localContext, adminAddr)

	if err != nil {
		fmt.Printf("Error retrieving account number/sequence: ", err)
		return nil, err
	} else {
		fmt.Printf("Account retrieved: %v with seq: %v", num, seq)
	}

	//txBuilder := localContext.TxConfig.NewTxBuilder()
	txBuilder.SetMsgs(&group.MsgCreateGroup{
		Admin:    adminAddr.String(),
		Members:  memberList,
		Metadata: nil,
	})

	txBuilder.SetFeeAmount(sdk.Coins{sdk.NewInt64Coin("uregen", 5000)})
	txBuilder.SetGasLimit(50000)

	txFactory := clienttx.Factory{}
	txFactory = txFactory.
		WithChainID(localContext.ChainID).
		WithKeybase(localContext.Keyring).
		WithTxConfig(encCfg.TxConfig)

	// Only needed for "offline" accounts
	//.WithAccountNumber(num).WithSequence(seq)

	info, err := txFactory.Keybase().Key("delegator")

	if err != nil {
		return nil, err
	}

	fmt.Printf("%v", info)

	// NOT NEEDED IF USING SignTx from the x/auth/client, which
	// does all these things
	// signerData := xauthsigning.SignerData{
	// 	ChainID:       MY_CHAIN,
	// 	AccountNumber: num,
	// 	Sequence:      seq,
	// }

	// signBytes, err := localContext.TxConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder.GetTx())

	// if err != nil{
	// 	fmt.Println("Error getting signed bytes: ", err)
	// 	return nil, err
	// }

	txJSON, err := localContext.TxConfig.TxJSONEncoder()(txBuilder.GetTx())

	if err != nil {
		fmt.Println("Error getting JSON: ", err)
		return nil, err
	}

	fmt.Printf("Unsigned TX %s\n", txJSON)

	err = authclient.SignTx(txFactory, *localContext, "validator", txBuilder, false, true)

	if err != nil {
		fmt.Println("Error signing: ", err)
		return nil, err
	}

	txBytes, err := localContext.TxConfig.TxEncoder()(txBuilder.GetTx())

	if err != nil {
		fmt.Println("Error encoding transaction: ", err)
		return nil, err
	}

	txJSON, err = localContext.TxConfig.TxJSONEncoder()(txBuilder.GetTx())

	if err != nil {
		fmt.Println("Error getting JSON: ", err)
		return nil, err
	}

	fmt.Printf("Signed TX %s\n", txJSON)

	//res, err := localContext.BroadcastTx(txBytes)

	opts := grpc.WithInsecure()

	grpcConn, err := grpc.Dial("127.0.0.1:9090", opts)

	if err != nil {
		fmt.Println("Err doing grpc dial: ", err)
		return nil, err
	}

	defer grpcConn.Close()

	txClient := tx.NewServiceClient(grpcConn)
	// We then call the BroadcastTx method on this client.
	res, err := txClient.BroadcastTx(
		context.Background(),
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
		},
	)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.TxResponse.Code) // Should be `0` if the tx is successful

	if err != nil {
		fmt.Println("Error broadcasting ", err)
	}

	fmt.Printf("Result: %v", res)

	return []byte{}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	keystonepb.RegisterKeystoneServiceServer(s, &server{})

	s.Serve(lis)

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
