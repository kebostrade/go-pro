package main

import (
	"fmt"
	"lesson-03/exercises"
)

func main() {
	// Lesson 03: Control Structures and Loops

	fmt.Println("=== Control Structures in Go ===")
	fmt.Println()

	// If/Else
	fmt.Println("1. If/Else Basic:")
	fmt.Printf("   IfElseBasic(5) = %s\n", exercises.IfElseBasic(5))
	fmt.Printf("   IfElseBasic(-3) = %s\n", exercises.IfElseBasic(-3))
	fmt.Printf("   IfElseBasic(0) = %s\n", exercises.IfElseBasic(0))
	fmt.Println()

	// If with initialization
	fmt.Println("2. If with Initialization:")
	fmt.Printf("   IfWithInit(-5) = %d\n", exercises.IfWithInit(-5))
	fmt.Printf("   IfWithInit(10) = %d\n", exercises.IfWithInit(10))
	fmt.Println()

	// Switch
	fmt.Println("3. Switch Day:")
	for d := 1; d <= 7; d++ {
		fmt.Printf("   SwitchDay(%d) = %s\n", d, exercises.SwitchDay(d))
	}
	fmt.Println()

	// Switch even/odd
	fmt.Println("4. Switch Even/Odd:")
	fmt.Printf("   SwitchEvenOdd(7) = %s\n", exercises.SwitchEvenOdd(7))
	fmt.Printf("   SwitchEvenOdd(4) = %s\n", exercises.SwitchEvenOdd(4))
	fmt.Println()

	// For loop sum
	fmt.Println("5. For Loop Sum:")
	fmt.Printf("   ForSum(5) = %d\n", exercises.ForSum(5))
	fmt.Printf("   ForSum(10) = %d\n", exercises.ForSum(10))
	fmt.Println()

	// For sum even
	fmt.Println("6. For Sum Even:")
	fmt.Printf("   ForSumEven(10) = %d\n", exercises.ForSumEven(10))
	fmt.Println()

	// Range count vowels
	fmt.Println("7. Range Count Vowels:")
	fmt.Printf("   RangeCountVowels(\"hello\") = %d\n", exercises.RangeCountVowels("hello"))
	fmt.Printf("   RangeCountVowels(\"AEIOU\") = %d\n", exercises.RangeCountVowels("AEIOU"))
	fmt.Println()

	// Skip multiples of three
	fmt.Println("8. Skip Multiples of Three:")
	fmt.Printf("   SumSkipMultiplesOfThree(10) = %d\n", exercises.SumSkipMultiplesOfThree(10))
	fmt.Println()

	// Defer
	fmt.Println("9. Defer:")
	fmt.Printf("   DeferGreeting() = %s\n", exercises.DeferGreeting())
	fmt.Println()

	// Two Sum
	fmt.Println("10. Two Sum:")
	fmt.Printf("   TwoSum([2,7,11,15], 9) = %s\n", exercises.TwoSum([]int{2, 7, 11, 15}, 9))
	fmt.Printf("   TwoSum([3,2,4], 6) = %s\n", exercises.TwoSum([]int{3, 2, 4}, 6))
	fmt.Println()

	fmt.Println("=== All exercises completed! ===")
}
