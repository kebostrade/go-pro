package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goproject/blockchain/internal/ethereum"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	log.Println("Ethereum Wallet CLI")

	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://localhost:8545"
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("PRIVATE_KEY environment variable is required")
	}

	// Create client
	client, err := ethereum.NewClient(rpcURL)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create wallet
	wallet, err := ethereum.NewWallet(client, privateKey)
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	// Parse command
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "balance":
		balance, err := wallet.Balance(context.Background())
		if err != nil {
			log.Fatalf("Failed to get balance: %v", err)
		}
		fmt.Printf("Balance: %s ETH\n", balance.Div(balance, big.NewInt(1e18)))

	case "send":
		if len(os.Args) < 4 {
			fmt.Println("Usage: wallet send <to_address> <amount_wei>")
			os.Exit(1)
		}
		toAddr := common.HexToAddress(os.Args[2])
		amount := new(big.Int)
		amount.SetString(os.Args[3], 10)

		txHash, err := wallet.SendETH(context.Background(), toAddr, amount)
		if err != nil {
			log.Fatalf("Failed to send: %v", err)
		}
		fmt.Printf("Transaction sent: %s\n", txHash.Hex())

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  wallet balance              - Check ETH balance")
	fmt.Println("  wallet send <address> <wei> - Send ETH")
}
