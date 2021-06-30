package main

import (
	"context"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/client"
)

func main() {
	c := client.NewClient("https://api.devnet.solana.com")
	balance, err := c.GetTokenAccountBalance(context.Background(),
		"AKpvREJ4JpdtHYWuk4W5J2Sro1ZfN11q7kuN7omGRCAC", client.CommitmentFinalized)
	if err != nil {
		log.Fatalln("get spl balance error", err)
	}
	fmt.Println(balance.UIAmountString)
}
