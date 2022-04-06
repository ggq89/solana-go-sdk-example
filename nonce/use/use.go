package main

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/types"
)

var feePayer = types.AccountFromPrivateKeyBytes([]byte{178, 244, 76, 4, 247, 41, 113, 40, 111, 103, 12, 76, 195, 4, 100, 123, 88, 226, 37, 56, 209, 180, 92, 77, 39, 85, 78, 202, 121, 162, 88, 29, 125, 155, 223, 107, 139, 223, 229, 82, 89, 209, 27, 43, 108, 205, 144, 2, 74, 159, 215, 57, 198, 4, 193, 36, 161, 50, 160, 119, 89, 240, 102, 184})

func main() {
	// you created before
	nonceAccountPubkey := common.PublicKeyFromString("8NoQpwqqEwHiTPA8WrFEQUEJ425VcCJu9wzKSwA1oUwC")

	c := client.NewClient("https://api.devnet.solana.com")

	// fetch nonce
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

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		FeePayer: feePayer.PublicKey,
		Instructions: []types.Instruction{
			sysprog.AdvanceNonceAccount(
				nonceAccountPubkey,
				feePayer.PublicKey,
			),
			sysprog.Transfer(
				//alice.PublicKey, // from
				feePayer.PublicKey,
				common.PublicKeyFromString("83R5RVHMEEmHtj9QydfAX958JDoNHKREmQhw8k24ryMj"), // to
				1e9, // 1 SOL
			),
		},
		Signers: []types.Account{feePayer},
		RecentBlockHash: nonceAccount.Nonce.ToBase58(),
	},
	)
	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	txhash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	log.Println("txhash:", txhash)
}
