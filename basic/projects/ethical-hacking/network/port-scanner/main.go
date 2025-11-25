package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ⚠️ LEGAL WARNING ⚠️
// This tool is for EDUCATIONAL PURPOSES ONLY.
// Only scan systems you own or have explicit permission to test.
// Unauthorized port scanning may be illegal in your jurisdiction.

type ScanResult struct {
	Port    int
	State   string
	Service string
}

var commonPorts = map[int]string{
	20:    "FTP Data",
	21:    "FTP Control",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	80:    "HTTP",
	110:   "POP3",
	143:   "IMAP",
	443:   "HTTPS",
	445:   "SMB",
	3306:  "MySQL",
	3389:  "RDP",
	5432:  "PostgreSQL",
	5900:  "VNC",
	6379:  "Redis",
	8080:  "HTTP Proxy",
	8443:  "HTTPS Alt",
	27017: "MongoDB",
}

func main() {
	printBanner()
	printWarning()

	// Parse flags
	host := flag.String("host", "localhost", "Target host to scan")
	portRange := flag.String("ports", "1-1000", "Port range (e.g., 1-1000 or 80,443,8080)")
	timeout := flag.Int("timeout", 1000, "Connection timeout in milliseconds")
	workers := flag.Int("workers", 100, "Number of concurrent workers")
	flag.Parse()

	// Parse port range
	ports, err := parsePortRange(*portRange)
	if err != nil {
		fmt.Printf("❌ Error parsing port range: %v\n", err)
		return
	}

	fmt.Printf("🎯 Target: %s\n", *host)
	fmt.Printf("📊 Scanning %d ports with %d workers...\n\n", len(ports), *workers)

	// Scan ports
	results := scanPorts(*host, ports, time.Duration(*timeout)*time.Millisecond, *workers)

	// Display results
	displayResults(results)
}

func parsePortRange(portRange string) ([]int, error) {
	var ports []int

	// Check if it's a range (e.g., "1-1000")
	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid port range format")
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		for i := start; i <= end; i++ {
			ports = append(ports, i)
		}
	} else {
		// Comma-separated ports (e.g., "80,443,8080")
		parts := strings.Split(portRange, ",")
		for _, p := range parts {
			port, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				return nil, err
			}
			ports = append(ports, port)
		}
	}

	return ports, nil
}

func scanPorts(host string, ports []int, timeout time.Duration, workers int) []ScanResult {
	var results []ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create job channel
	jobs := make(chan int, len(ports))
	for _, port := range ports {
		jobs <- port
	}
	close(jobs)

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				result := scanPort(host, port, timeout)
				if result.State == "open" {
					mu.Lock()
					results = append(results, result)
					mu.Unlock()
					fmt.Printf("✅ Port %d (%s) - OPEN\n", result.Port, result.Service)
				}
			}
		}()
	}

	wg.Wait()

	// Sort results by port number
	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

func scanPort(host string, port int, timeout time.Duration) ScanResult {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)

	result := ScanResult{
		Port:    port,
		State:   "closed",
		Service: getServiceName(port),
	}

	if err == nil {
		result.State = "open"
		conn.Close()
	}

	return result
}

func getServiceName(port int) string {
	if service, ok := commonPorts[port]; ok {
		return service
	}
	return "Unknown"
}

func displayResults(results []ScanResult) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📋 SCAN RESULTS")
	fmt.Println(strings.Repeat("=", 60))

	if len(results) == 0 {
		fmt.Println("❌ No open ports found")
		return
	}

	fmt.Printf("\n✅ Found %d open port(s):\n\n", len(results))
	fmt.Printf("%-10s %-10s %-20s\n", "PORT", "STATE", "SERVICE")
	fmt.Println(strings.Repeat("-", 60))

	for _, result := range results {
		fmt.Printf("%-10d %-10s %-20s\n", result.Port, result.State, result.Service)
	}

	fmt.Println(strings.Repeat("=", 60))
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  🔍 Port Scanner                            ║
║                                                              ║
║              Network Port Scanning Tool                     ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

func printWarning() {
	warning := `
⚠️  LEGAL WARNING ⚠️

This tool is for EDUCATIONAL and AUTHORIZED TESTING purposes ONLY.

You MUST have explicit permission to scan any system you do not own.
Unauthorized port scanning may be illegal and could result in:
  • Criminal prosecution
  • Civil liability
  • Network access termination

By using this tool, you agree to use it responsibly and legally.

Press Ctrl+C to cancel or wait 3 seconds to continue...
`
	fmt.Println(warning)
	time.Sleep(3 * time.Second)
	fmt.Println()
}

