package main

import (
	"fmt"
	"log"

	"github.com/goproject/ml-gorgonia/internal/model"
	"github.com/goproject/ml-gorgonia/internal/tensor"
)

func main() {
	log.Println("ML-Gorgonia Examples")

	// Example 1: Create and manipulate tensors
	fmt.Println("\n=== Tensor Operations ===")

	// Create a 2x3 matrix
	data := []float64{1, 2, 3, 4, 5, 6}

	matrix := tensor.CreateTensor(data, 2, 3)
	fmt.Printf("Created matrix:\n%v\n", matrix)

	// Transpose
	transposed := tensor.Transpose(matrix)
	fmt.Printf("Transposed:\n%v\n", transposed)

	// Matrix multiplication
	a := tensor.CreateTensor([]float64{1, 2, 3, 4}, 2, 2)
	b := tensor.CreateTensor([]float64{5, 6, 7, 8}, 2, 2)

	product, err := tensor.MatrixMultiply(a, b)
	if err != nil {
		log.Fatalf("Failed to multiply: %v", err)
	}
	fmt.Printf("Matrix product:\n%v\n", product)

	// Example 2: Model inference
	fmt.Println("\n=== Model Inference ===")

	// Create mock model
	m := model.MockModel()

	// Run inference
	input := []float64{0.5, 0.3, 0.8, 0.1}

	output, err := model.RunInference(m, input, 2, 2)
	if err != nil {
		log.Fatalf("Inference failed: %v", err)
	}

	fmt.Printf("Input: %v\n", input)
	fmt.Printf("Output: %v\n", output)

	fmt.Println("\nExamples completed successfully!")
}
