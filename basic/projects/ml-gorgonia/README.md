# ML-Gorgonia: Tensor Operations and ONNX Inference

Production-ready ML inference template using Gorgonia for tensor operations and ONNX model inference in Go.

## Features

- **Tensor Operations**: Create, multiply, reshape, transpose tensors using Gorgonia
- **ONNX Inference**: Load and run ONNX models for ML inference
- **HTTP API**: REST API server for tensor operations and model inference
- **Docker Support**: Containerized deployment with Docker and Docker Compose

## Installation

```bash
# Clone the repository
git clone https://github.com/goproject/ml-gorgonia.git
cd ml-gorgonia

# Download dependencies
go mod download

# Build the server
go build ./cmd/server
```

## Usage

### HTTP API Server

```bash
# Start the server
./server

# Or with Docker
docker-compose up -d
```

### API Endpoints

- `GET /api/v1/health` - Health check
- `POST /api/v1/tensor` - Perform tensor operations
- `POST /api/v1/inference` - Run model inference

### Tensor Operations Example

```bash
curl -X POST http://localhost:8080/api/v1/tensor \
  -H "Content-Type: application/json" \
  -d '{
    "operation": "add",
    "data": [1, 2, 3, 4],
    "shape": [2, 2],
    "operand": [5, 6, 7, 8]
  }'
```

### Model Inference Example

```bash
curl -X POST http://localhost:8080/api/v1/inference \
  -H "Content-Type: application/json" \
  -d '{
    "input": [0.5, 0.3, 0.8, 0.1],
    "shape": [2, 2]
  }'
```

## Tensor Operations

The template provides the following tensor operations:

| Operation | Description |
|------------|-------------|
| `CreateTensor` | Create tensor from flat slice with shape |
| `Zeros` | Create tensor filled with zeros |
| `Ones` | Create tensor filled with ones |
| `MatrixMultiply` | Matrix multiplication |
| `Add` | Element-wise addition |
| `Subtract` | Element-wise subtraction |
| `Reshape` | Change tensor shape |
| `Transpose` | Transpose a 2D tensor |
| `Scale` | Multiply by scalar |
| `Sum` | Sum all elements |

## Running Examples

```bash
# Run MNIST example
go run ./examples/mnist.go
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Project Structure

```
ml-gorgonia/
├── cmd/server/main.go        # HTTP inference server
├── internal/
│   ├── tensor/              # Tensor operations
│   ├── model/               # ONNX inference
│   └── api/                 # HTTP handlers
├── examples/                 # Example usage
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `MODEL_PATH` | Path to ONNX model file | (none, uses mock) |

## License

MIT
