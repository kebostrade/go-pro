//go:build benchmark
// +build benchmark

package testingexamples

// This file contains benchmark examples
// Run with: go test -tags=benchmark -bench=. -v ./...

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// ========================================
// EXAMPLE 1: BASIC BENCHMARK
// ========================================

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func BenchmarkFib10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(10)
	}
}

func BenchmarkFib20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(20)
	}
}

// ========================================
// EXAMPLE 2: BENCHMARK WITH SETUP
// ========================================

type DataProcessor struct {
	data []int
}

func NewDataProcessor(size int) *DataProcessor {
	return &DataProcessor{
		data: make([]int, size),
	}
}

func (p *DataProcessor) Process() int {
	sum := 0
	for _, v := range p.data {
		sum += v
	}
	return sum
}

func (p *DataProcessor) ProcessWithGoroutine() int {
	sum := 0
	ch := make(chan int)

	for _, v := range p.data {
		go func(val int) {
			ch <- val
		}(v)
	}

	for i := 0; i < len(p.data); i++ {
		sum += <-ch
	}

	return sum
}

func BenchmarkDataProcess(b *testing.B) {
	processor := NewDataProcessor(1000)

	b.ResetTimer() // Reset timer after setup
	for i := 0; i < b.N; i++ {
		processor.Process()
	}
}

func BenchmarkDataProcessWithGoroutine(b *testing.B) {
	processor := NewDataProcessor(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.ProcessWithGoroutine()
	}
}

// ========================================
// EXAMPLE 3: BENCHMARK WITH PAUSE/RESUME
// ========================================

func BenchmarkWithPause(b *testing.B) {
	data := generateLargeDataset()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer() // Pause timer
		// Do expensive setup that shouldn't be measured
		processed := prepareData(data)

		b.StartTimer() // Resume timer
		// Benchmark the actual operation
		result := processData(processed)
		_ = result
	}
}

func generateLargeDataset() []int {
	return make([]int, 1000)
}

func prepareData(data []int) []int {
	// Expensive preparation
	return data
}

func processData(data []int) int {
	sum := 0
	for _, v := range data {
		sum += v
	}
	return sum
}

// ========================================
// EXAMPLE 4: BENCHMARK SUB-ALLOCATIONS
// ========================================

func SprintfConcat(strs []string) string {
	var result string
	for _, s := range strs {
		result += fmt.Sprintf("%s ", s)
	}
	return result
}

func StringsBuilderConcat(strs []string) string {
	var builder strings.Builder
	for _, s := range strs {
		builder.WriteString(s)
		builder.WriteString(" ")
	}
	return builder.String()
}

func BytesBufferConcat(strs []string) string {
	var buffer bytes.Buffer
	for _, s := range strs {
		buffer.WriteString(s)
		buffer.WriteString(" ")
	}
	return buffer.String()
}

func BenchmarkSprintfConcat(b *testing.B) {
	strs := []string{"hello", "world", "foo", "bar", "baz"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SprintfConcat(strs)
	}
}

func BenchmarkStringsBuilderConcat(b *testing.B) {
	strs := []string{"hello", "world", "foo", "bar", "baz"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringsBuilderConcat(strs)
	}
}

func BenchmarkBytesBufferConcat(b *testing.B) {
	strs := []string{"hello", "world", "foo", "bar", "baz"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BytesBufferConcat(strs)
	}
}

// ========================================
// EXAMPLE 5: BENCHMARK WITH DIFFERENT SIZES
// ========================================

func BenchmarkSort10(b *testing.B) {
	benchmarkSort(b, 10)
}

func BenchmarkSort100(b *testing.B) {
	benchmarkSort(b, 100)
}

func BenchmarkSort1000(b *testing.B) {
	benchmarkSort(b, 1000)
}

func benchmarkSort(b *testing.B, size int) {
	b.StopTimer()
	data := generateRandomData(size)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// Create copy for fair comparison
		d := make([]int, len(data))
		copy(d, data)
		sort.Ints(d)
	}
}

func generateRandomData(size int) []int {
	data := make([]int, size)
	for i := range data {
		data[i] = rand.Intn(1000)
	}
	return data
}

// ========================================
// EXAMPLE 6: PARALLEL BENCHMARKS
// ========================================

func CounterMapSequential(m map[int]int, iterations int) {
	for i := 0; i < iterations; i++ {
		m[i%100]++
	}
}

func CounterMapParallel(m map[int]int, iterations int) {
	ch := make(chan int, iterations)

	// Writers
	for i := 0; i < 4; i++ {
		go func() {
			for j := range ch {
				m[j%100]++
			}
		}()
	}

	// Send work
	for i := 0; i < iterations; i++ {
		ch <- i
	}
	close(ch)
}

func BenchmarkCounterSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := make(map[int]int)
		b.StartTimer()

		CounterMapSequential(m, 1000)
	}
}

func BenchmarkCounterParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := make(map[int]int)
		b.StartTimer()

		CounterMapParallel(m, 1000)
	}
}

// ========================================
// EXAMPLE 7: BENCHMARK ALLOCATIONS
// ========================================

// High allocation version
func AppendNaive() []int {
	var result []int
	for i := 0; i < 1000; i++ {
		result = append(result, i)
	}
	return result
}

// Pre-allocated version (fewer allocations)
func AppendPreAllocated() []int {
	result := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		result = append(result, i)
	}
	return result
}

func BenchmarkAppendNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AppendNaive()
	}
}

func BenchmarkAppendPreAllocated(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AppendPreAllocated()
	}
}

// ========================================
// EXAMPLE 8: BENCHMARK STRING OPERATIONS
// ========================================

// String concatenation with +
func StringPlus(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s
	}
	return result
}

// String concatenation with strings.Builder
func StringBuilder(strs []string) string {
	var builder strings.Builder
	for _, s := range strs {
		builder.WriteString(s)
	}
	return builder.String()
}

func BenchmarkStringPlus10(b *testing.B) {
	strs := make([]string, 10)
	for i := range strs {
		strs[i] = "test"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringPlus(strs)
	}
}

func BenchmarkStringBuilder10(b *testing.B) {
	strs := make([]string, 10)
	for i := range strs {
		strs[i] = "test"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringBuilder(strs)
	}
}

func BenchmarkStringPlus100(b *testing.B) {
	strs := make([]string, 100)
	for i := range strs {
		strs[i] = "test"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringPlus(strs)
	}
}

func BenchmarkStringBuilder100(b *testing.B) {
	strs := make([]string, 100)
	for i := range strs {
		strs[i] = "test"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringBuilder(strs)
	}
}

// ========================================
// EXAMPLE 9: BENCHMARK JSON MARSHALLING
// ========================================

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Active   bool   `json:"active"`
	Balance  float64 `json:"balance"`
}

func GenerateUsers(count int) []User {
	users := make([]User, count)
	for i := 0; i < count; i++ {
		users[i] = User{
			ID:      i,
			Name:    fmt.Sprintf("User%d", i),
			Email:   fmt.Sprintf("user%d@example.com", i),
			Age:     20 + (i % 60),
			Active:  i%2 == 0,
			Balance: float64(i * 100),
		}
	}
	return users
}

func BenchmarkJSONMarshalSingle(b *testing.B) {
	user := User{ID: 1, Name: "John", Email: "john@example.com", Age: 30}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshalSingle(b *testing.B) {
	user := User{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	data, _ := json.Marshal(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var u User
		err := json.Unmarshal(data, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONMarshalArray100(b *testing.B) {
	users := GenerateUsers(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(users)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ========================================
// EXAMPLE 10: BENCHMARK WITH REPORT METRICS
// ========================================

func BenchmarkWithMetrics(b *testing.B) {
	data := make([]int, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Process data
		sum := 0
		for _, v := range data {
			sum += v
		}
		_ = sum
	}

	// Report additional metrics
	b.ReportMetric(float64(b.Elapsed().Milliseconds()), "ms/op")
}

// ========================================
// EXAMPLE 11: MEMORY ALLOCATION BENCHMARKS
// ========================================

func AllocateInLoop() []int {
	result := make([]int, 0)
	for i := 0; i < 100; i++ {
		// Allocates new memory each append
		result = append(result, i)
	}
	return result
}

func AllocateOnce() []int {
	// Single allocation
	result := make([]int, 100)
	for i := 0; i < 100; i++ {
		result[i] = i
	}
	return result
}

func BenchmarkAllocateInLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AllocateInLoop()
	}
}

func BenchmarkAllocateOnce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AllocateOnce()
	}
}

// ========================================
// EXAMPLE 12: COMPARING DATA STRUCTURES
// ========================================

func BenchmarkSliceInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := make([]int, 0, 100)
		for j := 0; j < 100; j++ {
			slice = append(slice, j)
		}
	}
}

func BenchmarkMapInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[int]bool)
		for j := 0; j < 100; j++ {
			m[j] = true
		}
	}
}

// ========================================
// EXAMPLE 13: BENCHMARK WITH CPU PROFLING
// ========================================

func ExpensiveComputation(n int) int {
	result := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			result += i * j
		}
	}
	return result
}

func BenchmarkExpensiveComputation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExpensiveComputation(100)
	}
}

// Run with:
// go test -bench=BenchmarkExpensiveComputation -cpuprofile=cpu.prof
// go tool pprof cpu.prof
// (pprof) top
// (pprof) list ExpensiveComputation

// ========================================
// EXAMPLE 14: BENCHMARK WITH MEMORY PROFLING
// ========================================

func MemoryIntensiveOperation() []byte {
	data := make([][]byte, 100)
	for i := range data {
		data[i] = make([]byte, 1024) // 1KB each
	}

	result := make([]byte, 0, len(data)*1024)
	for _, d := range data {
		result = append(result, d...)
	}
	return result
}

func BenchmarkMemoryIntensive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MemoryIntensiveOperation()
	}
}

// Run with:
// go test -bench=BenchmarkMemoryIntensive -memprofile=mem.prof
// go tool pprof mem.prof
// (pprof) top
// (pprof) list MemoryIntensiveOperation

// ========================================
// EXAMPLE 15: BENCHMARK TIMER CONTROL
// ========================================

func BenchmarkWithTimerControl(b *testing.B) {
	// Setup phase (not timed)
	data := prepareBenchmarkData()

	b.ResetTimer() // Reset timer to exclude setup time

	for i := 0; i < b.N; i++ {
		// Benchmark this
		result := processBenchmarkData(data)
		_ = result
	}

	b.StopTimer() // Stop timer
	// Teardown (not timed)
	cleanupBenchmarkData()
}

func prepareBenchmarkData() []int {
	return make([]int, 1000)
}

func processBenchmarkData(data []int) int {
	sum := 0
	for _, v := range data {
		sum += v
	}
	return sum
}

func cleanupBenchmarkData() {
	// Cleanup resources
}

// ========================================
// EXAMPLE 16: BENCHMARK WITH VARYING INPUTS
// ========================================

func BenchmarkStringToInt(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{"small", "123"},
		{"medium", "123456789"},
		{"large", "1234567890123456789"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				strconv.Atoi(tt.input)
			}
		})
	}
}

// ========================================
// EXAMPLE 17: BENCHMARK CONCURRENT OPERATIONS
// ========================================

func SerialSum(data []int) int {
	sum := 0
	for _, v := range data {
		sum += v
	}
	return sum
}

func ConcurrentSum(data []int) int {
	sum := 0
	ch := make(chan int, len(data))

	for _, v := range data {
		go func(val int) {
			ch <- val
		}(v)
	}

	for i := 0; i < len(data); i++ {
		sum += <-ch
	}

	return sum
}

func BenchmarkSerialSum(b *testing.B) {
	data := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SerialSum(data)
	}
}

func BenchmarkConcurrentSum(b *testing.B) {
	data := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcurrentSum(data)
	}
}

// ========================================
// EXAMPLE 18: BENCHMARK GC PRESSURE
// ========================================

func HighGCPressure() {
	for i := 0; i < 1000; i++ {
		data := make([]byte, 1024)
		_ = data
	}
}

func LowGCPressure() {
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = make([]byte, 1024)
	}
	_ = data
}

func BenchmarkHighGCPressure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HighGCPressure()
	}
}

func BenchmarkLowGCPressure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LowGCPressure()
	}
}

// ========================================
// EXAMPLE 19: BENCHMARK WITH CUSTOM ALLOCATION
// ========================================

type BufferPool struct {
	buffers chan *bytes.Buffer
}

func NewBufferPool(size int) *BufferPool {
	pool := &BufferPool{
		buffers: make(chan *bytes.Buffer, size),
	}
	for i := 0; i < size; i++ {
		pool.buffers <- &bytes.Buffer{}
	}
	return pool
}

func (p *BufferPool) Get() *bytes.Buffer {
	return <-p.buffers
}

func (p *BufferPool) Put(buf *bytes.Buffer) {
	buf.Reset()
	p.buffers <- buf
}

var globalPool = NewBufferPool(10)

func WriteWithPool(data string) string {
	buf := globalPool.Get()
	defer globalPool.Put(buf)

	buf.WriteString(data)
	return buf.String()
}

func WriteWithoutPool(data string) string {
	var buf bytes.Buffer
	buf.WriteString(data)
	return buf.String()
}

func BenchmarkWriteWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WriteWithPool("test data")
	}
}

func BenchmarkWriteWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WriteWithoutPool("test data")
	}
}

// ========================================
// EXAMPLE 20: MEASURING CACHE EFFECTS
// ========================================

func SequentialAccess(data [][]int) int {
	sum := 0
	for _, row := range data {
		for _, val := range row {
			sum += val
		}
	}
	return sum
}

func RandomAccess(data [][]int) int {
	sum := 0
	for _, row := range data {
		for _, val := range row {
			sum += val
		}
	}
	return sum
}

func BenchmarkSequentialAccess(b *testing.B) {
	data := make2DArray(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SequentialAccess(data)
	}
}

func BenchmarkRandomAccess(b *testing.B) {
	data := make2DArray(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RandomAccess(data)
	}
}

func make2DArray(rows, cols int) [][]int {
	data := make([][]int, rows)
	for i := range data {
		data[i] = make([]int, cols)
		for j := range data[i] {
			data[i][j] = i + j
		}
	}
	return data
}

// ========================================
// CPU PROFILING EXAMPLE
// ========================================

// Enable CPU profiling:
// go test -bench=. -cpuprofile=cpu.prof
// go tool pprof cpu.prof
//
// Common pprof commands:
// (pprof) top          - Show top functions
// (pprof) top10        - Show top 10 functions
// (pprof) list FuncName - Show function source
// (pprof) web          - Generate call graph (requires graphviz)
// (pprof) pdf          - Generate PDF

// ========================================
// MEMORY PROFILING EXAMPLE
// ========================================

// Enable memory profiling:
// go test -bench=. -memprofile=mem.prof
// go tool pprof mem.prof
//
// Common commands:
// (pprof) top          - Show top allocations
// (pprof) list FuncName - Show function with allocations
// (pprof) web          - Generate allocation graph

// ========================================
// BLOCK PROFILING EXAMPLE
// ========================================

// Enable block profiling:
// go test -bench=. -blockprofile=block.prof
// go tool pprof block.prof

// ========================================
// GOROUTINE PROFILING
// ========================================

// View goroutine profile in running application:
// curl http://localhost:6060/debug/pprof/goroutine?debug=2

// Or use pprof:
// go tool pprof http://localhost:6060/debug/pprof/goroutine

// ========================================
// EXAMPLE: USING runtime/pprof DIRECTLY
// ========================================

func DemonstrateCPUProfiling() {
	// Create CPU profile file
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	// Code to profile
	for i := 0; i < 1000000; i++ {
		_ = i * i
	}
}

func DemonstrateMemoryProfiling() {
	// Code to profile
	data := make([]byte, 1024*1024) // 1MB allocation
	_ = data

	// Create memory profile
	f, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write memory profile
	runtime.GC() // Run GC to get accurate profile
	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}
}

// ========================================
// BENCHMARK HELPER FUNCTIONS
// ========================================

// SkipBenchmark skips a benchmark based on a condition
func SkipBenchmark(b *testing.B, condition bool, reason string) {
	if condition {
		b.Skip(reason)
	}
}

// ReportAllocs reports memory allocations in benchmark
func BenchmarkReportAllocs(b *testing.B) {
	b.ReportAllocs() // Report memory allocations

	data := make([]int, 1000)
	for i := 0; i < b.N; i++ {
		sum := 0
		for _, v := range data {
			sum += v
		}
		_ = sum
	}
}
