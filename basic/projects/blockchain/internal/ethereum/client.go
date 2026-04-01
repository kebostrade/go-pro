// Package ethereum provides Ethereum client connection and interaction capabilities.
package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Client wraps the Ethereum client with additional functionality.
type Client struct {
	ethclient *ethclient.Client
	chainID   *big.Int
}

// NewClient creates a new Ethereum client connection.
func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return &Client{
		ethclient: client,
		chainID:   chainID,
	}, nil
}

// Close closes the client connection.
func (c *Client) Close() {
	c.ethclient.Close()
}

// GetBalance returns the ETH balance of an address.
func (c *Client) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	balance, err := c.ethclient.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

// GetBlockNumber returns the current block number.
func (c *Client) GetBlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := c.ethclient.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}
	return blockNumber, nil
}

// GetTransactionReceipt returns the receipt of a transaction.
func (c *Client) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := c.ethclient.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}
	return receipt, nil
}

// GetGasPrice returns the current gas price.
func (c *Client) GetGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := c.ethclient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	return gasPrice, nil
}

// ChainID returns the chain ID of the connected network.
func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	return c.chainID, nil
}

// NetworkID returns the network ID.
func (c *Client) NetworkID(ctx context.Context) (*big.Int, error) {
	networkID, err := c.ethclient.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}
	return networkID, nil
}

// GetCode returns the contract code at an address.
func (c *Client) GetCode(ctx context.Context, address common.Address) ([]byte, error) {
	code, err := c.ethclient.CodeAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get code: %w", err)
	}
	return code, nil
}

// CallContract calls a read-only contract method.
func (c *Client) CallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	result, err := c.ethclient.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}
	return result, nil
}

// PendingNonceAt returns the pending nonce for an account.
func (c *Client) PendingNonceAt(ctx context.Context, address common.Address) (uint64, error) {
	nonce, err := c.ethclient.PendingNonceAt(ctx, address)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending nonce: %w", err)
	}
	return nonce, nil
}

// SuggestGasPrice is alias for GetGasPrice.
func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.GetGasPrice(ctx)
}

// EstimateGas estimates the gas needed for a call.
func (c *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	gas, err := c.ethclient.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}
	return gas, nil
}
