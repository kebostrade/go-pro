// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package main provides functionality for the GO-PRO Learning Platform.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("🧪 Testing GO-PRO API Endpoints...")

	// Wait a moment for server to be ready.
	time.Sleep(1 * time.Second)

	// Test endpoints.
	testEndpoint("Health Check", "http://localhost:8080/api/v1/health")
	testEndpoint("All Courses", "http://localhost:8080/api/v1/courses")
	testEndpoint("GO-PRO Course", "http://localhost:8080/api/v1/courses/go-pro")
	testEndpoint("Course Lessons", "http://localhost:8080/api/v1/courses/go-pro/lessons")
	testEndpoint("Demo User Progress", "http://localhost:8080/api/v1/progress/demo-user")

	fmt.Println("\n✅ API testing completed!")
}

func testEndpoint(name, url string) {
	fmt.Printf("\n📡 Testing: %s\n", name)
	fmt.Printf("URL: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", resp.Status)

	// Pretty print JSON if possible.
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err == nil {
		prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
		fmt.Printf("Response: %s\n", string(prettyJSON))
	} else {
		fmt.Printf("Response: %s\n", string(body))
	}

	if resp.StatusCode == 200 {
		fmt.Println("✅ Success")
	} else {
		fmt.Printf("⚠️  Status Code: %d\n", resp.StatusCode)
	}
}
