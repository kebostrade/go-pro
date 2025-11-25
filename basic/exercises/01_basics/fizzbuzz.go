package main

import (
	"fmt"
	"strings"
)

// Exercise: FizzBuzz
// Write a program that prints numbers from 1 to 100
// For multiples of 3, print "Fizz" instead of the number
// For multiples of 5, print "Buzz" instead of the number
// For multiples of both 3 and 5, print "FizzBuzz"

// TODO: Implement the FizzBuzz function
func FizzBuzz(n int) string {
	// Your code here
	// Hint: Check divisibility by 15 first, then 3, then 5
	return ""
}

func main() {
	fmt.Println("FizzBuzz Challenge")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()

	for i := 1; i <= 100; i++ {
		result := FizzBuzz(i)
		if result == "" {
			fmt.Printf("%3d: (not implemented)\n", i)
		} else {
			fmt.Printf("%3d: %s\n", i, result)
		}
	}
}
