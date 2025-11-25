//go:build ignore

package main

import "fmt"

// Animal interface defines behavior for animals
type Animal interface {
	Speak() string
	GetName() string
}

// Dog implements Animal interface
type Dog struct {
	name string
}

func (d Dog) Speak() string {
	return "Woof!"
}

func (d Dog) GetName() string {
	return d.name
}

// Cat implements Animal interface
type Cat struct {
	name string
}

func (c Cat) Speak() string {
	return "Meow!"
}

func (c Cat) GetName() string {
	return c.name
}

func main() {
	var animals []Animal
	animals = append(animals, Dog{"Fido"})
	animals = append(animals, Cat{"Whiskers"})

	for _, animal := range animals {
		fmt.Printf("%s says: %s\n", animal.GetName(), animal.Speak())
	}
}
