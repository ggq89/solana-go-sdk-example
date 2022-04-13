package main

import (
	"context"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/portto/solana-go-sdk/rpc"
)

func main() {
	// NFT in solana is a normal mint but only mint 1.
	// If you want to get its metadata, you need to know where it stored.
	// and you can use `tokenmeta.GetTokenMetaPubkey` to get the metadata account key
	// here I take a random Degenerate Ape Academy as an example
	mint := common.PublicKeyFromString("6YSqRm9Zu272ZQqXyQnekaCNzxNcayRFzAKBe51Q6k7Z")
	metadataAccount, err := tokenmeta.GetTokenMetaPubkey(mint)
	if err != nil {
		log.Fatalf("faield to get metadata account, err: %v", err)
	}

	// new a client
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	// get data which stored in metadataAccount
	accountInfo, err := c.GetAccountInfo(context.Background(), metadataAccount.ToBase58())
	if err != nil {
		log.Fatalf("failed to get accountInfo, err: %v", err)
	}

	// parse it
	metadata, err := tokenmeta.MetadataDeserialize(accountInfo.Data)
	if err != nil {
		log.Fatalf("failed to parse metaAccount, err: %v", err)
	}

	log.Printf("metadata: %+v", metadata)
	log.Printf("Creators: %+v", metadata.Data.Creators)
	log.Printf("EditionNonce: %+v", *metadata.EditionNonce)
}
