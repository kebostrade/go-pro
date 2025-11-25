//go:build ignore

package main

import (
	"fmt"
	"strings"
	"sync"
)

// Solution: Concurrent Web Crawler

type SafeCache struct {
	mu      sync.Mutex
	visited map[string]bool
}

func (c *SafeCache) Visit(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.visited[url] {
		return false
	}
	c.visited[url] = true
	return true
}

func CrawlSolution(url string, depth int, fetcher FetcherSolution, cache *SafeCache, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	// Check if already visited
	if !cache.Visit(url) {
		return
	}

	// Fetch URLs
	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}

	// Crawl found URLs concurrently
	for _, u := range urls {
		wg.Add(1)
		go CrawlSolution(u, depth-1, fetcher, cache, wg)
	}
}

// FetcherSolution interface
type FetcherSolution interface {
	Fetch(url string) (urls []string, err error)
}

// fakeFetcherSolution implementation
type fakeFetcherSolution map[string]*fakeResultSolution

type fakeResultSolution struct {
	body string
	urls []string
}

func (f fakeFetcherSolution) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("✓ Crawled: %s\n", url)
		return res.urls, nil
	}
	fmt.Printf("✗ Not found: %s\n", url)
	return nil, fmt.Errorf("not found: %s", url)
}

var fetcherSolution = fakeFetcherSolution{
	"https://golang.org/": &fakeResultSolution{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResultSolution{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResultSolution{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResultSolution{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

func main() {
	fmt.Println("Web Crawler Solution")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	cache := &SafeCache{visited: make(map[string]bool)}
	var wg sync.WaitGroup

	wg.Add(1)
	go CrawlSolution("https://golang.org/", 3, fetcherSolution, cache, &wg)
	wg.Wait()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("Total URLs crawled: %d\n", len(cache.visited))
}
