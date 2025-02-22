package main

import (
	"context"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/assotokenprog"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

//var feePayer = types.AccountFromPrivateKeyBytes([]byte{209, 25, 205, 128, 167, 154, 29, 44, 118, 161, 221, 56, 195, 253, 106, 150, 177, 135, 44, 222, 168, 0, 34, 39, 135, 108, 78, 107, 243, 139, 41, 55, 236, 166, 114, 58, 45, 235, 17, 246, 124, 48, 45, 88, 100, 215, 166, 123, 17, 131, 143, 26, 155, 72, 7, 203, 91, 245, 125, 83, 153, 238, 50, 186})

//var alice = types.AccountFromPrivateKeyBytes([]byte{196, 114, 86, 165, 59, 177, 63, 87, 43, 10, 176, 101, 225, 42, 129, 158, 167, 43, 81, 214, 254, 28, 196, 158, 159, 64, 55, 123, 48, 211, 78, 166, 127, 96, 107, 250, 152, 133, 208, 224, 73, 251, 113, 151, 128, 139, 86, 80, 101, 70, 138, 50, 141, 153, 218, 110, 56, 39, 122, 181, 120, 55, 86, 185})
var alice = types.AccountFromPrivateKeyBytes([]byte{})

var mintPubkey = common.PublicKeyFromString("9WKBr1Gzt6YS1g4XMGHxTwZsmpyi5DHCyxqi7GuyNysC")

func main() {
	c := client.NewClient("https://api.devnet.solana.com")

	ata, _, err := common.FindAssociatedTokenAddress(alice.PublicKey, mintPubkey)
	if err != nil {
		log.Fatalf("find ata error, err: %v", err)
	}
	fmt.Println("ata:", ata.ToBase58())

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}
	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			assotokenprog.CreateAssociatedTokenAccount(
				//feePayer.PublicKey,
				alice.PublicKey,
				alice.PublicKey,
				mintPubkey,
			),
		},
		//Signers:         []types.Account{feePayer},
		Signers:         []types.Account{alice},
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
