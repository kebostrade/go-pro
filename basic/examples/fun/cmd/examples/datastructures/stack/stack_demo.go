package main

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/datastructures"
	"github.com/DimaJoyti/go-pro/basic/examples/fun/pkg/utils"
)

func main() {
	utils.PrintHeader("Stack (LIFO) Data Structure Demo")

	// Demo 1: Basic Stack Operations
	demo1BasicOperations()

	// Demo 2: String Stack
	demo2StringStack()

	// Demo 3: Stack Methods
	demo3StackMethods()

	// Demo 4: Practical Example - Undo/Redo
	demo4UndoRedo()

	// Demo 5: Balanced Parentheses
	demo5BalancedParentheses()
}

func demo1BasicOperations() {
	utils.PrintSubHeader("1. Basic Stack Operations")

	stack := datastructures.NewStack[int]()

	// Push elements
	fmt.Println("Pushing elements: 1, 2, 3, 4, 5")
	for i := 1; i <= 5; i++ {
		stack.Push(i)
		fmt.Printf("  Pushed %d, size: %d\n", i, stack.Size())
	}

	// Peek
	if top, err := stack.Peek(); err == nil {
		fmt.Printf("\nTop element (peek): %d\n", top)
	}

	// Pop elements
	fmt.Println("\nPopping elements:")
	for !stack.IsEmpty() {
		if val, err := stack.Pop(); err == nil {
			fmt.Printf("  Popped %d, remaining size: %d\n", val, stack.Size())
		}
	}

	// Try to pop from empty stack
	if _, err := stack.Pop(); err != nil {
		fmt.Printf("\nError popping from empty stack: %v\n", err)
	}
}

func demo2StringStack() {
	utils.PrintSubHeader("2. String Stack - Reverse Words")

	stack := datastructures.NewStack[string]()

	words := []string{"Hello", "world", "from", "Go"}
	fmt.Printf("Original words: %v\n", words)

	// Push all words
	for _, word := range words {
		stack.Push(word)
	}

	// Pop to reverse
	fmt.Print("Reversed: ")
	for !stack.IsEmpty() {
		if word, err := stack.Pop(); err == nil {
			fmt.Print(word, " ")
		}
	}
	fmt.Println()
}

func demo3StackMethods() {
	utils.PrintSubHeader("3. Stack Methods")

	stack := datastructures.NewStack[int]()

	// Add elements
	for i := 1; i <= 5; i++ {
		stack.Push(i * 10)
	}

	fmt.Printf("Stack: %s\n", stack.String())
	fmt.Printf("Size: %d\n", stack.Size())

	// ToSlice
	slice := stack.ToSlice()
	fmt.Printf("As slice (top to bottom): %v\n", slice)

	// Clone
	cloned := stack.Clone()
	fmt.Printf("Cloned stack: %s\n", cloned.String())

	// Filter - keep only values > 25
	filtered := stack.Filter(func(val int) bool {
		return val > 25
	})
	fmt.Printf("Filtered (>25): %s\n", filtered.String())

	// Map - double all values
	doubled := stack.Map(func(val int) int {
		return val * 2
	})
	fmt.Printf("Doubled: %s\n", doubled.String())

	// ForEach
	fmt.Print("ForEach (top to bottom): ")
	stack.ForEach(func(val int) {
		fmt.Printf("%d ", val)
	})
	fmt.Println()

	// Reverse
	stack.Reverse()
	fmt.Printf("After reverse: %s\n", stack.String())

	// Clear
	stack.Clear()
	fmt.Printf("After clear: %s (isEmpty: %v)\n", stack.String(), stack.IsEmpty())
}

func demo4UndoRedo() {
	utils.PrintSubHeader("4. Practical Example - Undo/Redo System")

	type Action struct {
		Type  string
		Value string
	}

	undoStack := datastructures.NewStack[Action]()
	redoStack := datastructures.NewStack[Action]()

	// Perform actions
	actions := []Action{
		{Type: "INSERT", Value: "Hello"},
		{Type: "INSERT", Value: " World"},
		{Type: "DELETE", Value: " World"},
		{Type: "INSERT", Value: " Go"},
	}

	fmt.Println("Performing actions:")
	for _, action := range actions {
		fmt.Printf("  %s: %s\n", action.Type, action.Value)
		undoStack.Push(action)
	}

	// Undo last 2 actions
	fmt.Println("\nUndo last 2 actions:")
	for i := 0; i < 2; i++ {
		if action, err := undoStack.Pop(); err == nil {
			fmt.Printf("  Undoing %s: %s\n", action.Type, action.Value)
			redoStack.Push(action)
		}
	}

	// Redo 1 action
	fmt.Println("\nRedo 1 action:")
	if action, err := redoStack.Pop(); err == nil {
		fmt.Printf("  Redoing %s: %s\n", action.Type, action.Value)
		undoStack.Push(action)
	}

	fmt.Printf("\nUndo stack size: %d\n", undoStack.Size())
	fmt.Printf("Redo stack size: %d\n", redoStack.Size())
}

func demo5BalancedParentheses() {
	utils.PrintSubHeader("5. Balanced Parentheses Checker")

	testCases := []string{
		"()",
		"()[]{}",
		"(]",
		"([)]",
		"{[]}",
		"((()))",
		"((())",
	}

	for _, test := range testCases {
		result := isBalanced(test)
		status := "✓"
		if !result {
			status = "✗"
		}
		fmt.Printf("%s %s: %v\n", status, test, result)
	}
}

func isBalanced(s string) bool {
	stack := datastructures.NewStack[rune]()
	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range s {
		switch char {
		case '(', '[', '{':
			stack.Push(char)
		case ')', ']', '}':
			if stack.IsEmpty() {
				return false
			}
			if top, err := stack.Pop(); err != nil || top != pairs[char] {
				return false
			}
		}
	}

	return stack.IsEmpty()
}
