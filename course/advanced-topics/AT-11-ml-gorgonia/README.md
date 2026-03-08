# Building Machine Learning Applications with Go and Gorgonia

Develop machine learning applications using Go and the Gorgonia library.

## Learning Objectives

- Understand tensor operations
- Build neural networks in Go
- Implement gradient descent
- Train and evaluate models
- Deploy ML models in production
- Optimize model performance

## Theory

### Basic Tensor Operations

```go
import (
    "gorgonia.org/tensor"
)

func basicTensors() {
    t1 := tensor.New(tensor.WithBacking([]float64{1, 2, 3, 4}), tensor.WithShape(2, 2))
    
    t2 := tensor.New(tensor.WithBacking([]float64{5, 6, 7, 8}), tensor.WithShape(2, 2))

    result, _ := tensor.Add(t1, t2)
    fmt.Println(result)

    result, _ = tensor.MatMul(t1, t2)
    fmt.Println(result)
}
```

### Building a Simple Neural Network

```go
import (
    "gorgonia.org/gorgonia"
    "gorgonia.org/tensor"
)

type NeuralNetwork struct {
    g      *gorgonia.ExprGraph
    w1     *gorgonia.Node
    w2     *gorgonia.Node
    pred   *gorgonia.Node
    loss   *gorgonia.Node
}

func NewNN(inputSize, hiddenSize, outputSize int) *NeuralNetwork {
    g := gorgonia.NewGraph()

    w1 := gorgonia.NewMatrix(g,
        tensor.Float64,
        gorgonia.WithShape(inputSize, hiddenSize),
        gorgonia.WithName("w1"),
        gorgonia.WithInit(gorgonia.GlorotN(1.0)),
    )

    w2 := gorgonia.NewMatrix(g,
        tensor.Float64,
        gorgonia.WithShape(hiddenSize, outputSize),
        gorgonia.WithName("w2"),
        gorgonia.WithInit(gorgonia.GlorotN(1.0)),
    )

    return &NeuralNetwork{g: g, w1: w1, w2: w2}
}

func (nn *NeuralNetwork) Forward(x *gorgonia.Node) (*gorgonia.Node, error) {
    hidden, err := gorgonia.Mul(x, nn.w1)
    if err != nil {
        return nil, err
    }
    hidden = gorgonia.Rectify(hidden)

    output, err := gorgonia.Mul(hidden, nn.w2)
    if err != nil {
        return nil, err
    }

    nn.pred = gorgonia.Sigmoid(output)
    return nn.pred, nil
}

func (nn *NeuralNetwork) Train(x, y tensor.Tensor, learningRate float64, epochs int) error {
    xNode := gorgonia.NewMatrix(nn.g, tensor.Float64, gorgonia.WithShape(x.Shape().Clone()...), gorgonia.WithName("x"))
    yNode := gorgonia.NewMatrix(nn.g, tensor.Float64, gorgonia.WithShape(y.Shape().Clone()...), gorgonia.WithName("y"))

    pred, err := nn.Forward(xNode)
    if err != nil {
        return err
    }

    loss, err := gorgonia.BinaryXent(pred, yNode)
    if err != nil {
        return err
    }
    nn.loss = loss

    _, err = gorgonia.Grad(loss, nn.w1, nn.w2)
    if err != nil {
        return err
    }

    vm := gorgonia.NewTapeMachine(nn.g)
    defer vm.Close()

    if err := gorgonia.Let(xNode, x); err != nil {
        return err
    }
    if err := gorgonia.Let(yNode, y); err != nil {
        return err
    }

    for epoch := 0; epoch < epochs; epoch++ {
        if err := vm.RunAll(); err != nil {
            return err
        }

        lossVal := nn.loss.Value().Data().(float64)
        fmt.Printf("Epoch %d: Loss = %.4f\n", epoch, lossVal)

        vm.Reset()
    }

    return nil
}
```

### Linear Regression Example

```go
func linearRegression() {
    g := gorgonia.NewGraph()

    x := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(100, 1), gorgonia.WithName("x"))
    y := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(100, 1), gorgonia.WithName("y"))

    w := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(1, 1), gorgonia.WithName("w"), gorgonia.WithInit(gorgonia.Zeroes()))
    b := gorgonia.NewScalar(g, tensor.Float64, gorgonia.WithName("b"))

    pred, _ := gorgonia.Add(gorgonia.Must(gorgonia.Mul(x, w)), b)

    loss, _ := gorgonia.Mean(gorgonia.Must(gorgonia.Square(gorgonia.Must(gorgonia.Sub(pred, y)))))

    _, _ = gorgonia.Grad(loss, w, b)

    vm := gorgonia.NewTapeMachine(g)
    defer vm.Close()

    xData := generateLinearData(100)
    yData := generateLabels(xData)

    gorgonia.Let(x, xData)
    gorgonia.Let(y, yData)

    for i := 0; i < 1000; i++ {
        vm.RunAll()
        vm.Reset()
    }

    fmt.Printf("Trained weight: %v\n", w.Value())
    fmt.Printf("Trained bias: %v\n", b.Value())
}
```

### Model Serialization

```go
func (nn *NeuralNetwork) Save(path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    enc := gob.NewEncoder(f)
    return enc.Encode(map[string]tensor.Tensor{
        "w1": nn.w1.Value().(tensor.Tensor),
        "w2": nn.w2.Value().(tensor.Tensor),
    })
}

func (nn *NeuralNetwork) Load(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()

    var weights map[string]tensor.Tensor
    dec := gob.NewDecoder(f)
    if err := dec.Decode(&weights); err != nil {
        return err
    }

    gorgonia.Let(nn.w1, weights["w1"])
    gorgonia.Let(nn.w2, weights["w2"])
    return nil
}
```

## Security Considerations

```go
func validateInput(input []float64, expectedLen int) error {
    if len(input) != expectedLen {
        return fmt.Errorf("expected %d values, got %d", expectedLen, len(input))
    }
    for i, v := range input {
        if math.IsNaN(v) || math.IsInf(v, 0) {
            return fmt.Errorf("invalid value at index %d", i)
        }
    }
    return nil
}

func sanitizeModelPath(path string) error {
    clean := filepath.Clean(path)
    if strings.Contains(clean, "..") {
        return errors.New("invalid path")
    }
    return nil
}
```

## Performance Tips

```go
func batchPredict(model *NeuralNetwork, inputs tensor.Tensor, batchSize int) (tensor.Tensor, error) {
    shape := inputs.Shape()
    n := shape[0]
    outputShape := []int{n, model.outputSize}
    results := tensor.New(tensor.WithBacking(make([]float64, n*model.outputSize)), tensor.WithShape(outputShape...))

    for i := 0; i < n; i += batchSize {
        end := i + batchSize
        if end > n {
            end = n
        }

        batch := inputs.(tensor.Slicer).Slice(gorgonia.S(i), gorgonia.S(end))
        pred, err := model.Predict(batch)
        if err != nil {
            return nil, err
        }

        results.(tensor.Slicer).Slice(gorgonia.S(i), gorgonia.S(end)).(*tensor.Dense).Copy(pred)
    }

    return results, nil
}
```

## Exercises

1. Implement logistic regression
2. Build a CNN for image classification
3. Create a recommendation system
4. Deploy model as REST API

## Validation

```bash
cd exercises
go test -v ./...
go run examples/linear_regression.go
```

## Key Takeaways

- Tensors are fundamental building blocks
- Build computational graphs for models
- Use gradient descent for training
- Validate inputs before prediction
- Serialize trained models

## Next Steps

**[AT-12: Blockchain Ethereum](../AT-12-blockchain-ethereum/README.md)**

---

ML in Go: fast, compiled, deployable. 🧠
