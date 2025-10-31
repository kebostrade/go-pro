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
	serviceName = "order-service"
	servicePort = "8083"
)

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var orders = make(map[string]*Order)

func main() {
	if err := logger.Init(serviceName, true); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Order Service", zap.String("port", servicePort))

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RateLimitMiddleware(100, 200))

	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/orders", listOrdersHandler).Methods("GET")
	r.HandleFunc("/orders", createOrderHandler).Methods("POST")
	r.HandleFunc("/orders/{id}", getOrderHandler).Methods("GET")
	r.HandleFunc("/orders/{id}/status", updateOrderStatusHandler).Methods("PUT")

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
		logger.Info("Order Service listening", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Order Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Order Service stopped")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func listOrdersHandler(w http.ResponseWriter, r *http.Request) {
	orderList := make([]*Order, 0, len(orders))
	for _, o := range orders {
		orderList = append(orderList, o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"orders": orderList,
		"total":  len(orderList),
	})
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order.ID = uuid.New().String()
	order.Status = "pending"
	order.CreatedAt = time.Now()
	orders[order.ID] = &order

	logger.Info("Order created", zap.String("order_id", order.ID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	order, ok := orders[id]
	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func updateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	order, ok := orders[id]
	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order.Status = req.Status
	logger.Info("Order status updated", zap.String("order_id", id), zap.String("status", req.Status))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

