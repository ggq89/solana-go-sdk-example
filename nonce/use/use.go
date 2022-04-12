package main

import (
	"context"
	"encoding/base64"
	"fmt"
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

	nonceAccount, err := getNonceAccount(c, nonceAccountPubkey.ToBase58())
	if err != nil {
		log.Fatalf("failed to getNonceAccount, err: %v", err)
	}

	log.Printf("nonce: %+v\n", nonceAccount.Nonce.ToBase58())

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
				1e7, // 1 SOL
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

	log.Println("txhash: ", txhash)

	for true {
		resp, err := c.GetSignatureStatuses(context.Background(), []string{txhash})
		if err!= nil {
			log.Fatalf("GetSignatureStatuses err: %v\n", err)
			continue
		}

		log.Printf("%+v\n", resp)

		if len(resp) == 0 {
			log.Fatalf("GetSignatureStatuses result len 0")
			continue
		}

		if resp[0].ConfirmationStatus != nil {
			status := *(resp[0].ConfirmationStatus)
			log.Printf("tx ConfirmationStatus: %s\n", status)
			if status == client.CommitmentFinalized {
				break
			}
		}


		nonceAccount, err := getNonceAccount(c, nonceAccountPubkey.ToBase58())
		if err != nil {
			log.Fatalf("failed to deserialize nonce account, err: %v", err)
			continue
		}

		log.Printf("nonce account state: %v\n", nonceAccount.State)
		log.Printf("nonce: %v\n", nonceAccount.Nonce.ToBase58())
	}
}

func getNonceAccount(c *client.Client, nonceAccountAddr string) (*sysprog.NonceAccount, error) {
	// fetch nonce
	cfg := client.GetAccountInfoConfig{
		client.GetAccountInfoConfigEncodingBase64,
		client.GetAccountInfoConfigDataSlice{},
	}
	accountInfo, err := c.GetAccountInfo(
		context.Background(),
		nonceAccountAddr,
		cfg,
	)
	if err != nil {
		log.Fatalf("failed to get account info, err: %v", err)
		return nil, err
	}

	data, ok := accountInfo.Data.([]interface{})
	if !ok {
		errStr := "failed to cast raw response to []interface{}"
		log.Fatalf(errStr)
		return nil, fmt.Errorf(errStr)
	}

	rawData, err := base64.StdEncoding.DecodeString(data[0].(string))
	if err != nil {
		log.Fatalf("failed to base64 decode data")
		return nil, err
	}

	nonceAccount, err := sysprog.NonceAccountDeserialize(rawData)
	if err != nil {
		log.Fatalf("failed to deserialize nonce account, err: %v", err)
		return nil, err
	}

	return &nonceAccount, nil
}
