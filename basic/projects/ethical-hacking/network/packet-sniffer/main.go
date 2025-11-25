package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// ⚠️ LEGAL WARNING ⚠️
// This tool is for EDUCATIONAL PURPOSES ONLY.
// Only capture packets on networks you own or have explicit permission to monitor.
// Unauthorized packet sniffing may be illegal.

type PacketStats struct {
	Total int
	TCP   int
	UDP   int
	ICMP  int
	HTTP  int
	HTTPS int
	DNS   int
	Other int
}

func main() {
	printBanner()

	// Parse flags
	iface := flag.String("interface", "eth0", "Network interface to capture from")
	filter := flag.String("filter", "", "BPF filter (e.g., 'tcp port 80')")
	count := flag.Int("count", 0, "Number of packets to capture (0 = unlimited)")
	output := flag.String("output", "", "Output file for packet capture (.pcap)")
	flag.Parse()

	// Check if running as root
	if os.Geteuid() != 0 {
		fmt.Println("❌ This tool requires root/administrator privileges")
		fmt.Println("   Run with: sudo go run main.go")
		return
	}

	fmt.Printf("🔍 Interface: %s\n", *iface)
	if *filter != "" {
		fmt.Printf("🔧 Filter: %s\n", *filter)
	}
	if *output != "" {
		fmt.Printf("💾 Output: %s\n", *output)
	}
	fmt.Println()

	// Open device
	handle, err := pcap.OpenLive(*iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("❌ Error opening device: %v", err)
	}
	defer handle.Close()

	// Set filter
	if *filter != "" {
		if err := handle.SetBPFFilter(*filter); err != nil {
			log.Fatalf("❌ Error setting BPF filter: %v", err)
		}
	}

	// Create packet source
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Statistics
	stats := &PacketStats{}

	// Handle graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("📡 Capturing packets... (Press Ctrl+C to stop)\n")

	// Capture packets
	packetCount := 0
	for {
		select {
		case packet := <-packetSource.Packets():
			packetCount++
			stats.Total++

			analyzePacket(packet, stats)

			if *count > 0 && packetCount >= *count {
				fmt.Println("\n✅ Capture limit reached")
				displayStats(stats)
				return
			}

		case <-sigterm:
			fmt.Println("\n🛑 Stopping capture...")
			displayStats(stats)
			return
		}
	}
}

func analyzePacket(packet gopacket.Packet, stats *PacketStats) {
	// Get timestamp
	timestamp := packet.Metadata().Timestamp.Format("15:04:05.000")

	// Analyze layers
	var srcIP, dstIP, protocol string
	var srcPort, dstPort string

	// Network layer
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		srcIP = ip.SrcIP.String()
		dstIP = ip.DstIP.String()
	}

	// Transport layer
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort = tcp.SrcPort.String()
		dstPort = tcp.DstPort.String()
		protocol = "TCP"
		stats.TCP++

		// Check for HTTP/HTTPS
		if tcp.DstPort == 80 || tcp.SrcPort == 80 {
			stats.HTTP++
			protocol = "HTTP"
		} else if tcp.DstPort == 443 || tcp.SrcPort == 443 {
			stats.HTTPS++
			protocol = "HTTPS"
		}
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		srcPort = udp.SrcPort.String()
		dstPort = udp.DstPort.String()
		protocol = "UDP"
		stats.UDP++

		// Check for DNS
		if udp.DstPort == 53 || udp.SrcPort == 53 {
			stats.DNS++
			protocol = "DNS"
		}
	} else if packet.Layer(layers.LayerTypeICMPv4) != nil {
		protocol = "ICMP"
		stats.ICMP++
	} else {
		protocol = "Other"
		stats.Other++
	}

	// Print packet info
	fmt.Printf("[%s] %s: %s:%s → %s:%s (%d bytes)\n",
		timestamp, protocol, srcIP, srcPort, dstIP, dstPort, len(packet.Data()))
}

func displayStats(stats *PacketStats) {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("📊 CAPTURE STATISTICS")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Printf("Total packets: %d\n", stats.Total)
	fmt.Printf("  TCP:   %d (%.1f%%)\n", stats.TCP, percentage(stats.TCP, stats.Total))
	fmt.Printf("  UDP:   %d (%.1f%%)\n", stats.UDP, percentage(stats.UDP, stats.Total))
	fmt.Printf("  ICMP:  %d (%.1f%%)\n", stats.ICMP, percentage(stats.ICMP, stats.Total))
	fmt.Printf("  HTTP:  %d (%.1f%%)\n", stats.HTTP, percentage(stats.HTTP, stats.Total))
	fmt.Printf("  HTTPS: %d (%.1f%%)\n", stats.HTTPS, percentage(stats.HTTPS, stats.Total))
	fmt.Printf("  DNS:   %d (%.1f%%)\n", stats.DNS, percentage(stats.DNS, stats.Total))
	fmt.Printf("  Other: %d (%.1f%%)\n", stats.Other, percentage(stats.Other, stats.Total))
	fmt.Println(strings.Repeat("═", 60))
}

func percentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  📡 Packet Sniffer                          ║
║                                                              ║
║              Network Packet Capture Tool                    ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

⚠️  LEGAL WARNING: Only capture packets on networks you own or have
   explicit permission to monitor. Unauthorized packet sniffing may
   be illegal in your jurisdiction.

⚠️  REQUIRES: Root/Administrator privileges

`
	fmt.Println(banner)
	time.Sleep(2 * time.Second)
}

