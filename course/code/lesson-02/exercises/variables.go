package exercises

import "fmt"

// Suppress unused import warning - fmt is used in exercise solutions
var _ = fmt.Sprint

// Exercise 1: Variable Declarations and Scope

// DeclareVariables demonstrates different variable declaration methods
func DeclareVariables() (string, int, float64, bool) {
	var name string = "John Doe"
	var age = 30
	salary := 75000.50
	var isEmployed bool = true
	return name, age, salary, isEmployed
}

// MultipleDeclarations practices declaring multiple variables at once
func MultipleDeclarations() (int, int, int, string, string) {
	var x, y, z = 10, 20, 30
	a, b := "Hello", "World"
	return x, y, z, a, b
}

// BlockDeclaration uses block declaration syntax
func BlockDeclaration() (string, float64, bool) {
	var (
		projectName = "GO-PRO"
		version     = 2.1
		isStable    = true
	)
	return projectName, version, isStable
}

// Package-level variable for scope testing
var packageVar = "I'm at package level"

// TestVariableScope demonstrates variable scope rules
func TestVariableScope(input string) string {
	functionVar := "I'm in function"
	blockVar := ""
	if true {
		blockVar = "I'm in block"
	}
	return fmt.Sprintf("Input: %s, Package: %s, Function: %s, Block: %s",
		input, packageVar, functionVar, blockVar)
}

// SwapVariables swaps the values of two variables
func SwapVariables(a, b int) (int, int) {
	return b, a
}

// ZeroValues returns the zero values of different types
func ZeroValues() (int, float64, string, bool) {
	var zeroInt int
	var zeroFloat float64
	var zeroString string
	var zeroBool bool
	return zeroInt, zeroFloat, zeroString, zeroBool
}

// ConstantUsage demonstrates working with constants
func ConstantUsage() (float64, int, string) {
	const pi = 3.14159
	const maxRetries = 5
	const appName = "Learning Go"
	return pi, maxRetries, appName
}

// VariableReassignment demonstrates variable reassignment
func VariableReassignment(initial int) int {
	value := initial
	value += 10
	value *= 2
	value -= 5
	return value
}

// TypeInference demonstrates Go's type inference
func TypeInference() (int, float64, string, bool) {
	var inferredInt = 42
	var inferredFloat = 3.14
	var inferredString = "Go"
	var inferredBool = true
	return inferredInt, inferredFloat, inferredString, inferredBool
}

// ShadowingExample demonstrates variable shadowing
func ShadowingExample(outer int) string {
	value := outer
	outerValue := value
	innerValue := 0
	if true {
		value := outer * 2
		innerValue = value
	}
	return fmt.Sprintf("Outer: %d, Inner: %d", outerValue, innerValue)
}
