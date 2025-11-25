//go:build ignore

package main

import "fmt"

func main() {
	i := 5
	if i < 5 {
		fmt.Println("i is less than 5")
	} else if i < 10 {
		fmt.Println("i is less then 10")
	} else {
		fmt.Println("i is at least 10")
	}
	fmt.Println("after the if statement")
}

func play() {
	player := Player{
		name:   "Jack",
		health: 100,
		power:  45.1,
	}
	fmt.Printf("this is the player: %+v\n", player)
}

var (
	floatVar32 float32 = 0.1
	floatVar64 float64 = 0.17
	name       string  = "Foo"
	intVar32   int32   = 1
	intVar64   int64   = 1343
	intVar     int     = -1
	unitVar    uint    = 1
	unit32Var  uint32  = 1
	unit64Var  uint64  = 10
	unit8Var   uint8   = 0x1
	byteVar    byte    = 0x2
)

type Player struct {
	name   string
	health int
	power  float64
}
