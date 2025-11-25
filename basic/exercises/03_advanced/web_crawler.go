package main

import (
	"fmt"
	"strings"
)

// Exercise: Concurrent Web Crawler
// Build a concurrent web crawler that:
// 1. Crawls URLs concurrently using goroutines
// 2. Avoids visiting the same URL twice
// 3. Limits the depth of crawling
// 4. Uses channels for communication

// Fetcher interface for fetching URLs
type Fetcher interface {
	Fetch(url string) (urls []string, err error)
}

// TODO: Implement the Crawl function
// Crawl uses fetcher to recursively crawl pages starting with url, to a maximum of depth
func Crawl(url string, depth int, fetcher Fetcher) {
	// Your code here
	// Hints:
	// - Use a map to track visited URLs
	// - Use sync.Mutex to protect the map
	// - Use sync.WaitGroup to wait for goroutines
	// - Use channels to communicate between goroutines
}

// fakeFetcher is a mock fetcher for testing
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("Found: %s\n", url)
		return res.urls, nil
	}
	fmt.Printf("Not found: %s\n", url)
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

func main() {
	fmt.Println("Web Crawler Challenge")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	Crawl("https://golang.org/", 3, fetcher)

	fmt.Println("\nNote: Implement the Crawl function to see results!")
}
