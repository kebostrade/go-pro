package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Variables & Types Demo")

	demo1BasicTypes()
	demo2TypeInference()
	demo3Constants()
	demo4ZeroValues()
	demo5TypeConversion()
}

func demo1BasicTypes() {
	utils.PrintSubHeader("1. Basic Types")

	// Integer types
	var intVar int = 42
	var int8Var int8 = 127
	var int16Var int16 = 32767
	var int32Var int32 = 2147483647
	var int64Var int64 = 9223372036854775807

	fmt.Println("Integer types:")
	fmt.Printf("  int:   %d\n", intVar)
	fmt.Printf("  int8:  %d\n", int8Var)
	fmt.Printf("  int16: %d\n", int16Var)
	fmt.Printf("  int32: %d\n", int32Var)
	fmt.Printf("  int64: %d\n", int64Var)

	// Unsigned integer types
	var uintVar uint = 42
	var uint8Var uint8 = 255
	var uint16Var uint16 = 65535
	var uint32Var uint32 = 4294967295
	var uint64Var uint64 = 18446744073709551615

	fmt.Println("\nUnsigned integer types:")
	fmt.Printf("  uint:   %d\n", uintVar)
	fmt.Printf("  uint8:  %d (also known as byte)\n", uint8Var)
	fmt.Printf("  uint16: %d\n", uint16Var)
	fmt.Printf("  uint32: %d\n", uint32Var)
	fmt.Printf("  uint64: %d\n", uint64Var)

	// Floating point types
	var float32Var float32 = 3.14159
	var float64Var float64 = 3.141592653589793

	fmt.Println("\nFloating point types:")
	fmt.Printf("  float32: %.5f\n", float32Var)
	fmt.Printf("  float64: %.15f\n", float64Var)

	// String and boolean
	var stringVar string = "Hello, Go!"
	var boolVar bool = true

	fmt.Println("\nOther types:")
	fmt.Printf("  string: %s\n", stringVar)
	fmt.Printf("  bool:   %v\n", boolVar)
}

func demo2TypeInference() {
	utils.PrintSubHeader("2. Type Inference (Short Declaration)")

	// Short variable declaration with type inference
	name := "Alice"
	age := 30
	height := 5.6
	isStudent := false

	fmt.Printf("name:      %s (type: %T)\n", name, name)
	fmt.Printf("age:       %d (type: %T)\n", age, age)
	fmt.Printf("height:    %.1f (type: %T)\n", height, height)
	fmt.Printf("isStudent: %v (type: %T)\n", isStudent, isStudent)

	// Multiple variable declaration
	x, y, z := 1, 2, 3
	fmt.Printf("\nMultiple declaration: x=%d, y=%d, z=%d\n", x, y, z)
}

func demo3Constants() {
	utils.PrintSubHeader("3. Constants")

	const Pi = 3.14159
	const AppName = "MyApp"
	const MaxConnections = 100

	fmt.Printf("Pi:             %.5f\n", Pi)
	fmt.Printf("AppName:        %s\n", AppName)
	fmt.Printf("MaxConnections: %d\n", MaxConnections)

	// Grouped constants
	const (
		StatusOK       = 200
		StatusNotFound = 404
		StatusError    = 500
	)

	fmt.Println("\nHTTP Status codes:")
	fmt.Printf("  OK:        %d\n", StatusOK)
	fmt.Printf("  Not Found: %d\n", StatusNotFound)
	fmt.Printf("  Error:     %d\n", StatusError)
}

func demo4ZeroValues() {
	utils.PrintSubHeader("4. Zero Values")

	var i int
	var f float64
	var b bool
	var s string

	fmt.Println("Uninitialized variables have zero values:")
	fmt.Printf("  int:     %d\n", i)
	fmt.Printf("  float64: %f\n", f)
	fmt.Printf("  bool:    %v\n", b)
	fmt.Printf("  string:  %q (empty string)\n", s)
}

func demo5TypeConversion() {
	utils.PrintSubHeader("5. Type Conversion")

	var i int = 42
	var f float64 = float64(i)
	var u uint = uint(f)

	fmt.Printf("int to float64: %d → %.1f\n", i, f)
	fmt.Printf("float64 to uint: %.1f → %d\n", f, u)

	// String conversion
	var num int = 65
	var char string = string(rune(num))

	fmt.Printf("\nASCII conversion: %d → %s\n", num, char)

	// Be careful with type conversion
	var largeInt int64 = 9223372036854775807
	var smallInt int32 = int32(largeInt) // Overflow!

	fmt.Printf("\nOverflow example:\n")
	fmt.Printf("  int64:  %d\n", largeInt)
	fmt.Printf("  int32:  %d (overflow!)\n", smallInt)
}
