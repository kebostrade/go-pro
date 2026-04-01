// Package tensor provides tensor operations using Gonum.
package tensor

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// CreateTensor creates a new dense tensor from a flat slice with the given shape.
func CreateTensor(data []float64, shape ...int) *mat.Dense {
	if len(shape) == 0 {
		shape = []int{1, len(data)}
	}

	rows, cols := shape[0], shape[1]
	if len(shape) == 1 {
		rows = 1
		cols = len(data)
	}

	// Pad with zeros if needed
	if len(data) < rows*cols {
		padded := make([]float64, rows*cols)
		copy(padded, data)
		data = padded
	}

	return mat.NewDense(rows, cols, data)
}

// Zeros creates a matrix filled with zeros.
func Zeros(rows, cols int) *mat.Dense {
	return mat.NewDense(rows, cols, make([]float64, rows*cols))
}

// Ones creates a matrix filled with ones.
func Ones(rows, cols int) *mat.Dense {
	data := make([]float64, rows*cols)
	for i := range data {
		data[i] = 1
	}
	return mat.NewDense(rows, cols, data)
}

// MatrixMultiply performs matrix multiplication of two matrices.
func MatrixMultiply(a, b *mat.Dense) (*mat.Dense, error) {
	var result mat.Dense
	result.Mul(a, b)
	return &result, nil
}

// Add performs element-wise addition of two matrices.
func Add(a, b *mat.Dense) (*mat.Dense, error) {
	rows, cols := a.Dims()
	resultRows, resultCols := b.Dims()
	if rows != resultRows || cols != resultCols {
		return nil, fmt.Errorf("dimension mismatch: %dx%d + %dx%d", rows, cols, resultRows, resultCols)
	}

	result := mat.NewDense(rows, cols, make([]float64, rows*cols))
	result.Add(a, b)
	return result, nil
}

// Subtract performs element-wise subtraction of two matrices.
func Subtract(a, b *mat.Dense) (*mat.Dense, error) {
	rows, cols := a.Dims()
	resultRows, resultCols := b.Dims()
	if rows != resultRows || cols != resultCols {
		return nil, fmt.Errorf("dimension mismatch: %dx%d - %dx%d", rows, cols, resultRows, resultCols)
	}

	result := mat.NewDense(rows, cols, make([]float64, rows*cols))
	result.Sub(a, b)
	return result, nil
}

// Reshape reshapes a matrix to a new shape.
func Reshape(t *mat.Dense, rows, cols int) *mat.Dense {
	r, c := t.Dims()
	data := make([]float64, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			data[i*c+j] = t.At(i, j)
		}
	}
	return mat.NewDense(rows, cols, data)
}

// Transpose transposes a matrix.
func Transpose(t *mat.Dense) *mat.Dense {
	rows, cols := t.Dims()
	result := mat.NewDense(cols, rows, make([]float64, rows*cols))
	result.CloneFrom(t.T())
	return result
}

// Scale multiplies all elements by a scalar.
func Scale(t *mat.Dense, scalar float64) *mat.Dense {
	rows, cols := t.Dims()
	result := mat.NewDense(rows, cols, make([]float64, rows*cols))
	result.Scale(scalar, t)
	return result
}

// Sum calculates the sum of all elements.
func Sum(t *mat.Dense) float64 {
	rows, cols := t.Dims()
	total := 0.0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			total += t.At(i, j)
		}
	}
	return total
}
