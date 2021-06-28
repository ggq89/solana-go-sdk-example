package main

import (
	"context"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/client"
)

func main() {
	c := client.NewClient("https://api.devnet.solana.com")
	balance, err := c.GetBalance(context.Background(), "83R5RVHMEEmHtj9QydfAX958JDoNHKREmQhw8k24ryMj")
	if err != nil {
		log.Fatalln("get balance error", err)
	}
	fmt.Println(balance)
}
