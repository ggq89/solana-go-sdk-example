package main

import (
	"fmt"
	"github.com/cosmos/go-bip39"
	"golang.org/x/crypto/ed25519"

	"github.com/portto/solana-go-sdk/types"
)

func main() {
	newAccount := types.NewAccount()
	fmt.Println(newAccount.PublicKey.ToBase58())
	fmt.Println(newAccount.PrivateKey)

	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println(mnemonic)
	seed := bip39.NewSeed(mnemonic, "")
	pri := ed25519.NewKeyFromSeed(seed[:32])
	//newAccount2 := types.AccountFromPrivateKeyBytes(newAccount.PrivateKey)
	newAccount2 := types.AccountFromPrivateKeyBytes(pri)
	fmt.Println(newAccount2.PublicKey.ToBase58())
}
