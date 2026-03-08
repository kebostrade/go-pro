# Building Blockchain Applications with Go and Ethereum

Create blockchain applications using Go and the Ethereum platform.

## Learning Objectives

- Understand Ethereum concepts
- Connect to Ethereum networks
- Interact with smart contracts
- Build transaction handling
- Implement wallet functionality
- Develop DApp backends

## Theory

### Connecting to Ethereum

```go
import (
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/common"
)

func connectEthereum(rpcURL string) (*ethclient.Client, error) {
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }

    networkID, err := client.NetworkID(context.Background())
    if err != nil {
        return nil, fmt.Errorf("failed to get network: %w", err)
    }

    log.Printf("Connected to network: %s", networkID)
    return client, nil
}

func getBlockInfo(client *ethclient.Client, blockNumber *big.Int) (*types.Block, error) {
    block, err := client.BlockByNumber(context.Background(), blockNumber)
    if err != nil {
        return nil, fmt.Errorf("failed to get block: %w", err)
    }

    fmt.Printf("Block #%d: %x\n", block.Number(), block.Hash())
    fmt.Printf("Transactions: %d\n", len(block.Transactions()))
    return block, nil
}
```

### Wallet Management

```go
import (
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/accounts/keystore"
)

func CreateWallet() (string, *ecdsa.PrivateKey, error) {
    privateKey, err := crypto.GenerateKey()
    if err != nil {
        return "", nil, fmt.Errorf("generate key: %w", err)
    }

    address := crypto.PubkeyToAddress(privateKey.PublicKey)
    return address.Hex(), privateKey, nil
}

func ImportWallet(privateKeyHex string) (string, *ecdsa.PrivateKey, error) {
    privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
    if err != nil {
        return "", nil, fmt.Errorf("import key: %w", err)
    }

    address := crypto.PubkeyToAddress(privateKey.PublicKey)
    return address.Hex(), privateKey, nil
}

func CreateKeystore(dir, password string) (string, error) {
    ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
    
    account, err := ks.NewAccount(password)
    if err != nil {
        return "", fmt.Errorf("create account: %w", err)
    }

    return account.Address.Hex(), nil
}
```

### Transaction Handling

```go
func getTransactor(client *ethclient.Client, privateKey *ecdsa.PrivateKey, chainID *big.Int) (*bind.TransactOpts, error) {
    fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
    
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        return nil, fmt.Errorf("get nonce: %w", err)
    }

    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        return nil, fmt.Errorf("suggest gas price: %w", err)
    }

    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
    if err != nil {
        return nil, fmt.Errorf("create transactor: %w", err)
    }

    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)
    auth.GasLimit = uint64(300000)
    auth.GasPrice = gasPrice

    return auth, nil
}

func sendTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, to common.Address, amount *big.Int) (common.Hash, error) {
    chainID, err := client.ChainID(context.Background())
    if err != nil {
        return common.Hash{}, err
    }

    fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        return common.Hash{}, err
    }

    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        return common.Hash{}, err
    }

    tx := types.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        return common.Hash{}, err
    }

    if err := client.SendTransaction(context.Background(), signedTx); err != nil {
        return common.Hash{}, err
    }

    return signedTx.Hash(), nil
}
```

### Smart Contract Interaction

```go
func deployContract(client *ethclient.Client, auth *bind.TransactOpts, bytecode string) (common.Address, *types.Transaction, error) {
    contractAddr, tx, _, err := bind.DeployContract(
        auth,
        common.FromHex(bytecode),
        nil,
        client,
    )
    if err != nil {
        return common.Address{}, nil, err
    }

    return contractAddr, tx, nil
}

func callContractMethod(client *ethclient.Client, contractAddr common.Address, methodID []byte, params []byte) ([]byte, error) {
    callMsg := ethereum.CallMsg{
        To:   &contractAddr,
        Data: append(methodID, params...),
    }

    result, err := client.CallContract(context.Background(), callMsg, nil)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func waitForTransaction(client *ethclient.Client, txHash common.Hash) error {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    for {
        receipt, err := client.TransactionReceipt(ctx, txHash)
        if err == nil {
            if receipt.Status == 0 {
                return errors.New("transaction failed")
            }
            return nil
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(2 * time.Second):
        }
    }
}
```

### Event Listening

```go
func subscribeToLogs(client *ethclient.Client, contractAddr common.Address) error {
    query := ethereum.FilterQuery{
        Addresses: []common.Address{contractAddr},
    }

    logs := make(chan types.Log)
    sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
    if err != nil {
        return err
    }
    defer sub.Unsubscribe()

    for {
        select {
        case err := <-sub.Err():
            return err
        case vLog := <-logs:
            log.Printf("Event: %+v", vLog)
            processEvent(vLog)
        }
    }
}

func processEvent(vLog types.Log) {
    transferSig := []byte("Transfer(address,address,uint256)")
    transferSigHash := crypto.Keccak256Hash(transferSig)

    if vLog.Topics[0] == transferSigHash {
        from := common.BytesToAddress(vLog.Topics[1].Bytes())
        to := common.BytesToAddress(vLog.Topics[2].Bytes())
        value := new(big.Int).SetBytes(vLog.Data)

        log.Printf("Transfer: from=%s to=%s value=%s", from, to, value)
    }
}
```

## Security Considerations

```go
func validateAddress(addr string) error {
    if !common.IsHexAddress(addr) {
        return errors.New("invalid address format")
    }
    return nil
}

func validateAmount(amount *big.Int) error {
    if amount == nil || amount.Sign() <= 0 {
        return errors.New("invalid amount")
    }
    if amount.Cmp(big.NewInt(1e18)) > 0 {
        return errors.New("amount too large")
    }
    return nil
}

func storePrivateKeySecurely(key *ecdsa.PrivateKey, path, password string) error {
    keyBytes := crypto.FromECDSA(key)
    encrypted, err := encrypt(keyBytes, password)
    if err != nil {
        return err
    }
    return os.WriteFile(path, encrypted, 0600)
}
```

## Performance Tips

```go
var clientPool = &sync.Pool{
    New: func() interface{} {
        client, err := ethclient.Dial(os.Getenv("RPC_URL"))
        if err != nil {
            return nil
        }
        return client
    },
}

func getBalancesBatch(addresses []string) ([]*big.Int, error) {
    client := clientPool.Get().(*ethclient.Client)
    defer clientPool.Put(client)

    var wg sync.WaitGroup
    results := make([]*big.Int, len(addresses))
    errs := make([]error, len(addresses))

    for i, addr := range addresses {
        wg.Add(1)
        go func(idx int, address string) {
            defer wg.Done()
            bal, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
            results[idx] = bal
            errs[idx] = err
        }(i, addr)
    }

    wg.Wait()

    for _, err := range errs {
        if err != nil {
            return nil, err
        }
    }

    return results, nil
}
```

## Exercises

1. Create a wallet service
2. Build a token transfer API
3. Listen to smart contract events
4. Implement transaction queue

## Validation

```bash
cd exercises
go test -v ./...
```

## Key Takeaways

- Always validate addresses and amounts
- Never store private keys in plaintext
- Handle transaction failures gracefully
- Use connection pooling for RPC
- Subscribe to events for real-time updates

## Next Steps

**[AT-13: IoT MQTT](../AT-13-iot-mqtt/README.md)**

---

Blockchain: trust through transparency. ⛓️
