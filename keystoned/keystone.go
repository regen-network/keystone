package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"flag"

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

type server struct{
	ServerAddress    string
	ChainID          string
	KeyringType      string
	KeyringDir       string
	RpcURI           string
}

// address will eventually create a keypair (in an HSM via the
// key/keyring struct) and then create the address, derived from the
// public key in the usual way
func address() (address string){
	return MY_ADDRESS
}

//adminMembers returns a []group.Member with two members
func adminMembers( addr1 string, addr2 string ) []group.Member{
	
	member1 := group.Member{
		Address:  addr1,
		Weight:   strconv.Itoa(1),
	}

	member2 := group.Member{
		Address:  addr2,
		Weight:   strconv.Itoa(1),
	}
	
	return []group.Member{member1, member2}
}

// Register implements the method given in the protobuf definition for
// the Keystone service (proto/keystone.proto)
func (s *server) Register(ctx context.Context, in *keystonepb.RegisterRequest) (*keystonepb.RegisterResponse, error) {
	log.Printf("Receive message body from client: %s %v", in.Address, in.EncryptedKey)

	var addr1 sdk.AccAddress = nil
	var err error

	// If an address is passed in via the request, then use that address
	// to create the group, otherwise create a new address for the group
	
	if len(in.Address) > 0 {
		log.Printf("Address passed in request")
		addr1, err = sdk.AccAddressFromBech32(in.Address)

		if err != nil {
			log.Println("Address conversion from bech32 failed")
			return nil, err
		}
	} else {
		
		addr1, err = sdk.AccAddressFromBech32(s.ServerAddress)
		
		if err != nil {
			log.Println("Address conversion from bech32 failed")
			return nil, err
		}
	}	
	
	localContext, err := getLocalContext(*s)

	if err != nil {
		fmt.Println("Error getting local node context: ", err)
		return nil, err
	}

	// @@todo: create key/address and use that as creator address
	// may need to first create a key/address if one not sent in request
	// as "from" address (creator) must already exist
	groupAddress, err := createGroup([]byte(addr1.String()), adminMembers(addr1.String(), addr1.String()), "", localContext)

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
func getLocalContext(s server) (*client.Context, error) {

	addr, err := sdk.AccAddressFromBech32(s.ServerAddress)
	encodingConfig := makeEncodingConfig()

	if err != nil {
		return nil, err
	}

	rpcclient, err := client.NewClientFromNode(s.RpcURI)

	if err != nil {
		return nil, err
	}

	//@@TODO configure keyring.BackendTest using the server-global context, not hardcode
	k, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, s.KeyringDir, nil)

	// l, err := k.List()

	// fmt.Printf("%v", l)

	if err != nil {
		fmt.Printf("error opening keyring: ", err)
		return nil, err
	}

	c := client.Context{FromAddress: addr, ChainID: s.ChainID}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithBroadcastMode(flags.BroadcastSync).
		WithNodeURI(s.RpcURI).
		WithAccountRetriever(acc.AccountRetriever{}).
		WithClient(rpcclient).
		WithKeyringDir(s.KeyringDir).
		WithKeyring(k)

	return &c, nil
}

func createAdminGroup(creatorAddress []byte, memberList []group.Member, metadata string, localContext *client.Context) ([]byte, error) {
	return createGroup( creatorAddress, memberList, metadata, localContext )
}
	
// something like this to abstract the tx building for multiple messages -- func createTx( txcfg params.EncodingConfig,

// CreateGroup creates a Cosmos Group using the MsgCreateGroup, filling the message with the input fields
func createGroup(creatorAddress []byte, memberList []group.Member, metadata string, localContext *client.Context) ([]byte, error) {

	encCfg := makeEncodingConfig()
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	// @@todo, how to get the private key from the keyring
	// associated with this address?
	adminAddr, err := sdk.AccAddressFromBech32(string(creatorAddress))

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

	// Only needed for "offline" accounts?
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

	// @@TODO: use secure connection?
	opts := grpc.WithInsecure()

	// @@TODO: configure the dial location from server context
	grpcConn, err := grpc.Dial("127.0.0.1:9090", opts)

	if err != nil {
		fmt.Println("Err doing grpc dial: ", err)
		return nil, err
	}

	defer grpcConn.Close()

	// @@TODO: configure broadcast mode from server-global context?
	
	txClient := tx.NewServiceClient(grpcConn)

	res, err := txClient.BroadcastTx(
		context.Background(),
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
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

	// Retrieve the command line parameters passed in to configure the server
	// Most have likely-reasonable defaults.
	keystoneAddress := flag.String("key-addr", "", "the address associated with the key used to sign transactions on behalf of Keystone")
	blockchain := flag.String("chain-id", "test-chain", "the blockchain that Keystone should connect to")
	keyringType := flag.String("keyring-type", "test", "the keyring backend type where keys should be read from")
	keyringDir := flag.String("keyring-dir", "~/.regen/", "the directory where the keys are")
	chainRpcURI := flag.String("chain-rpc", "tcp://localhost:26657", "the address of the RPC endpoint to communicate with the blockchain")
	grpcListenPort := flag.String("listen-port", "8080", "the port where the server will listen for connections")

	flag.Parse()

	if len(*keystoneAddress) <= 0 {
		log.Fatalln("Keystone server blockchain address may not be left empty")
		return
	}

	lis, err := net.Listen("tcp", ":" + *grpcListenPort)

	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create new server context, used for passing server-global state
	ss := server{
		ServerAddress: *keystoneAddress,
		ChainID: *blockchain,
		KeyringType: *keyringType,
		KeyringDir: *keyringDir,
		RpcURI: *chainRpcURI,
	}
	
	s := grpc.NewServer()
	keystonepb.RegisterKeystoneServiceServer(s, &ss)

	s.Serve(lis)
	return

}
