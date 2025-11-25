package golang

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// Analyzer provides Go-specific code analysis
type Analyzer struct {
	language types.Language
}

// NewAnalyzer creates a new Go analyzer
func NewAnalyzer() *Analyzer {
	lang, _ := types.GetLanguage("go")
	return &Analyzer{
		language: lang,
	}
}

// GetLanguage returns the Go language
func (a *Analyzer) GetLanguage() types.Language {
	return a.language
}

// Analyze performs static analysis on Go code
func (a *Analyzer) Analyze(ctx context.Context, code string) (*types.CodeAnalysis, error) {
	analysis := &types.CodeAnalysis{
		Issues:            make([]types.CodeIssue, 0),
		SecurityIssues:    make([]types.SecurityIssue, 0),
		PerformanceIssues: make([]types.PerformanceIssue, 0),
		BestPractices:     make([]types.BestPracticeViolation, 0),
	}

	// Parse the code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "code.go", code, parser.ParseComments)
	if err != nil {
		// Syntax error
		analysis.Issues = append(analysis.Issues, types.CodeIssue{
			Severity:    types.SeverityError,
			Type:        types.IssueTypeSyntax,
			Message:     fmt.Sprintf("Syntax error: %v", err),
			Suggestion:  "Fix the syntax error before proceeding",
			CodeSnippet: code,
		})
		analysis.OverallScore = 0
		return analysis, nil
	}

	// Calculate metrics
	analysis.Metrics = a.calculateMetrics(code, file)

	// Check for common issues
	a.checkCommonIssues(file, fset, analysis)

	// Check for security issues
	a.checkSecurityIssues(file, analysis)

	// Check for performance issues
	a.checkPerformanceIssues(file, analysis)

	// Check for best practice violations
	a.checkBestPractices(file, analysis)

	// Calculate overall score
	analysis.OverallScore = a.calculateScore(analysis)

	return analysis, nil
}

// Lint runs linting on Go code
func (a *Analyzer) Lint(ctx context.Context, code string) ([]types.CodeIssue, error) {
	issues := make([]types.CodeIssue, 0)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "code.go", code, parser.ParseComments)
	if err != nil {
		issues = append(issues, types.CodeIssue{
			Severity: types.SeverityError,
			Type:     types.IssueTypeSyntax,
			Message:  fmt.Sprintf("Syntax error: %v", err),
		})
		return issues, nil
	}

	// Check for unused imports
	a.checkUnusedImports(file, &issues)

	// Check for naming conventions
	a.checkNamingConventions(file, &issues)

	// Check for error handling
	a.checkErrorHandling(file, &issues)

	return issues, nil
}

// Format formats Go code according to gofmt standards
func (a *Analyzer) Format(ctx context.Context, code string) (string, error) {
	// In production, use go/format package
	// For now, return as-is
	return code, nil
}

// Validate checks if Go code is syntactically valid
func (a *Analyzer) Validate(ctx context.Context, code string) error {
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "code.go", code, parser.ParseComments)
	if err != nil {
		return &types.ToolError{
			Code:    "INVALID_SYNTAX",
			Message: fmt.Sprintf("Invalid Go syntax: %v", err),
		}
	}
	return nil
}

// ExtractImports extracts import statements from Go code
func (a *Analyzer) ExtractImports(ctx context.Context, code string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "code.go", code, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	imports := make([]string, 0)
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		imports = append(imports, path)
	}

	return imports, nil
}

// GetComplexity calculates cyclomatic complexity
func (a *Analyzer) GetComplexity(ctx context.Context, code string) (int, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "code.go", code, 0)
	if err != nil {
		return 0, err
	}

	complexity := 1 // Base complexity

	ast.Inspect(file, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			complexity++
		}
		return true
	})

	return complexity, nil
}

// calculateMetrics calculates code quality metrics
func (a *Analyzer) calculateMetrics(code string, file *ast.File) types.CodeMetrics {
	lines := strings.Split(code, "\n")
	loc := len(lines)

	// Count comments
	commentLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			commentLines++
		}
	}

	commentRatio := 0.0
	if loc > 0 {
		commentRatio = float64(commentLines) / float64(loc)
	}

	// Calculate complexity
	complexity := 1
	ast.Inspect(file, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause:
			complexity++
		}
		return true
	})

	// Calculate maintainability (simplified)
	maintainability := 100
	if complexity > 10 {
		maintainability -= (complexity - 10) * 5
	}
	if maintainability < 0 {
		maintainability = 0
	}

	return types.CodeMetrics{
		LinesOfCode:     loc,
		Complexity:      complexity,
		Maintainability: maintainability,
		CommentRatio:    commentRatio,
	}
}

// checkCommonIssues checks for common coding issues
func (a *Analyzer) checkCommonIssues(file *ast.File, fset *token.FileSet, analysis *types.CodeAnalysis) {
	ast.Inspect(file, func(n ast.Node) bool {
		// Check for empty if statements
		if ifStmt, ok := n.(*ast.IfStmt); ok {
			if ifStmt.Body == nil || len(ifStmt.Body.List) == 0 {
				pos := fset.Position(ifStmt.Pos())
				analysis.Issues = append(analysis.Issues, types.CodeIssue{
					Severity:   types.SeverityWarning,
					Type:       types.IssueTypeLogic,
					Message:    "Empty if statement",
					Line:       pos.Line,
					Suggestion: "Remove empty if statement or add implementation",
				})
			}
		}
		return true
	})
}

// checkSecurityIssues checks for security vulnerabilities
func (a *Analyzer) checkSecurityIssues(file *ast.File, analysis *types.CodeAnalysis) {
	ast.Inspect(file, func(n ast.Node) bool {
		// Check for SQL injection risks
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Exec" || sel.Sel.Name == "Query" {
					analysis.SecurityIssues = append(analysis.SecurityIssues, types.SecurityIssue{
						Type:           "SQL Injection Risk",
						Severity:       types.SeverityWarning,
						Description:    "Potential SQL injection vulnerability",
						Recommendation: "Use parameterized queries instead of string concatenation",
						CWE:            "CWE-89",
					})
				}
			}
		}
		return true
	})
}

// checkPerformanceIssues checks for performance problems
func (a *Analyzer) checkPerformanceIssues(file *ast.File, analysis *types.CodeAnalysis) {
	// Check for string concatenation in loops
	ast.Inspect(file, func(n ast.Node) bool {
		if forStmt, ok := n.(*ast.ForStmt); ok {
			ast.Inspect(forStmt.Body, func(inner ast.Node) bool {
				if binExpr, ok := inner.(*ast.BinaryExpr); ok {
					if binExpr.Op == token.ADD {
						analysis.PerformanceIssues = append(analysis.PerformanceIssues, types.PerformanceIssue{
							Type:        "String Concatenation in Loop",
							Description: "String concatenation in loop can be inefficient",
							Impact:      "O(n²) time complexity",
							Suggestion:  "Use strings.Builder for efficient string concatenation",
						})
					}
				}
				return true
			})
		}
		return true
	})
}

// checkBestPractices checks for best practice violations
func (a *Analyzer) checkBestPractices(file *ast.File, analysis *types.CodeAnalysis) {
	// Check for exported functions without comments
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.IsExported() && funcDecl.Doc == nil {
				analysis.BestPractices = append(analysis.BestPractices, types.BestPracticeViolation{
					Practice:       "Document Exported Functions",
					Description:    fmt.Sprintf("Exported function %s lacks documentation", funcDecl.Name.Name),
					Recommendation: "Add a comment describing what the function does",
				})
			}
		}
	}
}

// checkUnusedImports checks for unused imports
func (a *Analyzer) checkUnusedImports(file *ast.File, issues *[]types.CodeIssue) {
	// Simplified check - in production use go/types
	for _, imp := range file.Imports {
		*issues = append(*issues, types.CodeIssue{
			Severity: types.SeverityInfo,
			Type:     types.IssueTypeStyle,
			Message:  fmt.Sprintf("Check if import %s is used", imp.Path.Value),
		})
	}
}

// checkNamingConventions checks naming conventions
func (a *Analyzer) checkNamingConventions(file *ast.File, issues *[]types.CodeIssue) {
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			name := funcDecl.Name.Name
			if strings.Contains(name, "_") && funcDecl.Name.IsExported() {
				*issues = append(*issues, types.CodeIssue{
					Severity:   types.SeverityWarning,
					Type:       types.IssueTypeStyle,
					Message:    fmt.Sprintf("Exported function %s uses underscores", name),
					Suggestion: "Use camelCase for exported functions",
				})
			}
		}
		return true
	})
}

// checkErrorHandling checks error handling patterns
func (a *Analyzer) checkErrorHandling(file *ast.File, issues *[]types.CodeIssue) {
	// Check for ignored errors
	ast.Inspect(file, func(n ast.Node) bool {
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			for _, lhs := range assignStmt.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if ident.Name == "_" {
						*issues = append(*issues, types.CodeIssue{
							Severity:   types.SeverityWarning,
							Type:       types.IssueTypeBestPractice,
							Message:    "Error is being ignored",
							Suggestion: "Handle errors explicitly",
						})
					}
				}
			}
		}
		return true
	})
}

// calculateScore calculates overall code quality score
func (a *Analyzer) calculateScore(analysis *types.CodeAnalysis) int {
	score := 100

	// Deduct points for issues
	for _, issue := range analysis.Issues {
		switch issue.Severity {
		case types.SeverityError:
			score -= 20
		case types.SeverityWarning:
			score -= 10
		case types.SeverityInfo:
			score -= 2
		}
	}

	// Deduct points for security issues
	score -= len(analysis.SecurityIssues) * 15

	// Deduct points for performance issues
	score -= len(analysis.PerformanceIssues) * 5

	// Deduct points for best practice violations
	score -= len(analysis.BestPractices) * 3

	if score < 0 {
		score = 0
	}

	return score
}

