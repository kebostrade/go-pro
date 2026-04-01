package main

import (
	"context"
	"fmt"
	"log"

	"github.com/goproject/blockchain/internal/ethereum"
)

func main() {
	log.Println("Blockchain Examples")
	fmt.Println("====================")

	rpcURL := "http://localhost:8545" // Default Ganache URL

	// Create client
	client, err := ethereum.NewClient(rpcURL)
	if err != nil {
		log.Printf("Note: Could not connect to %s (expected if not running)", rpcURL)
		log.Printf("Error: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()

	// Get block number
	blockNum, err := client.GetBlockNumber(ctx)
	if err != nil {
		log.Printf("Could not get block number: %v", err)
		return
	}
	fmt.Printf("Connected to network, current block: %d\n", blockNum)

	// Example: Working with storage contract
	fmt.Println("\nStorage Contract Example:")
	fmt.Println("-------------------------")
	fmt.Println("Contract address: 0x... (deploy to get actual address)")
	fmt.Println("Use wallet deploy command to deploy SimpleStorage contract")

	// Example: ERC-20 token interaction pattern
	fmt.Println("\nERC-20 Token Example:")
	fmt.Println("---------------------")
	fmt.Println("Token contract: 0x... (use deployed token address)")
	fmt.Println("Standard methods: transfer, balanceOf, approve, transferFrom")

	// Network info
	if chainID, err := client.ChainID(ctx); err == nil {
		fmt.Printf("\nNetwork Chain ID: %s\n", chainID.String())
	}
}
