package main

import (
	"fmt"
	//"github.com/mr-tron/base58"
	//"github.com/cosmos/go-bip39"
	//"crypto/ed25519"

	"github.com/portto/solana-go-sdk/types"
)

func main() {
	//entropy, _ := bip39.NewEntropy(256)
	//mnemonic, _ := bip39.NewMnemonic(entropy)
	//seed := bip39.NewSeed(mnemonic, "")
	//pri := ed25519.NewKeyFromSeed(seed[:32])
	//pub := base58.Encode(pri.Public().(ed25519.PublicKey))

	newAccount := types.NewAccount()
	fmt.Println(newAccount.PublicKey.ToBase58())
	fmt.Println(newAccount.PrivateKey)

	newAccount2 := types.AccountFromPrivateKeyBytes(newAccount.PrivateKey)
	fmt.Println(newAccount2.PublicKey.ToBase58())
}
