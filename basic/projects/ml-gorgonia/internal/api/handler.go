// Package api provides HTTP handlers for the ML inference server.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/goproject/ml-gorgonia/internal/model"
	"github.com/goproject/ml-gorgonia/internal/tensor"

	"github.com/go-chi/chi/v5"
	"gonum.org/v1/gonum/mat"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	model model.Model
}

// NewHandler creates a new API handler.
func NewHandler(m model.Model) *Handler {
	return &Handler{model: m}
}

// TensorOpRequest represents a tensor operation request.
type TensorOpRequest struct {
	Operation string    `json:"operation"`
	Data      []float64 `json:"data"`
	Rows      int       `json:"rows"`
	Cols      int       `json:"cols"`
	Operand   []float64 `json:"operand,omitempty"`
}

// TensorOpResponse represents a tensor operation response.
type TensorOpResponse struct {
	Result []float64 `json:"result"`
	Rows   int       `json:"rows"`
	Cols   int       `json:"cols"`
}

// InferenceRequest represents an inference request.
type InferenceRequest struct {
	Input []float64 `json:"input"`
	Rows  int       `json:"rows"`
	Cols  int       `json:"cols"`
}

// InferenceResponse represents an inference response.
type InferenceResponse struct {
	Output []float64 `json:"output"`
}

// TensorOp handles tensor operations.
func (h *Handler) TensorOp(w http.ResponseWriter, r *http.Request) {
	var req TensorOpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t := tensor.CreateTensor(req.Data, req.Rows, req.Cols)

	var result *mat.Dense
	switch req.Operation {
	case "add":
		if len(req.Operand) == 0 {
			http.Error(w, "operand required for add", http.StatusBadRequest)
			return
		}
		operand := tensor.CreateTensor(req.Operand, req.Rows, req.Cols)
		result, _ = tensor.Add(t, operand)
	case "subtract":
		if len(req.Operand) == 0 {
			http.Error(w, "operand required for subtract", http.StatusBadRequest)
			return
		}
		operand := tensor.CreateTensor(req.Operand, req.Rows, req.Cols)
		result, _ = tensor.Subtract(t, operand)
	case "transpose":
		result = tensor.Transpose(t)
	case "reshape":
		// Reshape to 1D
		result = tensor.Reshape(t, 1, req.Rows*req.Cols)
	default:
		http.Error(w, "unknown operation", http.StatusBadRequest)
		return
	}

	// Extract data
	rows, cols := result.Dims()
	data := make([]float64, rows*cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			data[i*cols+j] = result.At(i, j)
		}
	}

	resp := TensorOpResponse{
		Result: data,
		Rows:   rows,
		Cols:   cols,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Inference handles model inference requests.
func (h *Handler) Inference(w http.ResponseWriter, r *http.Request) {
	var req InferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := model.RunInference(h.model, req.Input, req.Rows, req.Cols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := InferenceResponse{Output: output}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Health handles health check requests.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Routes returns the chi router with all routes configured.
func (h *Handler) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/v1/health", h.Health)
	r.Post("/api/v1/tensor", h.TensorOp)
	r.Post("/api/v1/inference", h.Inference)

	return r
}
