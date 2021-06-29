package main

import (
	"context"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

//var feePayer = types.AccountFromPrivateKeyBytes([]byte{209, 25, 205, 128, 167, 154, 29, 44, 118, 161, 221, 56, 195, 253, 106, 150, 177, 135, 44, 222, 168, 0, 34, 39, 135, 108, 78, 107, 243, 139, 41, 55, 236, 166, 114, 58, 45, 235, 17, 246, 124, 48, 45, 88, 100, 215, 166, 123, 17, 131, 143, 26, 155, 72, 7, 203, 91, 245, 125, 83, 153, 238, 50, 186})
var feePayer = types.AccountFromPrivateKeyBytes([]byte{52,97,70,29,114,88,217,212,239,177,207,64,16,4,223,160,132,156,96,47,142,2,175,155,58,69,252,17,54,178,216,101,104,160,49,117,132,180,116,226,115,206,138,198,14,92,86,2,163,251,80,137,146,193,213,39,155,97,161,62,43,99,11,74})

//var alice = types.AccountFromPrivateKeyBytes([]byte{196, 114, 86, 165, 59, 177, 63, 87, 43, 10, 176, 101, 225, 42, 129, 158, 167, 43, 81, 214, 254, 28, 196, 158, 159, 64, 55, 123, 48, 211, 78, 166, 127, 96, 107, 250, 152, 133, 208, 224, 73, 251, 113, 151, 128, 139, 86, 80, 101, 70, 138, 50, 141, 153, 218, 110, 56, 39, 122, 181, 120, 55, 86, 185})

var mintPubkey = common.PublicKeyFromString("9WKBr1Gzt6YS1g4XMGHxTwZsmpyi5DHCyxqi7GuyNysC")

var feeTokenATAPubkey = common.PublicKeyFromString("AKpvREJ4JpdtHYWuk4W5J2Sro1ZfN11q7kuN7omGRCAC")

var aliceTokenATAPubkey = common.PublicKeyFromString("ACzsMjMGJkdaUNQ1LtGa7bFe3GLWorv6PVdTCpXJxfdA")

func main() {
	c := client.NewClient("https://api.devnet.solana.com")

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}
	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			tokenprog.TransferChecked(
				feeTokenATAPubkey,
				aliceTokenATAPubkey,
				mintPubkey,
				//alice.PublicKey,
				feePayer.PublicKey,
				[]common.PublicKey{},
				3e8,
				9,
			),
		},
		//Signers:         []types.Account{feePayer, alice},
		Signers:         []types.Account{feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	txhash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	log.Println("txhash:", txhash)
}
