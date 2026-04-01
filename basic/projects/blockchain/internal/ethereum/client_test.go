package ethereum

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test with invalid RPC URL
	_, err := NewClient("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestWalletAddress(t *testing.T) {
	client := &Client{}

	wallet, err := NewWallet(client, "0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	if wallet.Address().Hex() == "" {
		t.Error("Wallet address should not be empty")
	}
}

func TestWalletInvalidKey(t *testing.T) {
	client := &Client{}

	_, err := NewWallet(client, "invalid-key")
	if err == nil {
		t.Error("Expected error for invalid private key")
	}
}
