package model

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestLoadONNXModel(t *testing.T) {
	// Test loading non-existent file
	_, err := LoadONNXModel("/nonexistent/model.onnx")
	if err == nil {
		t.Error("LoadONNXModel() expected error for non-existent file")
	}
}

func TestMockModel(t *testing.T) {
	m := MockModel()

	input := mat.NewDense(2, 3, []float64{1, 2, 3, 4, 5, 6})

	output, err := m.Run(input)
	if err != nil {
		t.Fatalf("MockModel.Run() error = %v", err)
	}

	r, c := output.Dims()
	if r != 2 || c != 3 {
		t.Errorf("MockModel.Run() dims = %dx%d, want 2x3", r, c)
	}
}

func TestRunInference(t *testing.T) {
	model := MockModel()

	input := []float64{1, 2, 3, 4, 5, 6}
	rows, cols := 2, 3

	result, err := RunInference(model, input, rows, cols)
	if err != nil {
		t.Fatalf("RunInference() error = %v", err)
	}

	if len(result) != len(input) {
		t.Errorf("RunInference() returned %d values, want %d", len(result), len(input))
	}

	for i, v := range result {
		if v != input[i] {
			t.Errorf("RunInference() value at %d = %v, want %v", i, v, input[i])
		}
	}
}

func TestPreprocessImage(t *testing.T) {
	img := []byte{255, 128, 0, 64}

	result := PreprocessImage(img, 2, 2)

	r, c := result.Dims()
	if r != 2 || c != 2 {
		t.Errorf("PreprocessImage() dims = %dx%d, want 2x2", r, c)
	}

	// Check normalization (use tolerance for floating point)
	expected := []float64{1.0, 0.502, 0.0, 0.251}
	data := make([]float64, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			data[i*c+j] = result.At(i, j)
		}
	}
	for i, v := range data {
		if abs(v-expected[i]) > 0.001 {
			t.Errorf("PreprocessImage() value at %d = %v, want ~%v", i, v, expected[i])
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestPostprocessOutput(t *testing.T) {
	output := mat.NewDense(1, 3, []float64{0.1, 0.7, 0.2})

	result, err := PostprocessOutput(output)
	if err != nil {
		t.Fatalf("PostprocessOutput() error = %v", err)
	}

	if len(result) != 3 {
		t.Errorf("PostprocessOutput() returned %d values, want 3", len(result))
	}
}
