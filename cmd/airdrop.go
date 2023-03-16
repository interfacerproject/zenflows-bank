package cmd

import (
	"context"
    "crypto/ecdsa"
	"log"
	"fmt"
	"github.com/interfacerproject/zenflows-bank/config"
	"github.com/interfacerproject/zenflows-bank/storage"
	"github.com/interfacerproject/zenflows-bank/fabcoin"
	"github.com/spf13/cobra"
	"math/big"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)


var minIdea int64
var minStrengths int64
var fabcoinsAmount int64

func init() {
	rootCmd.AddCommand(airdropCmd)
	airdropCmd.Flags().Int64Var(&minIdea, "minIdea", 10, "minimum amount of idea. Defaults to 10")
	airdropCmd.Flags().Int64Var(&minStrengths, "minStrengths", 10, "minimum amount of strengths. Defaults to 10")
	airdropCmd.Flags().Int64Var(&fabcoinsAmount, "fabcoinsAmount", 100, "fabcoinsAmount to transfer")
}

var airdropCmd = &cobra.Command{
	Use:   "airdrop",
	Short: "Make user earn fabcoins based on idea and strengths points",
	Run: func(cmd *cobra.Command, args []string) {
		balancesIO := storage.TokensFile{FileName: inputFile}
		balancesIO.Import()
		fmt.Println("Imported ", len(balancesIO.Tokens), "accounts")

		client, err := ethclient.Dial(config.Config.EthereumUrl)
		if err != nil {
			log.Fatal(err)
		}

		privateKey, err := crypto.HexToECDSA(config.Config.EthereumSk)
		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		chainID, err := client.ChainID(context.Background())
		if err != nil {
			panic(err)
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			panic(err)
		}
		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = big.NewInt(0)     // in wei
		auth.GasLimit = uint64(300000) // in units
		auth.GasPrice = gasPrice

		address := common.HexToAddress(config.Config.Fabcoin)
		instance, err := fabcoin.NewFabcoin(address, client)
		if err != nil {
			log.Fatal(err)
		}

		one := big.NewInt(1)

		fabcoinsReceipt := storage.FabcoinFile{FileName: outputFile}

		for _, v := range balancesIO.Tokens {
			var fabcoin int64 = 0
			if v.Idea > minIdea && v.Strengths > minStrengths {
				fabcoin = 100
			}

			if v.EthereumAddress == "" {
				continue
			}
			to := common.HexToAddress(v.EthereumAddress)

			txid := ""
			if fabcoin > 0 {
				amount := big.NewInt(int64(fabcoin))

				tx, err := instance.Transfer(auth, to, amount)
				if err != nil {
					log.Fatal(err)
				}
				txid = tx.Hash().Hex()
				auth.Nonce.Add(auth.Nonce, one)
			}

			fabcoinsReceipt.Fabcoins = append(fabcoinsReceipt.Fabcoins,
				&storage.Fabcoin{
					EthereumAddress: v.EthereumAddress,
					Idea: v.Idea,
					Strengths: v.Strengths,
					Fabcoin: fabcoin,
					TxId: txid,
				},
			)
		}
		fabcoinsReceipt.Export()
	},
}
