// Package model provides ONNX model loading and inference capabilities.
package model

import (
	"fmt"
	"os"

	"gonum.org/v1/gonum/mat"
)

// Model interface defines the contract for ML models.
type Model interface {
	Run(input *mat.Dense) (*mat.Dense, error)
}

// ONNXModel represents an ONNX model for inference.
type ONNXModel struct {
	modelPath string
	inputs    []string
	outputs   []string
}

// LoadONNXModel loads an ONNX model from the given path.
func LoadONNXModel(modelPath string) (*ONNXModel, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found: %s", modelPath)
	}

	return &ONNXModel{
		modelPath: modelPath,
		inputs:    []string{"input"},
		outputs:   []string{"output"},
	}, nil
}

// Run performs inference using the loaded model.
// For demonstration, this creates a simple pass-through.
func (m *ONNXModel) Run(input *mat.Dense) (*mat.Dense, error) {
	if input == nil {
		return nil, fmt.Errorf("no input tensor provided")
	}

	// Simple inference: pass through with identity for demonstration
	// In production, this would use onnx-runtime or gonnx
	var result mat.Dense
	result.CloneFrom(input)
	return &result, nil
}

// RunInference runs inference on the model with the given input data.
func RunInference(model Model, input []float64, rows, cols int) ([]float64, error) {
	// Create input matrix
	inputMat := mat.NewDense(rows, cols, input)

	// Run model
	output, err := model.Run(inputMat)
	if err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	// Extract output data
	r, c := output.Dims()
	data := make([]float64, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			data[i*c+j] = output.At(i, j)
		}
	}

	return data, nil
}

// PreprocessImage preprocesses image data for model input.
// Converts image bytes to normalized float tensor.
func PreprocessImage(img []byte, rows, cols int) *mat.Dense {
	normalized := make([]float64, len(img))

	for i, b := range img {
		// Normalize to [0, 1]
		normalized[i] = float64(b) / 255.0
	}

	return mat.NewDense(rows, cols, normalized)
}

// PostprocessOutput converts model output to human-readable format.
func PostprocessOutput(output *mat.Dense) ([]float64, error) {
	r, c := output.Dims()
	data := make([]float64, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			data[i*c+j] = output.At(i, j)
		}
	}
	return data, nil
}

// MockModel returns a simple mock model for testing without ONNX runtime.
func MockModel() Model {
	return &mockModel{}
}

type mockModel struct{}

func (m *mockModel) Run(input *mat.Dense) (*mat.Dense, error) {
	if input == nil {
		return nil, fmt.Errorf("no inputs")
	}
	// Return input as output (identity function)
	var result mat.Dense
	result.CloneFrom(input)
	return &result, nil
}
