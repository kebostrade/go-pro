package tensor

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestCreateTensor(t *testing.T) {
	tests := []struct {
		name string
		data []float64
		rows int
		cols int
	}{
		{
			name: "1x4 row vector",
			data: []float64{1, 2, 3, 4},
			rows: 1,
			cols: 4,
		},
		{
			name: "2x3 matrix",
			data: []float64{1, 2, 3, 4, 5, 6},
			rows: 2,
			cols: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateTensor(tt.data, tt.rows, tt.cols)
			r, c := result.Dims()
			if r != tt.rows || c != tt.cols {
				t.Errorf("CreateTensor() dims = %dx%d, want %dx%d", r, c, tt.rows, tt.cols)
			}
		})
	}
}

func TestZeros(t *testing.T) {
	z := Zeros(3, 4)

	r, c := z.Dims()
	if r != 3 || c != 4 {
		t.Errorf("Zeros() dims = %dx%d, want 3x4", r, c)
	}

	// Check all zeros
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if z.At(i, j) != 0 {
				t.Errorf("Zeros() contains non-zero value at (%d,%d): %v", i, j, z.At(i, j))
			}
		}
	}
}

func TestOnes(t *testing.T) {
	o := Ones(2, 3)

	r, c := o.Dims()
	if r != 2 || c != 3 {
		t.Errorf("Ones() dims = %dx%d, want 2x3", r, c)
	}

	// Check all ones
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if o.At(i, j) != 1 {
				t.Errorf("Ones() contains non-one value at (%d,%d): %v", i, j, o.At(i, j))
			}
		}
	}
}

func TestMatrixMultiply(t *testing.T) {
	// 2x3 matrix
	a := CreateTensor([]float64{1, 2, 3, 4, 5, 6}, 2, 3)
	// 3x2 matrix
	b := CreateTensor([]float64{7, 8, 9, 10, 11, 12}, 3, 2)

	result, err := MatrixMultiply(a, b)
	if err != nil {
		t.Fatalf("MatrixMultiply() error = %v", err)
	}

	r, c := result.Dims()
	if r != 2 || c != 2 {
		t.Errorf("MatrixMultiply() dims = %dx%d, want 2x2", r, c)
	}

	// Result should be [[58, 64], [139, 154]]
	expected := []float64{58, 64, 139, 154}
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if result.At(i, j) != expected[i*c+j] {
				t.Errorf("MatrixMultiply() at (%d,%d) = %v, want %v", i, j, result.At(i, j), expected[i*c+j])
			}
		}
	}
}

func TestAdd(t *testing.T) {
	a := CreateTensor([]float64{1, 2, 3, 4}, 2, 2)
	b := CreateTensor([]float64{5, 6, 7, 8}, 2, 2)

	result, err := Add(a, b)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	expected := []float64{6, 8, 10, 12}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if result.At(i, j) != expected[i*2+j] {
				t.Errorf("Add() at (%d,%d) = %v, want %v", i, j, result.At(i, j), expected[i*2+j])
			}
		}
	}
}

func TestSubtract(t *testing.T) {
	a := CreateTensor([]float64{10, 20, 30, 40}, 2, 2)
	b := CreateTensor([]float64{1, 2, 3, 4}, 2, 2)

	result, err := Subtract(a, b)
	if err != nil {
		t.Fatalf("Subtract() error = %v", err)
	}

	expected := []float64{9, 18, 27, 36}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if result.At(i, j) != expected[i*2+j] {
				t.Errorf("Subtract() at (%d,%d) = %v, want %v", i, j, result.At(i, j), expected[i*2+j])
			}
		}
	}
}

func TestReshape(t *testing.T) {
	original := CreateTensor([]float64{1, 2, 3, 4, 5, 6}, 2, 3)

	result := Reshape(original, 3, 2)

	r, c := result.Dims()
	if r != 3 || c != 2 {
		t.Errorf("Reshape() dims = %dx%d, want 3x2", r, c)
	}
}

func TestTranspose(t *testing.T) {
	m := CreateTensor([]float64{1, 2, 3, 4, 5, 6}, 2, 3)

	result := Transpose(m)

	r, c := result.Dims()
	if r != 3 || c != 2 {
		t.Errorf("Transpose() dims = %dx%d, want 3x2", r, c)
	}

	// Check values: [[1,2,3],[4,5,6]] transpose = [[1,4],[2,5],[3,6]]
	expected := []float64{1, 4, 2, 5, 3, 6}
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if result.At(i, j) != expected[i*c+j] {
				t.Errorf("Transpose() at (%d,%d) = %v, want %v", i, j, result.At(i, j), expected[i*c+j])
			}
		}
	}
}

func TestScale(t *testing.T) {
	m := CreateTensor([]float64{1, 2, 3, 4}, 2, 2)

	result := Scale(m, 2.0)

	expected := []float64{2, 4, 6, 8}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if result.At(i, j) != expected[i*2+j] {
				t.Errorf("Scale() at (%d,%d) = %v, want %v", i, j, result.At(i, j), expected[i*2+j])
			}
		}
	}
}

func TestSum(t *testing.T) {
	m := CreateTensor([]float64{1, 2, 3, 4, 5, 6}, 2, 3)

	sum := Sum(m)

	if sum != 21.0 {
		t.Errorf("Sum() = %v, want 21.0", sum)
	}
}

var _ = mat.Dense{} // Silence unused import warning
