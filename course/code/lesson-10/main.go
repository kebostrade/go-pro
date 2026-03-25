package main

import (
	"fmt"
	"lesson-10/exercises"
)

func main() {
	fmt.Println("=== Lesson 10: Packages and Modules ===")
	fmt.Println()

	// Package initialization
	fmt.Println("1. Package Initialization:")
	fmt.Printf("   Initialized: %v\n", exercises.IsInitialized())
	fmt.Printf("   Init count: %d\n", exercises.GetInitCount())
	fmt.Println()

	// Import aliases
	fmt.Println("2. Import Aliases:")
	result := exercises.UseAliases()
	fmt.Printf("   UseAliases: '%s'\n", result)
	fmt.Println()

	// External packages
	fmt.Println("3. External Packages:")
	_ = exercises.UseExternal()
	fmt.Println("   External packages work!")
	fmt.Println()

	// Package documentation
	fmt.Println("4. Package Documentation:")
	pd := exercises.NewPackageDocumentation("mymath", "Math utilities", "v2.0.0")
	fmt.Printf("%s\n", pd)
	fmt.Println()

	// Internal packages
	fmt.Println("5. Internal Packages:")
	fmt.Println("   Internal packages (named 'internal') can only")
	fmt.Println("   be imported by parent directory packages")
	fmt.Println()

	// Vendor directory
	fmt.Println("6. Vendor Directory:")
	fmt.Printf("   %s\n", exercises.VendorDirectory())
	fmt.Println()

	// Go mod commands
	fmt.Println("7. Go Module Commands:")
	fmt.Printf("   %s\n", exercises.ModuleCommands())
	fmt.Println()

	// Best practices
	fmt.Println("8. Best Practices:")
	fmt.Println("   - Use meaningful package names")
	fmt.Println("   - Keep packages focused and small")
	fmt.Println("   - Use 'internal' for private packages")
	fmt.Println("   - Document exported functions")
	fmt.Println("   - Use go mod tidy regularly")
	fmt.Println()

	fmt.Println("=== All exercises completed! ===")
}
