package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Wallet handles wallet operations and transactions.
type Wallet struct {
	client   *Client
	key      *ecdsa.PrivateKey
	keystore *keystore.KeyStore
	fromAddr common.Address
}

// NewWallet creates a new wallet from a private key.
func NewWallet(client *Client, privateKeyHex string) (*Wallet, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}

	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Wallet{
		client:   client,
		key:      privateKey,
		fromAddr: fromAddr,
	}, nil
}

// NewWalletFromKeystore creates a wallet from a keystore file.
func NewWalletFromKeystore(client *Client, ks *keystore.KeyStore, account accounts.Account) (*Wallet, error) {
	return &Wallet{
		client:   client,
		keystore: ks,
		fromAddr: account.Address,
	}, nil
}

// Address returns the wallet's address.
func (w *Wallet) Address() common.Address {
	return w.fromAddr
}

// SendETH sends ETH to an address.
func (w *Wallet) SendETH(ctx context.Context, to common.Address, amount *big.Int) (common.Hash, error) {
	// Get nonce
	nonce, err := w.client.PendingNonceAt(ctx, w.fromAddr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price
	gasPrice, err := w.client.GetGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Estimate gas
	msg := ethereum.CallMsg{
		From:  w.fromAddr,
		To:    &to,
		Value: amount,
	}
	gasLimit, err := w.client.EstimateGas(ctx, msg)
	if err != nil {
		gasLimit = 21000 // Default gas limit for ETH transfer
	}

	// Create transaction
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)

	// Sign transaction
	signedTx, err := w.SignTransaction(tx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = w.client.ethclient.SendTransaction(ctx, signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash(), nil
}

// SignTransaction signs a transaction with the wallet's private key.
func (w *Wallet) SignTransaction(tx *types.Transaction) (*types.Transaction, error) {
	if w.key == nil {
		return nil, fmt.Errorf("no private key available")
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(w.client.chainID), w.key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return signedTx, nil
}

// Balance returns the wallet's ETH balance.
func (w *Wallet) Balance(ctx context.Context) (*big.Int, error) {
	return w.client.GetBalance(ctx, w.fromAddr)
}

// TransferToken transfers ERC-20 tokens.
func (w *Wallet) TransferToken(ctx context.Context, tokenAddr, to common.Address, amount *big.Int) (common.Hash, error) {
	// This would require the token ABI and contract interaction
	// For now, return an error indicating this needs implementation
	return common.Hash{}, fmt.Errorf("token transfer not implemented - use contracts package")
}
