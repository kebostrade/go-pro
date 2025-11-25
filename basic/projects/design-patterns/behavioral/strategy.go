package behavioral

import "fmt"

/*
STRATEGY PATTERN

Purpose: Define a family of algorithms, encapsulate each one, and make them interchangeable.

Use Cases:
- Payment methods
- Sorting algorithms
- Compression algorithms
- Routing strategies

Go-Specific Implementation:
- Interface for strategy
- Context holds strategy reference
- Strategy can be changed at runtime
*/

// PaymentStrategy defines the payment interface
type PaymentStrategy interface {
	Pay(amount float64) error
}

// CreditCardStrategy implements credit card payment
type CreditCardStrategy struct {
	CardNumber string
	CVV        string
}

func (c *CreditCardStrategy) Pay(amount float64) error {
	fmt.Printf("💳 Paid $%.2f using Credit Card ending in %s\n", 
		amount, c.CardNumber[len(c.CardNumber)-4:])
	return nil
}

// PayPalStrategy implements PayPal payment
type PayPalStrategy struct {
	Email string
}

func (p *PayPalStrategy) Pay(amount float64) error {
	fmt.Printf("💰 Paid $%.2f using PayPal account %s\n", amount, p.Email)
	return nil
}

// CryptoStrategy implements cryptocurrency payment
type CryptoStrategy struct {
	WalletAddress string
	Currency      string
}

func (c *CryptoStrategy) Pay(amount float64) error {
	fmt.Printf("₿ Paid $%.2f using %s to wallet %s\n", 
		amount, c.Currency, c.WalletAddress[:10]+"...")
	return nil
}

// PaymentContext uses a payment strategy
type PaymentContext struct {
	strategy PaymentStrategy
}

func NewPaymentContext(strategy PaymentStrategy) *PaymentContext {
	return &PaymentContext{strategy: strategy}
}

func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
	p.strategy = strategy
}

func (p *PaymentContext) ExecutePayment(amount float64) error {
	return p.strategy.Pay(amount)
}

// Sorting Strategy Example
type SortStrategy interface {
	Sort([]int) []int
}

// BubbleSortStrategy implements bubble sort
type BubbleSortStrategy struct{}

func (b *BubbleSortStrategy) Sort(data []int) []int {
	fmt.Println("🔄 Using Bubble Sort")
	result := make([]int, len(data))
	copy(result, data)
	
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result)-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

// QuickSortStrategy implements quick sort
type QuickSortStrategy struct{}

func (q *QuickSortStrategy) Sort(data []int) []int {
	fmt.Println("⚡ Using Quick Sort")
	result := make([]int, len(data))
	copy(result, data)
	q.quickSort(result, 0, len(result)-1)
	return result
}

func (q *QuickSortStrategy) quickSort(arr []int, low, high int) {
	if low < high {
		pi := q.partition(arr, low, high)
		q.quickSort(arr, low, pi-1)
		q.quickSort(arr, pi+1, high)
	}
}

func (q *QuickSortStrategy) partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// SortContext uses a sorting strategy
type SortContext struct {
	strategy SortStrategy
}

func NewSortContext(strategy SortStrategy) *SortContext {
	return &SortContext{strategy: strategy}
}

func (s *SortContext) SetStrategy(strategy SortStrategy) {
	s.strategy = strategy
}

func (s *SortContext) Sort(data []int) []int {
	return s.strategy.Sort(data)
}

