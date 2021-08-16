package main

import (
	"log"
	"context"
	"fmt"
	
	"google.golang.org/grpc"
	
	keystonepb "github.com/regen-network/keystone/keystoned/proto"
)

func main() {
	
	fmt.Println("Keystone client ...")

	opts := grpc.WithInsecure()
	
	cc, err := grpc.Dial("localhost:8080", opts)
	
	if err != nil {
		log.Fatal(err)
	}
	
	defer cc.Close()

	client := keystonepb.NewKeystoneServiceClient(cc)

	request := &keystonepb.RegisterRequest{Address: "regen19m2337xhcdd9ylwsxklcdeyanf25p6h266dd9m"}

	resp, _ := client.Register(context.Background(), request)
	
	fmt.Printf("Receive response => [%v]", resp.Greeting)

	cleartext := "For signing"
	
	signRequest := &keystonepb.SignRequest{ForSigning: []byte(cleartext)}

	signResp, _ := client.Sign(context.Background(), signRequest)

	fmt.Printf("Signing response => [%v]", signResp.Status)
}
