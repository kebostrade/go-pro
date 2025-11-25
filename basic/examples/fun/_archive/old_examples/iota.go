//go:build ignore

package main

import "fmt"

const (
	Readable   = 1 << iota // 1 << 0 = 001
	Writable               // 1 << 1 = 010
	Executable             // 1 << 2 = 100
)

const (
	Monday = iota + 1
	_      // Skip 0
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	fmt.Println("Days of the week using iota:")
	fmt.Println("Monday:", Monday)
	fmt.Println("Tuesday:", Tuesday)
	fmt.Println("Wednesday:", Wednesday)
	fmt.Println("Thursday:", Thursday)
	fmt.Println("Friday:", Friday)
	fmt.Println("Saturday:", Saturday)
	fmt.Println("Sunday:", Sunday)
}
