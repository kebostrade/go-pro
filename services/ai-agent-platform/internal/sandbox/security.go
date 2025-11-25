package sandbox

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// SecurityPolicy defines security rules for code execution
type SecurityPolicy struct {
	// AllowedImports list of allowed import patterns
	AllowedImports []string

	// BlockedImports list of blocked import patterns
	BlockedImports []string

	// BlockedPatterns dangerous code patterns to block
	BlockedPatterns []CodePattern

	// MaxCodeSize maximum code size in bytes
	MaxCodeSize int

	// AllowNetworkAccess whether network access is allowed
	AllowNetworkAccess bool

	// AllowFileSystemAccess whether file system access is allowed
	AllowFileSystemAccess bool

	// AllowExternalCommands whether external commands are allowed
	AllowExternalCommands bool
}

// CodePattern represents a dangerous code pattern
type CodePattern struct {
	// Pattern regex pattern to match
	Pattern string

	// Description what this pattern does
	Description string

	// Severity how dangerous this is
	Severity string

	// Recommendation what to do instead
	Recommendation string
}

// SecurityValidator validates code against security policies
type SecurityValidator struct {
	policies map[string]*SecurityPolicy
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator() *SecurityValidator {
	return &SecurityValidator{
		policies: make(map[string]*SecurityPolicy),
	}
}

// RegisterPolicy registers a security policy for a language
func (v *SecurityValidator) RegisterPolicy(language string, policy *SecurityPolicy) {
	v.policies[language] = policy
}

// Validate validates code against security policies
func (v *SecurityValidator) Validate(language, code string) error {
	policy, ok := v.policies[language]
	if !ok {
		// Use default policy if language-specific not found
		policy = v.getDefaultPolicy()
	}

	// Check code size
	if policy.MaxCodeSize > 0 && len(code) > policy.MaxCodeSize {
		return &types.ToolError{
			Code:    "CODE_TOO_LARGE",
			Message: fmt.Sprintf("Code size exceeds limit: %d > %d", len(code), policy.MaxCodeSize),
		}
	}

	// Check for blocked patterns
	for _, pattern := range policy.BlockedPatterns {
		matched, err := regexp.MatchString(pattern.Pattern, code)
		if err != nil {
			continue
		}
		if matched {
			return &types.ToolError{
				Code:    "UNSAFE_CODE_PATTERN",
				Message: fmt.Sprintf("Blocked pattern detected: %s", pattern.Description),
				Details: pattern.Recommendation,
			}
		}
	}

	// Check imports (language-specific)
	if err := v.validateImports(language, code, policy); err != nil {
		return err
	}

	return nil
}

// validateImports checks if imports are allowed
func (v *SecurityValidator) validateImports(language, code string, policy *SecurityPolicy) error {
	switch language {
	case "go":
		return v.validateGoImports(code, policy)
	case "python":
		return v.validatePythonImports(code, policy)
	case "javascript", "typescript":
		return v.validateJSImports(code, policy)
	default:
		return nil
	}
}

// validateGoImports validates Go imports
func (v *SecurityValidator) validateGoImports(code string, policy *SecurityPolicy) error {
	// Extract import statements
	importRegex := regexp.MustCompile(`import\s+(?:\(([^)]+)\)|"([^"]+)")`)
	matches := importRegex.FindAllStringSubmatch(code, -1)

	for _, match := range matches {
		var importPath string
		if match[1] != "" {
			// Multiple imports
			imports := strings.Split(match[1], "\n")
			for _, imp := range imports {
				imp = strings.TrimSpace(imp)
				imp = strings.Trim(imp, `"`)
				if imp != "" {
					if err := v.checkImport(imp, policy); err != nil {
						return err
					}
				}
			}
		} else if match[2] != "" {
			// Single import
			importPath = match[2]
			if err := v.checkImport(importPath, policy); err != nil {
				return err
			}
		}
	}

	return nil
}

// validatePythonImports validates Python imports
func (v *SecurityValidator) validatePythonImports(code string, policy *SecurityPolicy) error {
	// Extract import statements
	importRegex := regexp.MustCompile(`(?:from\s+(\S+)\s+)?import\s+(\S+)`)
	matches := importRegex.FindAllStringSubmatch(code, -1)

	for _, match := range matches {
		module := match[1]
		if module == "" {
			module = match[2]
		}
		if err := v.checkImport(module, policy); err != nil {
			return err
		}
	}

	return nil
}

// validateJSImports validates JavaScript/TypeScript imports
func (v *SecurityValidator) validateJSImports(code string, policy *SecurityPolicy) error {
	// Extract import statements
	importRegex := regexp.MustCompile(`(?:import|require)\s*\(?['"]([^'"]+)['"]`)
	matches := importRegex.FindAllStringSubmatch(code, -1)

	for _, match := range matches {
		module := match[1]
		if err := v.checkImport(module, policy); err != nil {
			return err
		}
	}

	return nil
}

// checkImport checks if an import is allowed
func (v *SecurityValidator) checkImport(importPath string, policy *SecurityPolicy) error {
	// Check blocked imports first
	for _, blocked := range policy.BlockedImports {
		if strings.Contains(importPath, blocked) {
			return &types.ToolError{
				Code:    "BLOCKED_IMPORT",
				Message: fmt.Sprintf("Import not allowed: %s", importPath),
				Details: "This import is blocked for security reasons",
			}
		}
	}

	// If allowed imports are specified, check against them
	if len(policy.AllowedImports) > 0 {
		allowed := false
		for _, allowedPattern := range policy.AllowedImports {
			if strings.Contains(importPath, allowedPattern) {
				allowed = true
				break
			}
		}
		if !allowed {
			return &types.ToolError{
				Code:    "IMPORT_NOT_ALLOWED",
				Message: fmt.Sprintf("Import not in allowed list: %s", importPath),
			}
		}
	}

	return nil
}

// getDefaultPolicy returns the default security policy
func (v *SecurityValidator) getDefaultPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		BlockedImports: []string{
			"os/exec",
			"syscall",
			"unsafe",
			"net/http",
			"net",
			"subprocess",
			"socket",
			"eval",
			"exec",
		},
		BlockedPatterns: []CodePattern{
			{
				Pattern:        `os\.RemoveAll|os\.Remove`,
				Description:    "File deletion operations",
				Severity:       "high",
				Recommendation: "File deletion is not allowed in sandbox",
			},
			{
				Pattern:        `exec\.Command|subprocess\.`,
				Description:    "External command execution",
				Severity:       "critical",
				Recommendation: "External commands are not allowed",
			},
			{
				Pattern:        `eval\(|exec\(`,
				Description:    "Dynamic code execution",
				Severity:       "critical",
				Recommendation: "Dynamic code execution is not allowed",
			},
			{
				Pattern:        `__import__|importlib`,
				Description:    "Dynamic imports",
				Severity:       "high",
				Recommendation: "Dynamic imports are restricted",
			},
		},
		MaxCodeSize:            100 * 1024, // 100KB
		AllowNetworkAccess:     false,
		AllowFileSystemAccess:  false,
		AllowExternalCommands:  false,
	}
}

// GetGoPolicy returns security policy for Go
func GetGoPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		AllowedImports: []string{
			"fmt",
			"strings",
			"strconv",
			"math",
			"sort",
			"time",
			"encoding/json",
			"errors",
		},
		BlockedImports: []string{
			"os/exec",
			"syscall",
			"unsafe",
			"net/http",
			"net",
			"plugin",
		},
		BlockedPatterns: []CodePattern{
			{
				Pattern:        `os\.RemoveAll|os\.Remove`,
				Description:    "File deletion",
				Severity:       "high",
				Recommendation: "File operations are restricted",
			},
			{
				Pattern:        `exec\.Command`,
				Description:    "Command execution",
				Severity:       "critical",
				Recommendation: "External commands not allowed",
			},
		},
		MaxCodeSize:            50 * 1024,
		AllowNetworkAccess:     false,
		AllowFileSystemAccess:  false,
		AllowExternalCommands:  false,
	}
}

// GetPythonPolicy returns security policy for Python
func GetPythonPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		AllowedImports: []string{
			"math",
			"random",
			"datetime",
			"json",
			"re",
			"collections",
			"itertools",
		},
		BlockedImports: []string{
			"os",
			"sys",
			"subprocess",
			"socket",
			"urllib",
			"requests",
			"eval",
			"exec",
			"__import__",
		},
		BlockedPatterns: []CodePattern{
			{
				Pattern:        `eval\(|exec\(`,
				Description:    "Dynamic code execution",
				Severity:       "critical",
				Recommendation: "eval/exec not allowed",
			},
			{
				Pattern:        `__import__`,
				Description:    "Dynamic imports",
				Severity:       "high",
				Recommendation: "Dynamic imports restricted",
			},
		},
		MaxCodeSize:            50 * 1024,
		AllowNetworkAccess:     false,
		AllowFileSystemAccess:  false,
		AllowExternalCommands:  false,
	}
}

// GetJavaScriptPolicy returns security policy for JavaScript
func GetJavaScriptPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		BlockedImports: []string{
			"child_process",
			"fs",
			"net",
			"http",
			"https",
			"eval",
			"vm",
		},
		BlockedPatterns: []CodePattern{
			{
				Pattern:        `eval\(|Function\(`,
				Description:    "Dynamic code execution",
				Severity:       "critical",
				Recommendation: "eval/Function not allowed",
			},
			{
				Pattern:        `require\(['"]child_process['"]\)`,
				Description:    "Process spawning",
				Severity:       "critical",
				Recommendation: "Process operations not allowed",
			},
		},
		MaxCodeSize:            50 * 1024,
		AllowNetworkAccess:     false,
		AllowFileSystemAccess:  false,
		AllowExternalCommands:  false,
	}
}

