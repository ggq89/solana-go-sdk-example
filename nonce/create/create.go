package main

import (
	"context"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/types"
)

var feePayer = types.AccountFromPrivateKeyBytes([]byte{178, 244, 76, 4, 247, 41, 113, 40, 111, 103, 12, 76, 195, 4, 100, 123, 88, 226, 37, 56, 209, 180, 92, 77, 39, 85, 78, 202, 121, 162, 88, 29, 125, 155, 223, 107, 139, 223, 229, 82, 89, 209, 27, 43, 108, 205, 144, 2, 74, 159, 215, 57, 198, 4, 193, 36, 161, 50, 160, 119, 89, 240, 102, 184})

func main() {
	c := client.NewClient("https://api.devnet.solana.com")

	nonceAccountRentFreeBalance, err := c.GetMinimumBalanceForRentExemption(
		context.Background(),
		sysprog.NonceAccountSize,
	)
	if err != nil {
		log.Fatalf("failed to get min balance for nonce account, err: %v", err)
	}

	nonceAccount := types.NewAccount()
	fmt.Println("nonce account:", nonceAccount.PublicKey)

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		FeePayer: feePayer.PublicKey,
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				nonceAccount.PublicKey,
				common.SystemProgramID,
				nonceAccountRentFreeBalance,
				sysprog.NonceAccountSize,
			),
			sysprog.InitializeNonceAccount(
				// nonce account
				nonceAccount.PublicKey,
				// nonce account's owner
				feePayer.PublicKey,
			),
		},
		Signers: []types.Account{feePayer, nonceAccount},
		RecentBlockHash: res.Blockhash,
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
