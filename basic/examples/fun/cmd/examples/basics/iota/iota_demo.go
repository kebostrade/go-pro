package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Iota & Constants Demo")

	demo1BasicIota()
	demo2SkippingValues()
	demo3BitFlags()
	demo4CustomExpressions()
	demo5RealWorldExamples()
}

func demo1BasicIota() {
	utils.PrintSubHeader("1. Basic Iota")

	const (
		First  = iota // 0
		Second        // 1
		Third         // 2
		Fourth        // 3
	)

	fmt.Println("Basic iota (starts at 0):")
	fmt.Printf("  First:  %d\n", First)
	fmt.Printf("  Second: %d\n", Second)
	fmt.Printf("  Third:  %d\n", Third)
	fmt.Printf("  Fourth: %d\n", Fourth)

	// Starting from 1
	const (
		Monday    = iota + 1 // 1
		Tuesday              // 2
		Wednesday            // 3
		Thursday             // 4
		Friday               // 5
		Saturday             // 6
		Sunday               // 7
	)

	fmt.Println("\nDays of the week (starting from 1):")
	fmt.Printf("  Monday:    %d\n", Monday)
	fmt.Printf("  Tuesday:   %d\n", Tuesday)
	fmt.Printf("  Wednesday: %d\n", Wednesday)
	fmt.Printf("  Thursday:  %d\n", Thursday)
	fmt.Printf("  Friday:    %d\n", Friday)
	fmt.Printf("  Saturday:  %d\n", Saturday)
	fmt.Printf("  Sunday:    %d\n", Sunday)
}

func demo2SkippingValues() {
	utils.PrintSubHeader("2. Skipping Values")

	const (
		Zero = iota // 0
		_           // 1 (skipped)
		Two         // 2
		_           // 3 (skipped)
		Four        // 4
	)

	fmt.Println("Skipping values with underscore:")
	fmt.Printf("  Zero: %d\n", Zero)
	fmt.Printf("  Two:  %d\n", Two)
	fmt.Printf("  Four: %d\n", Four)
}

func demo3BitFlags() {
	utils.PrintSubHeader("3. Bit Flags (Powers of 2)")

	const (
		FlagNone    = 0
		FlagRead    = 1 << iota // 1 << 0 = 1 (binary: 001)
		FlagWrite               // 1 << 1 = 2 (binary: 010)
		FlagExecute             // 1 << 2 = 4 (binary: 100)
	)

	fmt.Println("File permissions (bit flags):")
	fmt.Printf("  None:    %d (binary: %03b)\n", FlagNone, FlagNone)
	fmt.Printf("  Read:    %d (binary: %03b)\n", FlagRead, FlagRead)
	fmt.Printf("  Write:   %d (binary: %03b)\n", FlagWrite, FlagWrite)
	fmt.Printf("  Execute: %d (binary: %03b)\n", FlagExecute, FlagExecute)

	// Combining flags
	readWrite := FlagRead | FlagWrite
	readExecute := FlagRead | FlagExecute
	all := FlagRead | FlagWrite | FlagExecute

	fmt.Println("\nCombined permissions:")
	fmt.Printf("  Read+Write:   %d (binary: %03b)\n", readWrite, readWrite)
	fmt.Printf("  Read+Execute: %d (binary: %03b)\n", readExecute, readExecute)
	fmt.Printf("  All:          %d (binary: %03b)\n", all, all)

	// Checking flags
	fmt.Println("\nChecking permissions:")
	fmt.Printf("  Has Read? %v\n", readWrite&FlagRead != 0)
	fmt.Printf("  Has Write? %v\n", readWrite&FlagWrite != 0)
	fmt.Printf("  Has Execute? %v\n", readWrite&FlagExecute != 0)
}

func demo4CustomExpressions() {
	utils.PrintSubHeader("4. Custom Expressions")

	const (
		_  = iota             // 0 (skip)
		KB = 1 << (10 * iota) // 1 << 10 = 1024
		MB                    // 1 << 20 = 1048576
		GB                    // 1 << 30 = 1073741824
		TB                    // 1 << 40 = 1099511627776
	)

	fmt.Println("Data sizes:")
	fmt.Printf("  1 KB = %d bytes\n", KB)
	fmt.Printf("  1 MB = %d bytes\n", MB)
	fmt.Printf("  1 GB = %d bytes\n", GB)
	fmt.Printf("  1 TB = %d bytes\n", TB)

	// Multipliers
	const (
		One      = 1
		Ten      = 10 * iota // 0
		Hundred              // 10
		Thousand             // 20
	)

	fmt.Println("\nMultipliers:")
	fmt.Printf("  One:      %d\n", One)
	fmt.Printf("  Ten:      %d\n", Ten)
	fmt.Printf("  Hundred:  %d\n", Hundred)
	fmt.Printf("  Thousand: %d\n", Thousand)
}

func demo5RealWorldExamples() {
	utils.PrintSubHeader("5. Real-World Examples")

	// HTTP Status Codes
	const (
		StatusOK                  = 200
		StatusCreated             = 201
		StatusBadRequest          = 400
		StatusUnauthorized        = 401
		StatusForbidden           = 403
		StatusNotFound            = 404
		StatusInternalServerError = 500
	)

	fmt.Println("HTTP Status Codes:")
	fmt.Printf("  OK:                    %d\n", StatusOK)
	fmt.Printf("  Created:               %d\n", StatusCreated)
	fmt.Printf("  Bad Request:           %d\n", StatusBadRequest)
	fmt.Printf("  Unauthorized:          %d\n", StatusUnauthorized)
	fmt.Printf("  Not Found:             %d\n", StatusNotFound)
	fmt.Printf("  Internal Server Error: %d\n", StatusInternalServerError)

	// Log Levels
	const (
		LogLevelDebug = iota
		LogLevelInfo
		LogLevelWarning
		LogLevelError
		LogLevelFatal
	)

	fmt.Println("\nLog Levels:")
	fmt.Printf("  Debug:   %d\n", LogLevelDebug)
	fmt.Printf("  Info:    %d\n", LogLevelInfo)
	fmt.Printf("  Warning: %d\n", LogLevelWarning)
	fmt.Printf("  Error:   %d\n", LogLevelError)
	fmt.Printf("  Fatal:   %d\n", LogLevelFatal)

	// Connection States
	const (
		StateDisconnected = iota
		StateConnecting
		StateConnected
		StateDisconnecting
	)

	fmt.Println("\nConnection States:")
	fmt.Printf("  Disconnected:   %d\n", StateDisconnected)
	fmt.Printf("  Connecting:     %d\n", StateConnecting)
	fmt.Printf("  Connected:      %d\n", StateConnected)
	fmt.Printf("  Disconnecting:  %d\n", StateDisconnecting)
}
