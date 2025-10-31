package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// ⚠️ LEGAL WARNING ⚠️
// This tool is for EDUCATIONAL PURPOSES ONLY.
// Only scan networks you own or have explicit permission to test.

type Host struct {
	IP       string
	Hostname string
	Status   string
	Latency  time.Duration
}

func main() {
	printBanner()

	// Parse flags
	subnet := flag.String("subnet", "192.168.1.0/24", "Subnet to scan (CIDR notation)")
	timeout := flag.Int("timeout", 1000, "Ping timeout in milliseconds")
	workers := flag.Int("workers", 50, "Number of concurrent workers")
	flag.Parse()

	fmt.Printf("🌐 Scanning subnet: %s\n", *subnet)
	fmt.Printf("⏱️  Timeout: %dms\n", *timeout)
	fmt.Printf("👷 Workers: %d\n\n", *workers)

	// Parse subnet
	ips, err := parseSubnet(*subnet)
	if err != nil {
		fmt.Printf("❌ Error parsing subnet: %v\n", err)
		return
	}

	fmt.Printf("📊 Scanning %d IP addresses...\n\n", len(ips))

	// Scan network
	hosts := scanNetwork(ips, time.Duration(*timeout)*time.Millisecond, *workers)

	// Display results
	displayResults(hosts)
}

func parseSubnet(subnet string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// Remove network and broadcast addresses
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}

	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func scanNetwork(ips []string, timeout time.Duration, workers int) []Host {
	var hosts []Host
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create job channel
	jobs := make(chan string, len(ips))
	for _, ip := range ips {
		jobs <- ip
	}
	close(jobs)

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range jobs {
				host := pingHost(ip, timeout)
				if host.Status == "up" {
					mu.Lock()
					hosts = append(hosts, host)
					mu.Unlock()
					fmt.Printf("✅ %s - %s (%.2fms)\n", host.IP, host.Hostname, float64(host.Latency.Microseconds())/1000)
				}
			}
		}()
	}

	wg.Wait()
	return hosts
}

func pingHost(ip string, timeout time.Duration) Host {
	host := Host{
		IP:     ip,
		Status: "down",
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", ip+":80", timeout)
	if err == nil {
		host.Status = "up"
		host.Latency = time.Since(start)
		conn.Close()

		// Try to resolve hostname
		names, err := net.LookupAddr(ip)
		if err == nil && len(names) > 0 {
			host.Hostname = names[0]
		} else {
			host.Hostname = "Unknown"
		}
	}

	return host
}

func displayResults(hosts []Host) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 NETWORK MAP")
	fmt.Println(strings.Repeat("=", 80))

	if len(hosts) == 0 {
		fmt.Println("❌ No active hosts found")
		return
	}

	fmt.Printf("\n✅ Found %d active host(s):\n\n", len(hosts))
	fmt.Printf("%-20s %-10s %-30s %-15s\n", "IP ADDRESS", "STATUS", "HOSTNAME", "LATENCY")
	fmt.Println(strings.Repeat("-", 80))

	for _, host := range hosts {
		latency := fmt.Sprintf("%.2fms", float64(host.Latency.Microseconds())/1000)
		fmt.Printf("%-20s %-10s %-30s %-15s\n", host.IP, host.Status, host.Hostname, latency)
	}

	fmt.Println(strings.Repeat("=", 80))

	// Network statistics
	fmt.Println("\n📊 NETWORK STATISTICS")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Total hosts scanned: %d\n", len(hosts))
	fmt.Printf("Active hosts: %d\n", len(hosts))

	// Calculate average latency
	var totalLatency time.Duration
	for _, host := range hosts {
		totalLatency += host.Latency
	}
	avgLatency := totalLatency / time.Duration(len(hosts))
	fmt.Printf("Average latency: %.2fms\n", float64(avgLatency.Microseconds())/1000)

	fmt.Println(strings.Repeat("=", 80))
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  🗺️  Network Mapper                         ║
║                                                              ║
║              Network Discovery Tool                         ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

⚠️  LEGAL WARNING: Only scan networks you own or have permission to test.

`
	fmt.Println(banner)
}

