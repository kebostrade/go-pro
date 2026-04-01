// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go-pro-backend/internal/service"
)

// ExecuteCodeRequest represents a code execution request.
type ExecuteCodeRequest struct {
	Code      string             `json:"code"`
	TopicID   string             `json:"topic_id"`
	TestCases []service.TestCase `json:"test_cases"`
}

// ExecuteCode handles POST /api/execute
// Executes Go code against test cases with security validation and topic-specific package allowlists.
func (h *Handler) ExecuteCode(w http.ResponseWriter, r *http.Request) {
	var req ExecuteCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate code is not empty
	if req.Code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	// Validate test cases
	if len(req.TestCases) == 0 {
		http.Error(w, "At least one test case is required", http.StatusBadRequest)
		return
	}

	// Validate code has package main
	if !strings.Contains(req.Code, "package main") {
		http.Error(w, "Code must contain 'package main'", http.StatusBadRequest)
		return
	}

	// Validate code has func main()
	if !strings.Contains(req.Code, "func main()") {
		http.Error(w, "Code must contain 'func main()'", http.StatusBadRequest)
		return
	}

	// Check for blocked packages
	blocked := []string{"os", "net", "syscall", "unsafe", "runtime/debug"}
	for _, pkg := range blocked {
		if strings.Contains(req.Code, `"`+pkg+`"`) || strings.Contains(req.Code, "'"+pkg+"'") {
			http.Error(w, "Package '"+pkg+"' is not allowed", http.StatusBadRequest)
			return
		}
	}

	// Validate against topic-specific allowlist
	allowlist := getTopicAllowlist(req.TopicID)
	imports := extractImports(req.Code)
	for _, pkg := range imports {
		if !isAllowed(pkg, allowlist) {
			http.Error(w, "Package '"+pkg+"' is not allowed for this topic", http.StatusBadRequest)
			return
		}
	}

	// Execute using DockerExecutor with 30s context timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	result, err := h.services.Executor.ExecuteCode(ctx, &service.ExecuteRequest{
		Code:      req.Code,
		Language:  "go",
		TestCases: req.TestCases,
		Timeout:   15 * time.Second,
	})
	if err != nil {
		h.logger.Error(r.Context(), "Code execution failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write ExecuteResult as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error(r.Context(), "Failed to encode execution result", "error", err)
	}
}

// getTopicAllowlist returns topic-specific package allowlist.
func getTopicAllowlist(topicID string) []string {
	allowlists := map[string][]string{
		"grpc":       {"google.golang.org/grpc", "google.golang.org/protobuf"},
		"nats":       {"github.com/nats-io/nats.go"},
		"kafka":      {"github.com/IBM/sarama"},
		"prometheus": {"github.com/prometheus/client_golang"},
		// Add more topics as needed
	}
	return allowlists[topicID]
}

// extractImports extracts import paths from Go code.
func extractImports(code string) []string {
	var imports []string

	// Match single-line imports: import "fmt"
	singleImport := regexp.MustCompile(`import\s+"([^"]+)"`)
	for _, m := range singleImport.FindAllStringSubmatch(code, -1) {
		if len(m) > 1 {
			imports = append(imports, m[1])
		}
	}

	// Match multi-line imports: import ( "fmt" "os" )
	multiImport := regexp.MustCompile(`(?s)import\s+\(([^)]+)\)`)
	for _, m := range multiImport.FindAllStringSubmatch(code, -1) {
		if len(m) > 1 {
			// Extract individual imports from the group
			inner := m[1]
			for _, line := range strings.Split(inner, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, `"`) && strings.HasSuffix(line, `"`) {
					pkg := strings.Trim(line, `"`)
					imports = append(imports, pkg)
				}
			}
		}
	}

	return imports
}

// isAllowed checks if a package is allowed based on the allowlist.
func isAllowed(pkg string, allowlist []string) bool {
	if len(allowlist) == 0 {
		// Empty allowlist means only stdlib is allowed
		stdlib := []string{
			"fmt", "os", "io", "bufio", "strings", "strconv", "time", "math",
			"bytes", "encoding/json", "encoding/xml", "log", "sort", "container/list",
			"container/heap", "sync", "context", "crypto/aes", "crypto/cipher",
			"crypto/rand", "crypto/sha256", "hash", "regexp", "unicode", "unicode/utf8",
		}
		for _, p := range stdlib {
			if p == pkg {
				return true
			}
		}
		return false
	}
	for _, a := range allowlist {
		if a == pkg {
			return true
		}
	}
	return false
}
