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

//var feePayer = types.AccountFromPrivateKeyBytes([]byte{209, 25, 205, 128, 167, 154, 29, 44, 118, 161, 221, 56, 195, 253, 106, 150, 177, 135, 44, 222, 168, 0, 34, 39, 135, 108, 78, 107, 243, 139, 41, 55, 236, 166, 114, 58, 45, 235, 17, 246, 124, 48, 45, 88, 100, 215, 166, 123, 17, 131, 143, 26, 155, 72, 7, 203, 91, 245, 125, 83, 153, 238, 50, 186})

var alice = types.AccountFromPrivateKeyBytes([]byte{})

var aliceTokenAccount = types.AccountFromPrivateKeyBytes([]byte{})

var mintPubkey = common.PublicKeyFromString("9WKBr1Gzt6YS1g4XMGHxTwZsmpyi5DHCyxqi7GuyNysC")

func main() {
	c := client.NewClient("https://api.devnet.solana.com")

	//aliceTokenAccount := types.NewAccount()
	fmt.Println("alice account:", alice.PublicKey.ToBase58())
	fmt.Println("alice token account:", aliceTokenAccount.PublicKey.ToBase58())

	rentExemptionBalance, err := c.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.TokenAccountSize)
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
				//feePayer.PublicKey,
				alice.PublicKey,
				aliceTokenAccount.PublicKey,
				common.TokenProgramID,
				rentExemptionBalance,
				tokenprog.TokenAccountSize,
			),
			tokenprog.InitializeAccount(
				aliceTokenAccount.PublicKey,
				mintPubkey,
				alice.PublicKey,
			),
		},
		//Signers:         []types.Account{feePayer, aliceTokenAccount},
		Signers:         []types.Account{alice, aliceTokenAccount},
		//FeePayer:        feePayer.PublicKey,
		FeePayer:        alice.PublicKey,//确保付费地址放在数组最前面
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
