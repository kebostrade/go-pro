package internal

import (
	"encoding/json"
	"net/http"

	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/logger"
	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Handler handles HTTP requests for the user service
type Handler struct {
	repo Repository
}

// NewHandler creates a new handler
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// RegisterRoutes registers all routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/users", h.CreateUser).Methods("POST")
	r.HandleFunc("/users/login", h.Login).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/users").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("", h.ListUsers).Methods("GET")
	protected.HandleFunc("/{id}", h.GetUser).Methods("GET")
	protected.HandleFunc("/{id}", h.UpdateUser).Methods("PUT")
	protected.HandleFunc("/{id}", h.DeleteUser).Methods("DELETE")
}

// CreateUser creates a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.repo.Create(user); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("User created", zap.String("user_id", user.ID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login authenticates a user
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUser retrieves a user by ID
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a user
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := h.repo.Update(id, user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, _ := h.repo.GetByID(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListUsers lists all users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10

	users, total, err := h.repo.List(page, pageSize)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

