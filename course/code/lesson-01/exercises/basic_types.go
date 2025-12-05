package exercises

import "fmt"

// Suppress unused import warning - fmt is used in exercise solutions
var _ = fmt.Sprint

// Exercise 1: Basic Types Practice
// Complete the following functions to work with Go's basic types

// PersonInfo represents information about a person
type PersonInfo struct {
	Name      string
	Age       int
	Height    float64
	IsStudent bool
}

// CreatePersonInfo creates and returns a PersonInfo struct with the given values
// Parameters:
//   - name: person's name (string)
//   - age: person's age (int)
//   - height: person's height in meters (float64)
//   - isStudent: whether the person is a student (bool)
//
// Returns: PersonInfo struct with the provided values
func CreatePersonInfo(name string, age int, height float64, isStudent bool) PersonInfo {
	return PersonInfo{
		Name:      name,
		Age:       age,
		Height:    height,
		IsStudent: isStudent,
	}
}

// FormatPersonInfo formats a PersonInfo struct into a readable string
// Parameter: info - PersonInfo struct to format
// Returns: formatted string with person's information
// Expected format: "Name: John, Age: 25, Height: 1.75m, Student: true"
func FormatPersonInfo(info PersonInfo) string {
	return fmt.Sprintf("Name: %s, Age: %d, Height: %.2fm, Student: %t",
		info.Name, info.Age, info.Height, info.IsStudent)
}

// CalculateBMI calculates the Body Mass Index given weight and height
// Parameters:
//   - weight: weight in kilograms (float64)
//   - height: height in meters (float64)
//
// Returns: BMI value (float64)
// Formula: BMI = weight / (height * height)
func CalculateBMI(weight, height float64) float64 {
	return weight / (height * height)
}

// ClassifyBMI classifies BMI into categories
// Parameter: bmi - BMI value (float64)
// Returns: BMI category (string)
// Categories:
//   - BMI < 18.5: "Underweight"
//   - 18.5 <= BMI < 25: "Normal weight"
//   - 25 <= BMI < 30: "Overweight"
//   - BMI >= 30: "Obese"
func ClassifyBMI(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "Underweight"
	case bmi < 25:
		return "Normal weight"
	case bmi < 30:
		return "Overweight"
	default:
		return "Obese"
	}
}

// ConvertTemperature converts temperature between Celsius and Fahrenheit
// Parameters:
//   - temp: temperature value (float64)
//   - fromUnit: source unit, either "C" or "F" (string)
//
// Returns: converted temperature (float64)
// Formulas:
//   - Celsius to Fahrenheit: F = C * 9/5 + 32
//   - Fahrenheit to Celsius: C = (F - 32) * 5/9
func ConvertTemperature(temp float64, fromUnit string) float64 {
	if fromUnit == "C" {
		// Celsius to Fahrenheit: F = C * 9/5 + 32
		return temp*9/5 + 32
	} else if fromUnit == "F" {
		// Fahrenheit to Celsius: C = (F - 32) * 5/9
		return (temp - 32) * 5 / 9
	}
	return temp // Return original if unit is not recognized
}

// IsValidAge checks if an age is valid (between 0 and 150)
// Parameter: age - age to validate (int)
// Returns: true if age is valid, false otherwise
func IsValidAge(age int) bool {
	return age >= 0 && age <= 150
}

// GetAgeCategory returns the age category for a given age
// Parameter: age - age to categorize (int)
// Returns: age category (string)
// Categories:
//   - 0-12: "Child"
//   - 13-19: "Teenager"
//   - 20-64: "Adult"
//   - 65+: "Senior"
func GetAgeCategory(age int) string {
	switch {
	case age <= 12:
		return "Child"
	case age <= 19:
		return "Teenager"
	case age <= 64:
		return "Adult"
	default:
		return "Senior"
	}
}

// CalculateCircleArea calculates the area of a circle given its radius
// Parameter: radius - radius of the circle (float64)
// Returns: area of the circle (float64)
// Formula: Area = π * radius²
// Use 3.14159 for π
func CalculateCircleArea(radius float64) float64 {
	const pi = 3.14159
	return pi * radius * radius
}

// IsEven checks if a number is even
// Parameter: number - number to check (int)
// Returns: true if number is even, false if odd
func IsEven(number int) bool {
	return number%2 == 0
}

// MaxOfThree returns the maximum of three integers
// Parameters: a, b, c - three integers to compare
// Returns: the largest of the three integers
func MaxOfThree(a, b, c int) int {
	max := a
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	return max
}
