package main

import (
	"context"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

var feePayer = types.AccountFromPrivateKeyBytes([]byte{})

//var alice = types.AccountFromPrivateKeyBytes([]byte{196, 114, 86, 165, 59, 177, 63, 87, 43, 10, 176, 101, 225, 42, 129, 158, 167, 43, 81, 214, 254, 28, 196, 158, 159, 64, 55, 123, 48, 211, 78, 166, 127, 96, 107, 250, 152, 133, 208, 224, 73, 251, 113, 151, 128, 139, 86, 80, 101, 70, 138, 50, 141, 153, 218, 110, 56, 39, 122, 181, 120, 55, 86, 185})

var mint = types.AccountFromPrivateKeyBytes([]byte{})

func main() {
	c := client.NewClient("https://api.devnet.solana.com")

	//mint := types.NewAccount()
	//fmt.Println("mint:", mint.PublicKey.ToBase58())

	rentExemptionBalance, err := c.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		log.Fatalf("get min balacne for rent exemption, err: %v", err)
	}
	fmt.Println("rent exemption balance:", rentExemptionBalance)

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}
	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				mint.PublicKey,
				common.TokenProgramID,
				rentExemptionBalance,
				tokenprog.MintAccountSize,
			),
			tokenprog.InitializeMint(
				9,
				mint.PublicKey,
				//alice.PublicKey,
				common.PublicKeyFromString("9CfNE5H21Rqh1Wev28vZ3f41315GU6uz8khxt6ihh9Gs"),
				common.PublicKey{},
			),
		},
		//Signers:         []types.Account{feePayer, mint},
		Signers:         []types.Account{feePayer, mint},
		//FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	txhash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	fmt.Println("txhash:", txhash)
}
