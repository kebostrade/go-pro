package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/discovery"
	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/logger"
	"github.com/DimaJoyti/go-pro/basic/projects/microservices-demo/pkg/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	serviceName = "product-service"
	servicePort = "8082"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

var products = make(map[string]*Product)

func main() {
	if err := logger.Init(serviceName, true); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Product Service", zap.String("port", servicePort))

	// Seed some products
	seedProducts()

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RateLimitMiddleware(100, 200))

	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/products", listProductsHandler).Methods("GET")
	r.HandleFunc("/products", createProductHandler).Methods("POST")
	r.HandleFunc("/products/{id}", getProductHandler).Methods("GET")
	r.HandleFunc("/products/{id}", updateProductHandler).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProductHandler).Methods("DELETE")

	serviceAddr := fmt.Sprintf("localhost:%s", servicePort)
	discovery.Register(serviceName, serviceAddr)
	defer discovery.Deregister(serviceName)

	srv := &http.Server{
		Addr:         ":" + servicePort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Product Service listening", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Product Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Product Service stopped")
}

func seedProducts() {
	products["1"] = &Product{
		ID:          "1",
		Name:        "Laptop",
		Description: "High-performance laptop",
		Price:       999.99,
		Stock:       10,
	}
	products["2"] = &Product{
		ID:          "2",
		Name:        "Mouse",
		Description: "Wireless mouse",
		Price:       29.99,
		Stock:       50,
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	productList := make([]*Product, 0, len(products))
	for _, p := range products {
		productList = append(productList, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": productList,
		"total":    len(productList),
	})
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ID = uuid.New().String()
	products[product.ID] = &product

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, ok := products[id]
	if !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, ok := products[id]
	if !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var updates Product
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updates.Name != "" {
		product.Name = updates.Name
	}
	if updates.Description != "" {
		product.Description = updates.Description
	}
	if updates.Price > 0 {
		product.Price = updates.Price
	}
	if updates.Stock >= 0 {
		product.Stock = updates.Stock
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, ok := products[id]; !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	delete(products, id)
	w.WriteHeader(http.StatusNoContent)
}

