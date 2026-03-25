package exercises

// Exercise 1: Package structure
// Create a simple package with exported and unexported functions
// Note: In Go, unexported (lowercase) functions are private

// Exercise 2: Import with alias
// Demonstrate importing packages with aliases

// Exercise 3: Init function
// Create a package with init function that sets up state

// Exercise 4: Package-level variables
// Demonstrate package-level variables and initialization

// Exercise 5: Go modules
// Demonstrate module usage with go.mod

// Exercise 6: Import external package
// Use a simple external package like "strings" or "time"

// Exercise 7: Blank identifier imports
// Use _ to import for side effects only

// Exercise 8: Package documentation
// Write documentation for a package

// Exercise 9: Internal packages
// Create an internal package structure

// Exercise 10: Vendor directory concept
// Explain vendor directory usage

// ============================================
// Let's implement some of these concepts
// ============================================

import (
	"fmt"
	"strings"
)

// PackageVariables demonstrates package-level variables
var (
	initialized bool
	initCount   int
)

// InitFunction demonstrates init
func init() {
	initialized = true
	initCount++
	fmt.Println("Package initialized!")
}

// IsInitialized returns package initialization status
func IsInitialized() bool {
	return initialized
}

// GetInitCount returns how many times init was called
func GetInitCount() int {
	return initCount
}

// Exercise 2: Using import aliases
// Example: import f "fmt" or import "strings"

func UseAliases() string {
	// Using strings with alias (if imported as str)
	// str.TrimSpace("  hello  ")
	return strings.TrimSpace("  hello  ")
}

// Exercise 6: Using external packages
func UseExternal() string {
	// Using time package
	return "Time package available"
}

// PackageDocumentation represents a documented package
type PackageDocumentation struct {
	Name        string
	Version string
	Description     string
}

// String returns formatted documentation
func (pd PackageDocumentation) String() string {
	return fmt.Sprintf("# %s\n\n%s\n\nVersion: %s", pd.Name, pd.Description, pd.Version)
}

// NewPackageDocumentation creates a new package documentation
func NewPackageDocumentation(name, desc, version string) PackageDocumentation {
	return PackageDocumentation{
		Name:        name,
		Description: desc,
		Version:     version,
	}
}

// InternalPackageExample demonstrates internal package pattern
// In Go 1.4+, packages named "internal" can only be imported by parent directories

// BlankImport demonstrates blank identifier import
// Usage: import _ "package/name" (for side effects only, like init())
func BlankImport() string {
	// This would be used when you only want the init() function
	// from a package to run, without using any of its functions
	return "Blank import would trigger init()"
}

// VendorDirectory explains vendor concept
/*
In Go, vendor/ directory contains package dependencies.
Usage:
- Place dependencies in vendor/ folder
- go build -mod=vendor uses vendor instead of downloading
*/
func VendorDirectory() string {
	return "Use vendor/ directory to freeze dependencies"
}

// ModuleCommands shows common go mod commands
/*
go mod init <module-name>
go mod tidy
go mod download
go list -m all
go get <package>
go mod why <package>
*/
func ModuleCommands() string {
	return "See comments for module commands"
}
