//go:build ignore
// +build ignore

// This is a standalone example - run with: go run order.go
// Note: Rename mainOrder to main to run this file standalone

package main

import (
	"math/rand"
	"sort"
)

type Order struct {
	ID     int
	Status string
	Price  float64
}

func generateOrders(count int) []*Order {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = &Order{
			ID:     i + 1,
			Status: "pending",
			Price:  rand.Float64() * 1000,
		}
	}
	return orders
}

// Sort orders by price in ascending order
func sortOrders(orders []*Order) {
	// Sort orders using built-in sort package
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Price < orders[j].Price
	})
}
