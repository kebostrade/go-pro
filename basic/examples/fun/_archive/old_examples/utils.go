//go:build ignore

package main

import "fmt"

func Square(x int) int {
	return x * x
}

func sumNums(n int, y int, z int) int {
	sum := n + y + z
	return Square(sum)
}

func main() {
	fmt.Println(Square(5))        // Output: 25
	fmt.Println(sumNums(1, 2, 3)) // Output: 36
}
