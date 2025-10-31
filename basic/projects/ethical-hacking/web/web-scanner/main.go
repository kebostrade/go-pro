package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ⚠️ LEGAL WARNING ⚠️
// This tool is for EDUCATIONAL PURPOSES ONLY.
// Only scan websites you own or have explicit permission to test.

type Vulnerability struct {
	Type        string
	Severity    string
	Description string
	URL         string
}

func main() {
	printBanner()

	// Parse flags
	targetURL := flag.String("url", "", "Target URL to scan")
	timeout := flag.Int("timeout", 10, "Request timeout in seconds")
	flag.Parse()

	if *targetURL == "" {
		fmt.Println("❌ Error: URL is required")
		fmt.Println("Usage: web-scanner -url https://example.com")
		return
	}

	fmt.Printf("🎯 Target: %s\n", *targetURL)
	fmt.Printf("⏱️  Timeout: %ds\n\n", *timeout)

	// Create HTTP client
	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var vulnerabilities []Vulnerability

	// Run security checks
	fmt.Println("🔍 Running security checks...\n")

	vulnerabilities = append(vulnerabilities, checkHTTPSRedirect(*targetURL, client)...)
	vulnerabilities = append(vulnerabilities, checkSecurityHeaders(*targetURL, client)...)
	vulnerabilities = append(vulnerabilities, checkXSSVulnerability(*targetURL, client)...)
	vulnerabilities = append(vulnerabilities, checkSQLInjection(*targetURL, client)...)
	vulnerabilities = append(vulnerabilities, checkDirectoryListing(*targetURL, client)...)

	// Display results
	displayResults(vulnerabilities)
}

func checkHTTPSRedirect(targetURL string, client *http.Client) []Vulnerability {
	var vulns []Vulnerability

	fmt.Println("📋 Checking HTTPS redirect...")

	// Check if HTTP redirects to HTTPS
	if strings.HasPrefix(targetURL, "https://") {
		httpURL := strings.Replace(targetURL, "https://", "http://", 1)
		resp, err := client.Get(httpURL)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusMovedPermanently && resp.StatusCode != http.StatusFound {
				vulns = append(vulns, Vulnerability{
					Type:        "Missing HTTPS Redirect",
					Severity:    "Medium",
					Description: "HTTP does not redirect to HTTPS",
					URL:         httpURL,
				})
			}
		}
	}

	return vulns
}

func checkSecurityHeaders(targetURL string, client *http.Client) []Vulnerability {
	var vulns []Vulnerability

	fmt.Println("📋 Checking security headers...")

	resp, err := client.Get(targetURL)
	if err != nil {
		return vulns
	}
	defer resp.Body.Close()

	// Check for important security headers
	headers := map[string]string{
		"X-Frame-Options":           "Clickjacking protection",
		"X-Content-Type-Options":    "MIME-sniffing protection",
		"Strict-Transport-Security": "HSTS enforcement",
		"Content-Security-Policy":   "XSS protection",
		"X-XSS-Protection":          "XSS filter",
	}

	for header, description := range headers {
		if resp.Header.Get(header) == "" {
			vulns = append(vulns, Vulnerability{
				Type:        "Missing Security Header",
				Severity:    "Low",
				Description: fmt.Sprintf("Missing %s (%s)", header, description),
				URL:         targetURL,
			})
		}
	}

	return vulns
}

func checkXSSVulnerability(targetURL string, client *http.Client) []Vulnerability {
	var vulns []Vulnerability

	fmt.Println("📋 Checking for XSS vulnerabilities...")

	// Parse URL
	u, err := url.Parse(targetURL)
	if err != nil {
		return vulns
	}

	// Test XSS payloads
	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
	}

	for _, payload := range xssPayloads {
		// Add payload to query parameter
		q := u.Query()
		q.Set("test", payload)
		u.RawQuery = q.Encode()

		resp, err := client.Get(u.String())
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		// Check if payload is reflected in response
		if strings.Contains(string(body), payload) {
			vulns = append(vulns, Vulnerability{
				Type:        "Potential XSS",
				Severity:    "High",
				Description: "User input reflected in response without sanitization",
				URL:         u.String(),
			})
			break
		}
	}

	return vulns
}

func checkSQLInjection(targetURL string, client *http.Client) []Vulnerability {
	var vulns []Vulnerability

	fmt.Println("📋 Checking for SQL injection vulnerabilities...")

	// Parse URL
	u, err := url.Parse(targetURL)
	if err != nil {
		return vulns
	}

	// Test SQL injection payloads
	sqlPayloads := []string{
		"' OR '1'='1",
		"1' OR '1'='1' --",
		"' UNION SELECT NULL--",
	}

	for _, payload := range sqlPayloads {
		q := u.Query()
		q.Set("id", payload)
		u.RawQuery = q.Encode()

		resp, err := client.Get(u.String())
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		// Check for SQL error messages
		sqlErrors := []string{
			"SQL syntax",
			"mysql_fetch",
			"ORA-",
			"PostgreSQL",
			"SQLite",
		}

		for _, sqlError := range sqlErrors {
			if strings.Contains(string(body), sqlError) {
				vulns = append(vulns, Vulnerability{
					Type:        "Potential SQL Injection",
					Severity:    "Critical",
					Description: "SQL error message detected in response",
					URL:         u.String(),
				})
				return vulns
			}
		}
	}

	return vulns
}

func checkDirectoryListing(targetURL string, client *http.Client) []Vulnerability {
	var vulns []Vulnerability

	fmt.Println("📋 Checking for directory listing...")

	// Common directories to check
	directories := []string{
		"/admin/",
		"/backup/",
		"/config/",
		"/uploads/",
		"/files/",
	}

	u, err := url.Parse(targetURL)
	if err != nil {
		return vulns
	}

	for _, dir := range directories {
		u.Path = dir
		resp, err := client.Get(u.String())
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		// Check for directory listing indicators
		if strings.Contains(string(body), "Index of") || strings.Contains(string(body), "Parent Directory") {
			vulns = append(vulns, Vulnerability{
				Type:        "Directory Listing",
				Severity:    "Medium",
				Description: "Directory listing enabled",
				URL:         u.String(),
			})
		}
	}

	return vulns
}

func displayResults(vulnerabilities []Vulnerability) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 SCAN RESULTS")
	fmt.Println(strings.Repeat("=", 80))

	if len(vulnerabilities) == 0 {
		fmt.Println("\n✅ No vulnerabilities found!")
		fmt.Println("\nNote: This is a basic scanner. Consider using professional tools for")
		fmt.Println("comprehensive security testing.")
		return
	}

	// Count by severity
	critical := 0
	high := 0
	medium := 0
	low := 0

	for _, vuln := range vulnerabilities {
		switch vuln.Severity {
		case "Critical":
			critical++
		case "High":
			high++
		case "Medium":
			medium++
		case "Low":
			low++
		}
	}

	fmt.Printf("\n⚠️  Found %d potential vulnerabilities:\n", len(vulnerabilities))
	fmt.Printf("   Critical: %d | High: %d | Medium: %d | Low: %d\n\n", critical, high, medium, low)

	// Display vulnerabilities
	for i, vuln := range vulnerabilities {
		fmt.Printf("%d. [%s] %s\n", i+1, vuln.Severity, vuln.Type)
		fmt.Printf("   Description: %s\n", vuln.Description)
		fmt.Printf("   URL: %s\n\n", vuln.URL)
	}

	fmt.Println(strings.Repeat("=", 80))
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  🔍 Web Vulnerability Scanner               ║
║                                                              ║
║              Web Application Security Tool                  ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

⚠️  LEGAL WARNING: Only scan websites you own or have explicit
   permission to test. Unauthorized scanning may be illegal.

`
	fmt.Println(banner)
}

