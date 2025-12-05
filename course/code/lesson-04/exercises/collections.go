package exercises

// Exercise 1: Array Practice

// GetFirstFivePrimes returns an array of the first 5 prime numbers
func GetFirstFivePrimes() [5]int {
	return [5]int{2, 3, 5, 7, 11}
}

// FindMaxInArray finds the maximum value in an array
func FindMaxInArray(arr [10]int) (max int, index int) {
	max = arr[0]
	index = 0
	for i := 1; i < len(arr); i++ {
		if arr[i] > max {
			max = arr[i]
			index = i
		}
	}
	return max, index
}

// Exercise 2: Slice Manipulation

// RemoveDuplicates removes duplicates from a slice
func RemoveDuplicates(slice []int) []int {
	if len(slice) == 0 {
		return slice
	}
	seen := make(map[int]bool)
	result := []int{}
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// ReverseSlice reverses a slice in place
func ReverseSlice(slice []string) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - 1 - i
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// MergeSortedSlices merges two sorted slices into one sorted slice
func MergeSortedSlices(slice1, slice2 []int) []int {
	result := make([]int, 0, len(slice1)+len(slice2))
	i, j := 0, 0
	for i < len(slice1) && j < len(slice2) {
		if slice1[i] < slice2[j] {
			result = append(result, slice1[i])
			i++
		} else {
			result = append(result, slice2[j])
			j++
		}
	}
	result = append(result, slice1[i:]...)
	result = append(result, slice2[j:]...)
	return result
}

// Exercise 3: Map Operations

// CountCharacters counts the frequency of each character in a string
func CountCharacters(s string) map[rune]int {
	counts := make(map[rune]int)
	for _, ch := range s {
		counts[ch]++
	}
	return counts
}

// InvertMap inverts a map (keys become values, values become keys)
func InvertMap(m map[string]int) map[int]string {
	inverted := make(map[int]string)
	for k, v := range m {
		inverted[v] = k
	}
	return inverted
}

// MergeMaps merges two maps
func MergeMaps(map1, map2 map[string]int) map[string]int {
	result := make(map[string]int)
	for k, v := range map1 {
		result[k] = v
	}
	for k, v := range map2 {
		result[k] = v
	}
	return result
}

// Exercise 4: Advanced Collection Operations

// FindIntersection finds the intersection of two slices
func FindIntersection(slice1, slice2 []int) []int {
	set := make(map[int]bool)
	for _, v := range slice1 {
		set[v] = true
	}
	result := []int{}
	seen := make(map[int]bool)
	for _, v := range slice2 {
		if set[v] && !seen[v] {
			result = append(result, v)
			seen[v] = true
		}
	}
	return result
}

// GroupByLength groups strings by their length
func GroupByLength(words []string) map[int][]string {
	groups := make(map[int][]string)
	for _, word := range words {
		length := len(word)
		groups[length] = append(groups[length], word)
	}
	return groups
}

// Exercise 5: Real-World Scenario - Inventory Management

type Product struct {
	ID    string
	Name  string
	Price float64
	Stock int
}

type Inventory struct {
	products map[string]*Product
}

// NewInventory creates a new inventory
func NewInventory() *Inventory {
	return &Inventory{
		products: make(map[string]*Product),
	}
}

// AddProduct adds a product to the inventory
func (inv *Inventory) AddProduct(product *Product) {
	inv.products[product.ID] = product
}

// GetProduct gets a product by ID
func (inv *Inventory) GetProduct(id string) (*Product, bool) {
	product, exists := inv.products[id]
	return product, exists
}

// UpdateStock updates product stock
func (inv *Inventory) UpdateStock(id string, newStock int) bool {
	product, exists := inv.products[id]
	if !exists {
		return false
	}
	product.Stock = newStock
	return true
}

// GetLowStockProducts gets all products with stock below a threshold
func (inv *Inventory) GetLowStockProducts(threshold int) []*Product {
	lowStock := []*Product{}
	for _, product := range inv.products {
		if product.Stock < threshold {
			lowStock = append(lowStock, product)
		}
	}
	return lowStock
}

// GetTotalValue calculates total inventory value
func (inv *Inventory) GetTotalValue() float64 {
	total := 0.0
	for _, product := range inv.products {
		total += product.Price * float64(product.Stock)
	}
	return total
}

// Exercise 6: Memory Efficiency Challenge

// EfficientAppend efficiently appends to a slice
func EfficientAppend(initialSize int, values []int) []int {
	result := make([]int, 0, initialSize+len(values))
	result = append(result, values...)
	return result
}

// ProcessInChunks processes large slices in chunks
func ProcessInChunks(data []int, chunkSize int, processor func([]int) int) []int {
	results := []int{}
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		results = append(results, processor(chunk))
	}
	return results
}
