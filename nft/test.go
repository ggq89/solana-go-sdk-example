package main

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"log"
)

func main() {
	c := client.NewClient(rpc.TestnetRPCEndpoint)

	//4Gr5C8ryhMCZeHiKJ99jtb5SRcCofa55mCidVtvX6HrD4PTYi5eaxMu3hcmkDS2uBsCKDv461Qu9LR8jCrjDXPLz
	//2CVNNiGkPy7Zgs54bi9RN3MvRUjF7enY7KRzJpJxooNNtZ37TdrfvK7EQZQu5ifcfjAr4oJJiiA9evTXPTPkVgnH
	//3Eyq6c9okX3i4mRxatdwYasFRhfwK1SEqDVPUaMyTsPWPcAQfUWmqxdVsSaiqs9hfmcKbjVyasSBF2k8LbTw6F8N
	result, err := c.GetSignatureStatusWithConfig(context.Background(),
	"3Eyq6c9okX3i4mRxatdwYasFRhfwK1SEqDVPUaMyTsPWPcAQfUWmqxdVsSaiqs9hfmcKbjVyasSBF2k8LbTw6F8N",
		rpc.GetSignatureStatusesConfig{SearchTransactionHistory: true})

	if err != nil {
		log.Println(err.Error())
		return
	}

	if result == nil {
		log.Println("tx not existed")
		return
	}

	spew.Dump(result)

	if result.Confirmations == nil {
		if result.Err != nil {
			log.Println("tx failed")
			return
		}

		log.Println("tx finalized")
		return
	}
}
