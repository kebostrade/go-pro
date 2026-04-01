package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// StorageContract represents a simple key-value storage contract.
type StorageContract struct {
	Address common.Address
	client  *Client
}

// SetEvent represents a Set event from the contract.
type SetEvent struct {
	OldValue *big.Int
	NewValue *big.Int
}

// NewStorageContract creates a new storage contract instance.
func NewStorageContract(address common.Address, client *Client) *StorageContract {
	return &StorageContract{
		Address: address,
		client:  client,
	}
}

// Get retrieves the stored value.
func (s *StorageContract) Get(ctx context.Context) (*big.Int, error) {
	// Simplified - in production would use ABI binding
	data := []byte{
		0x6d, 0xa4, 0xd7, 0xdf, // keccak256("get()")
	}

	msg := ethereum.CallMsg{
		To:   &s.Address,
		Data: data,
	}

	result, err := s.client.CallContract(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to call get: %w", err)
	}

	// Parse result (simplified - assumes single uint256 return)
	if len(result) < 32 {
		return nil, fmt.Errorf("invalid result length")
	}

	return new(big.Int).SetBytes(result[:32]), nil
}

// Set sets a new value.
func (s *StorageContract) Set(ctx context.Context, value *big.Int) (common.Hash, error) {
	// This is a write operation - would need transaction signing
	// For demonstration, return error indicating need for wallet
	return common.Hash{}, fmt.Errorf("set requires wallet for transaction signing")
}

// DeployContract deploys a new storage contract.
func DeployContract(ctx context.Context, wallet *Wallet) (common.Address, common.Hash, error) {
	// In production, would use ABI binding or raw transaction
	// This is a placeholder demonstrating the pattern
	return common.Address{}, common.Hash{}, fmt.Errorf("deployment requires compiled contract and ABI")
}

// ParseABI parses an ABI JSON string.
func ParseABI(abiJSON string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(abiJSON))
}
