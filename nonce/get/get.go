package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/sysprog"
)

func main() {
	c := client.NewClient("https://api.devnet.solana.com")


	cfg := client.GetAccountInfoConfig{
		client.GetAccountInfoConfigEncodingBase64,
		client.GetAccountInfoConfigDataSlice{},
	}
	accountInfo, err := c.GetAccountInfo(
		context.Background(),
		"8NoQpwqqEwHiTPA8WrFEQUEJ425VcCJu9wzKSwA1oUwC",
		cfg,
	)
	if err != nil {
		log.Fatalf("failed to get account info, err: %v", err)
	}

	data, ok := accountInfo.Data.([]interface{})
	if !ok {
		log.Fatalf("failed to cast raw response to []interface{}")
	}

	rawData, err := base64.StdEncoding.DecodeString(data[0].(string))
	if err != nil {
		log.Fatalf("failed to base64 decode data")
	}

	nonceAccount, err := sysprog.NonceAccountDeserialize(rawData)
	if err != nil {
		log.Fatalf("failed to deserialize nonce account, err: %v", err)
	}

	/*
		type NonceAccount struct {
			Version          uint32
			State            uint32
			AuthorizedPubkey common.PublicKey
			Nonce            common.PublicKey
			FeeCalculator    FeeCalculator
		}
	*/
	fmt.Printf("%+v\n", nonceAccount)
	fmt.Printf("nonce: %+v\n", nonceAccount.Nonce.ToBase58())
}
