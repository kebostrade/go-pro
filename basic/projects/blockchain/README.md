# Blockchain-Ethereum: Smart Contracts and Wallet Operations

Production-ready Ethereum blockchain template using go-ethereum for wallet operations and smart contract interactions.

## Features

- **Ethereum Client**: Connect to Ethereum networks (mainnet, testnets, local)
- **Wallet Operations**: Send ETH, sign transactions, manage keys
- **Smart Contract Interaction**: Call contracts, listen to events
- **Docker Support**: Local development with Ganache

## Installation

```bash
# Clone the repository
git clone https://github.com/goproject/blockchain.git
cd blockchain

# Download dependencies
go mod download

# Build the wallet tool
go build ./cmd/wallet
```

## Usage

### CLI Wallet

```bash
# Check balance
PRIVATE_KEY=0x... ETH_RPC_URL=https://sepolia.infura.io/v3/YOUR_PROJECT_ID ./wallet balance

# Send ETH
PRIVATE_KEY=0x... ./wallet send 0x... 1000000000000000000
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `ETH_RPC_URL` | Ethereum RPC endpoint | Yes |
| `PRIVATE_KEY` | Wallet private key (hex) | Yes |
| `CONTRACT_ADDRESS` | Deployed contract address | No |

## Docker Development

```bash
# Start Ganache (local blockchain)
docker-compose up ganache -d

# Run wallet commands
docker-compose run wallet balance
```

## Smart Contracts

### SimpleStorage

A basic key-value storage contract:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleStorage {
    uint256 private value;

    event ValueChanged(uint256 oldValue, uint256 newValue);

    function store(uint256 newValue) public {
        emit ValueChanged(value, newValue);
        value = newValue;
    }

    function retrieve() public view returns (uint256) {
        return value;
    }
}
```

## Project Structure

```
blockchain/
├── cmd/wallet/main.go        # CLI wallet tool
├── internal/
│   └── ethereum/
│       ├── client.go        # Ethereum client
│       ├── wallet.go         # Wallet operations
│       └── contracts.go     # Contract interactions
├── abi/                      # Contract ABIs
├── examples/                 # Usage examples
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

## License

MIT
