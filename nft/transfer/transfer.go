package main

import (
	"context"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

// 83R5RVHMEEmHtj9QydfAX958JDoNHKREmQhw8k24ryMj
var  alice, _ = types.AccountFromBase58("23jvmogkErSMC9rR8nShPP15yWDdN8hDzcYzS2ixeE4QnvfPcyD8oBuS2mwK2v31nG24pqN9kevn6Xw3Mna8P4U1")

// 9TKpNubehUkr8PfP6vCoWUuNejiNw9y21MLcCsEUUSPq
var feePayer, _ = types.AccountFromBytes([]byte{178, 244, 76, 4, 247, 41, 113, 40, 111, 103, 12, 76, 195, 4, 100, 123, 88, 226, 37, 56, 209, 180, 92, 77, 39, 85, 78, 202, 121, 162, 88, 29, 125, 155, 223, 107, 139, 223, 229, 82, 89, 209, 27, 43, 108, 205, 144, 2, 74, 159, 215, 57, 198, 4, 193, 36, 161, 50, 160, 119, 89, 240, 102, 184})

var mintPubkey = common.PublicKeyFromString("6jrMWzLceTMKVWS22r1fNZWe7pEmqk41a9sRjuabn76w")

var toPubkey = common.PublicKeyFromString("9CfNE5H21Rqh1Wev28vZ3f41315GU6uz8khxt6ihh9Gs")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	nonceAccountPubkey := common.PublicKeyFromString("8NoQpwqqEwHiTPA8WrFEQUEJ425VcCJu9wzKSwA1oUwC")

	isCreated := true

	toATA, _, err := common.FindAssociatedTokenAddress(toPubkey, mintPubkey)
	if err != nil {
		log.Fatalf("find associated token account error, err: %s", err.Error())
	}
	fmt.Printf("toATA: %v\n", toATA.String())

	toATAInfo, err := c.GetAccountInfo(context.Background(),toATA.String())
	if err != nil {
		log.Fatalf("get to ATA info err:%s", err.Error())
	}

	if len(toATAInfo.Data) == 0 {
		isCreated = false
	}

	nftPubkey, err := getNftAccount(c, mintPubkey.String(), alice.PublicKey.String())
	if err != nil {
		log.Fatalf("getNftAccount err:%s", err.Error())
	}

	instructions := []types.Instruction{
		tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
			From:     common.PublicKeyFromString(nftPubkey),
			To:       toATA,
			Mint:     mintPubkey,
			Auth:     alice.PublicKey,
			Signers:  []common.PublicKey{},
			Amount:   1,
			Decimals: 0,
		}),
	}

	// 判断ata地址是否创建
	if !isCreated {
		instructions = append([]types.Instruction{
			assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
				Funder:                 feePayer.PublicKey,
				Owner:                  toPubkey,
				Mint:                   mintPubkey,
				AssociatedTokenAccount: toATA,
			}),
		}, instructions...)
	}

	// make nonce advance instruction is the first instruction
	instructions = append([]types.Instruction{
		sysprog.AdvanceNonceAccount(sysprog.AdvanceNonceAccountParam{
			Nonce: nonceAccountPubkey,
			Auth:  feePayer.PublicKey, // nonce account's owner
		}),
	}, instructions...)

	nonce, err := c.GetNonceFromNonceAccount(context.Background(), nonceAccountPubkey.String())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}
	fmt.Printf("nonce: %v\n", nonce)

	txParam := types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:     feePayer.PublicKey,
			RecentBlockhash: nonce,
			Instructions: instructions,
		}),
		Signers: []types.Account{alice, feePayer},
	}



	tx, err := types.NewTransaction(txParam)
	if err != nil {
		log.Fatalf("generate tx error, err: %s", err.Error())
	}

	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	fmt.Println(sig)
}

func getNftAccount(c *client.Client,  mint, owner string) (string, error) {
	resp, err  := c.RpcClient.GetTokenAccountsByOwner(context.Background(),owner,
		rpc.GetTokenAccountsByOwnerConfigFilter{
			Mint: mint,
			//ProgramId: common.TokenProgramID.ToBase58(),
		})
	if err != nil {
		log.Fatalf("failed GetTokenAccountsByOwner, err: %v", err)
		return "", err
	}
	if resp.Error != nil {
		errStr := fmt.Sprintf("rpc response error: %v", resp.Error)
		log.Fatalf("failed GetTokenAccountsByOwner, err: %v", errStr)
		return "", fmt.Errorf(errStr)
	}

	for _,programAccount := range resp.Result.Value {
		balance, _ ,err := c.GetTokenAccountBalance(context.Background(),programAccount.Pubkey)
		if err != nil {
			log.Fatalf("failed GetTokenAccountBalance, err: %v", err)
			return "", err
		}

		if balance >0 {
			return programAccount.Pubkey, nil
		}
	}

	return "", fmt.Errorf("cannot find account")
}
