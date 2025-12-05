package exercises

import (
	"fmt"
	"strconv"
	"strings"
)

// Exercise 1: Pointer Basics

// DoubleValue takes an integer pointer and doubles the value it points to
func DoubleValue(ptr *int) {
	*ptr = *ptr * 2
}

// SwapStrings swaps the values of two string pointers
func SwapStrings(a, b *string) {
	*a, *b = *b, *a
}

// GetLargerPointer returns a pointer to the larger of two integers
func GetLargerPointer(a, b *int) *int {
	if *a >= *b {
		return a
	}
	return b
}

// Exercise 2: Function Parameters with Pointers

// IncrementCounter increments a counter through a pointer and returns the new value
func IncrementCounter(counter *int) int {
	*counter++
	return *counter
}

// AppendString appends a string to another string through a pointer
func AppendString(target *string, suffix string) {
	*target = *target + suffix
}

// FindMinMax finds the minimum and maximum values in a slice
func FindMinMax(numbers []int, min, max *int) {
	if len(numbers) == 0 {
		return
	}
	*min = numbers[0]
	*max = numbers[0]
	for _, num := range numbers {
		if num < *min {
			*min = num
		}
		if num > *max {
			*max = num
		}
	}
}

// Exercise 3: Struct Methods with Pointer Receivers

type BankAccount struct {
	AccountNumber string
	Balance       float64
	Owner         string
}

// Deposit deposits money with pointer receiver
func (ba *BankAccount) Deposit(amount float64) {
	ba.Balance += amount
}

// Withdraw withdraws money with pointer receiver
func (ba *BankAccount) Withdraw(amount float64) bool {
	if ba.Balance >= amount {
		ba.Balance -= amount
		return true
	}
	return false
}

// GetAccountInfo returns account info with value receiver
func (ba BankAccount) GetAccountInfo() string {
	return fmt.Sprintf("Account: %s, Owner: %s, Balance: %.2f",
		ba.AccountNumber, ba.Owner, ba.Balance)
}

// TransferTo transfers money to another account with pointer receiver
func (ba *BankAccount) TransferTo(target *BankAccount, amount float64) bool {
	if ba.Withdraw(amount) {
		target.Deposit(amount)
		return true
	}
	return false
}

// Exercise 4: Memory Management and Allocation

// AllocateInt allocates and returns a pointer to a new integer
func AllocateInt(value int) *int {
	v := value
	return &v
}

// CreatePointerSlice creates a slice of pointers to integers
func CreatePointerSlice(size int) []*int {
	result := make([]*int, size)
	for i := 0; i < size; i++ {
		v := i
		result[i] = &v
	}
	return result
}

// SafeDereference safely dereferences a pointer
func SafeDereference(ptr *int, defaultValue int) int {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// Exercise 5: Linked List Implementation

type ListNode struct {
	Value int
	Next  *ListNode
}

type LinkedList struct {
	Head *ListNode
	Size int
}

// NewLinkedList creates a new empty linked list
func NewLinkedList() *LinkedList {
	return &LinkedList{
		Head: nil,
		Size: 0,
	}
}

// PrependNode adds a new node at the beginning of the list
func (ll *LinkedList) PrependNode(value int) {
	newNode := &ListNode{Value: value, Next: ll.Head}
	ll.Head = newNode
	ll.Size++
}

// AppendNode adds a new node at the end of the list
func (ll *LinkedList) AppendNode(value int) {
	newNode := &ListNode{Value: value, Next: nil}
	if ll.Head == nil {
		ll.Head = newNode
	} else {
		current := ll.Head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newNode
	}
	ll.Size++
}

// RemoveValue removes the first occurrence of a value from the list
func (ll *LinkedList) RemoveValue(value int) bool {
	if ll.Head == nil {
		return false
	}
	if ll.Head.Value == value {
		ll.Head = ll.Head.Next
		ll.Size--
		return true
	}
	current := ll.Head
	for current.Next != nil {
		if current.Next.Value == value {
			current.Next = current.Next.Next
			ll.Size--
			return true
		}
		current = current.Next
	}
	return false
}

// FindNode finds a value in the list and returns a pointer to the node
func (ll *LinkedList) FindNode(value int) *ListNode {
	current := ll.Head
	for current != nil {
		if current.Value == value {
			return current
		}
		current = current.Next
	}
	return nil
}

// ToSlice converts the linked list to a slice
func (ll *LinkedList) ToSlice() []int {
	result := make([]int, 0, ll.Size)
	current := ll.Head
	for current != nil {
		result = append(result, current.Value)
		current = current.Next
	}
	return result
}

// Exercise 6: Performance Optimization

type LargeData struct {
	Numbers [1000]int
	Text    string
	Active  bool
}

// ProcessByValue processes LargeData by value
func ProcessByValue(data LargeData) int {
	sum := 0
	for _, num := range data.Numbers {
		sum += num
	}
	return sum
}

// ProcessByPointer processes LargeData by pointer
func ProcessByPointer(data *LargeData) int {
	sum := 0
	for _, num := range data.Numbers {
		sum += num
	}
	return sum
}

// InitializeLargeData modifies LargeData efficiently
func InitializeLargeData(data *LargeData, text string) {
	for i := 0; i < len(data.Numbers); i++ {
		data.Numbers[i] = i
	}
	data.Text = text
	data.Active = true
}

// Exercise 7: Advanced Pointer Patterns

// ParseNameAge parses a "name:age" string and sets the values through pointers
func ParseNameAge(input string, name *string, age *int) bool {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return false
	}
	parsedAge, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return false
	}
	*name = strings.TrimSpace(parts[0])
	*age = parsedAge
	return true
}

// ModifyThroughPointers modifies all values in the slice by adding the given increment
func ModifyThroughPointers(ptrs []*int, increment int) {
	for _, ptr := range ptrs {
		if ptr != nil {
			*ptr += increment
		}
	}
}

// CreateCounter implements a simple reference counter
func CreateCounter(initialValue int) *int {
	counter := initialValue
	return &counter
}

// Exercise 8: Pointer Safety and Best Practices

// SafeStringCopy safely copies a string through pointers
func SafeStringCopy(src *string, dst *string) bool {
	if src == nil || dst == nil {
		return false
	}
	*dst = *src
	return true
}

// SumValidPointers validates and processes a slice of pointers
func SumValidPointers(ptrs []*int) int {
	sum := 0
	for _, ptr := range ptrs {
		if ptr != nil {
			sum += *ptr
		}
	}
	return sum
}

// DeepCopySlice creates a deep copy of a slice of integers
func DeepCopySlice(original []int) []int {
	if original == nil {
		return nil
	}
	copied := make([]int, len(original))
	copy(copied, original)
	return copied
}
