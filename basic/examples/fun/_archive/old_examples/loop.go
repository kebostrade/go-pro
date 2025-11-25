//go:build ignore

package main

import "fmt"

func main() {

	// // Simple iteration over a range
	// for i := 0; i < 5; i++ {
	// 	fmt.Println(i)
	// }

	// // iterate over collection
	// numbers := []int{1, 2, 3, 4, 5}
	// for index, value := range numbers {
	// 	fmt.Printf("index: %d, value: %d\n", index, value)
	// }

	// for i:=0; i<=15; i++ {
	// 	if i%2 == 0 {
	// 		continue
	// 	}
	// 	fmt.Println(i)
	// 	if i == 5 {
	// 		break
	// 	}
	// }

	rows := 5

	for i := 1; i <= rows; i++ {
		// inner loop for spaces before stars
		for j := 1; j <= rows-i; j++ {
			fmt.Print(" ")
		}
		// inner loop for stars
		for k := 1; k <= 2*i-1; k++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}
