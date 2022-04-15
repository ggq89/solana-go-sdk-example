package main

import (
	"context"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/pointer"
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

// 83R5RVHMEEmHtj9QydfAX958JDoNHKREmQhw8k24ryMj
var  alice, _ = types.AccountFromBase58("23jvmogkErSMC9rR8nShPP15yWDdN8hDzcYzS2ixeE4QnvfPcyD8oBuS2mwK2v31nG24pqN9kevn6Xw3Mna8P4U1")

// 9TKpNubehUkr8PfP6vCoWUuNejiNw9y21MLcCsEUUSPq
var feePayer, _ = types.AccountFromBytes([]byte{178, 244, 76, 4, 247, 41, 113, 40, 111, 103, 12, 76, 195, 4, 100, 123, 88, 226, 37, 56, 209, 180, 92, 77, 39, 85, 78, 202, 121, 162, 88, 29, 125, 155, 223, 107, 139, 223, 229, 82, 89, 209, 27, 43, 108, 205, 144, 2, 74, 159, 215, 57, 198, 4, 193, 36, 161, 50, 160, 119, 89, 240, 102, 184})

func main() {
	// you created before
	nonceAccountPubkey := common.PublicKeyFromString("8NoQpwqqEwHiTPA8WrFEQUEJ425VcCJu9wzKSwA1oUwC")

	c := client.NewClient(rpc.DevnetRPCEndpoint)

	mint := types.NewAccount()
	fmt.Printf("NFT: %v\n", mint.PublicKey.ToBase58())

	collection := types.NewAccount()
	//fmt.Println(base58.Encode(collection.PrivateKey))
	fmt.Printf("collection: %v\n", collection.PublicKey.ToBase58())

	ata, _, err := common.FindAssociatedTokenAddress(alice.PublicKey, mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid ata, err: %v", err)
	}
	fmt.Printf("ata: %v\n", ata.ToBase58())

	tokenMetadataPubkey, err := tokenmeta.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid token metadata, err: %v", err)

	}
	fmt.Printf("tokenMetadataPubkey: %v\n", tokenMetadataPubkey.ToBase58())

	tokenMasterEditionPubkey, err := tokenmeta.GetMasterEdition(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid master edition, err: %v", err)
	}
	fmt.Printf("tokenMasterEditionPubkey: %v\n", tokenMasterEditionPubkey.ToBase58())

	mintAccountRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		log.Fatalf("failed to get mint account rent, err: %v", err)
	}

	nonce, err := c.GetNonceFromNonceAccount(context.Background(), nonceAccountPubkey.String())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}
	fmt.Printf("nonce: %v\n", nonce)

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{mint, feePayer, alice},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: nonce,
			Instructions: []types.Instruction{
				sysprog.AdvanceNonceAccount(sysprog.AdvanceNonceAccountParam{
					Nonce: nonceAccountPubkey,
					Auth:  feePayer.PublicKey,
				}),
				sysprog.CreateAccount(sysprog.CreateAccountParam{
					From:     alice.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: mintAccountRent,
					Space:    tokenprog.MintAccountSize,
				}),
				tokenprog.InitializeMint(tokenprog.InitializeMintParam{
					Decimals: 0,
					Mint:     mint.PublicKey,
					MintAuth: alice.PublicKey,
				}),
				tokenmeta.CreateMetadataAccountV2(tokenmeta.CreateMetadataAccountV2Param{
					Metadata:                tokenMetadataPubkey,
					Mint:                    mint.PublicKey,
					MintAuthority:           alice.PublicKey,
					Payer:                   feePayer.PublicKey,
					UpdateAuthority:         alice.PublicKey,
					UpdateAuthorityIsSigner: true,
					IsMutable:               true,
					Data: tokenmeta.DataV2{
						Name:                 "Fake SMS #13055",
						Symbol:               "FSMB",
						Uri:                  "https://34c7ef24f4v2aejh75xhxy5z6ars4xv47gpsdrei6fiowptk2nqq.arweave.net/3wXyF1wvK6ARJ_9ue-O58CMuXrz5nyHEiPFQ6z5q02E",
						SellerFeeBasisPoints: 100,
						Creators: &[]tokenmeta.Creator{
							{
								Address:  alice.PublicKey,
								Verified: true,
								Share:    100,
							},
						},
						Collection: &tokenmeta.Collection{
							Verified: false,
							Key:      collection.PublicKey,
						},
						Uses: &tokenmeta.Uses{
							UseMethod: tokenmeta.Burn,
							Remaining: 10,
							Total:     10,
						},
					},
				}),
				assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
					Funder:                 feePayer.PublicKey,
					Owner:                  alice.PublicKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ata,
				}),
				tokenprog.MintTo(tokenprog.MintToParam{
					Mint:   mint.PublicKey,
					To:     ata,
					Auth:   alice.PublicKey,
					Amount: 1,
				}),
				tokenmeta.CreateMasterEditionV3(tokenmeta.CreateMasterEditionParam{
					Edition:         tokenMasterEditionPubkey,
					Mint:            mint.PublicKey,
					UpdateAuthority: alice.PublicKey,
					MintAuthority:   alice.PublicKey,
					Metadata:        tokenMetadataPubkey,
					Payer:           feePayer.PublicKey,
					MaxSupply:       pointer.Uint64(0),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a tx, err: %v", err)
	}

	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	fmt.Println(sig)
}
