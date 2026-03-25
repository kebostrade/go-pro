package exercises

// Exercise 1: Basic if/else
func IfElseBasic(n int) string {
	if n > 0 {
		return "positive"
	} else if n < 0 {
		return "negative"
	}
	return "zero"
}

// Exercise 2: If with initialization
func IfWithInit(n int) int {
	if n := n; n < 0 {
		return -n
	}
	return n
}

// Exercise 3: Switch basic
func SwitchDay(day int) string {
	switch day {
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	case 7:
		return "Sunday"
	default:
		return "invalid"
	}
}

// Exercise 4: Switch with expressions
func SwitchEvenOdd(n int) string {
	switch n % 2 {
	case 0:
		return "even"
	default:
		return "odd"
	}
}

// Exercise 5: For loop basic
func ForSum(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

// Exercise 6: For loop with condition
func ForSumEven(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		if i%2 == 0 {
			sum += i
		}
	}
	return sum
}

// Exercise 7: For range
func RangeCountVowels(s string) int {
	count := 0
	vowels := map[rune]bool{
		'a': true, 'e': true, 'i': true, 'o': true, 'u': true,
		'A': true, 'E': true, 'I': true, 'O': true, 'U': true,
	}
	for _, ch := range s {
		if vowels[ch] {
			count++
		}
	}
	return count
}

// Exercise 8: Break and Continue
func SumSkipMultiplesOfThree(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		if i%3 == 0 {
			continue
		}
		sum += i
	}
	return sum
}

// Exercise 9: Defer
func DeferGreeting() string {
	return "Hello"
}

package exercises

import "fmt"

// Exercise 10: Nested loops with label
func TwoSum(nums []int, target int) string {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return fmt.Sprintf("%d,%d", i, j)
			}
		}
	}
	return "-1,-1"
}
