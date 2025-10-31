package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"strings"
	"time"
)

func main() {
	printBanner()

	// Parse flags
	host := flag.String("host", "", "Host to check (e.g., example.com)")
	port := flag.String("port", "443", "Port to connect to")
	flag.Parse()

	if *host == "" {
		fmt.Println("❌ Error: Host is required")
		fmt.Println("Usage: ssl-checker -host example.com")
		return
	}

	fmt.Printf("🎯 Target: %s:%s\n\n", *host, *port)

	// Connect and get certificate
	address := fmt.Sprintf("%s:%s", *host, *port)
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		fmt.Printf("❌ Error connecting: %v\n", err)
		return
	}
	defer conn.Close()

	// Get certificate chain
	certs := conn.ConnectionState().PeerCertificates

	if len(certs) == 0 {
		fmt.Println("❌ No certificates found")
		return
	}

	// Analyze certificate
	analyzeCertificate(certs[0])

	// Check certificate chain
	fmt.Println("\n📋 CERTIFICATE CHAIN")
	fmt.Println(strings.Repeat("-", 80))
	for i, cert := range certs {
		fmt.Printf("%d. %s\n", i+1, cert.Subject.CommonName)
		fmt.Printf("   Issuer: %s\n", cert.Issuer.CommonName)
		fmt.Printf("   Valid: %s to %s\n", cert.NotBefore.Format("2006-01-02"), cert.NotAfter.Format("2006-01-02"))
	}

	// Security recommendations
	displayRecommendations(certs[0], conn.ConnectionState())
}

func analyzeCertificate(cert *x509.Certificate) {
	fmt.Println("📋 CERTIFICATE INFORMATION")
	fmt.Println(strings.Repeat("=", 80))

	// Basic info
	fmt.Printf("Common Name: %s\n", cert.Subject.CommonName)
	fmt.Printf("Organization: %s\n", cert.Subject.Organization)
	fmt.Printf("Issuer: %s\n", cert.Issuer.CommonName)

	// Validity
	fmt.Printf("\nValidity:\n")
	fmt.Printf("  Not Before: %s\n", cert.NotBefore.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("  Not After:  %s\n", cert.NotAfter.Format("2006-01-02 15:04:05 MST"))

	// Check expiration
	daysUntilExpiry := int(time.Until(cert.NotAfter).Hours() / 24)
	if daysUntilExpiry < 0 {
		fmt.Printf("  Status: ❌ EXPIRED (%d days ago)\n", -daysUntilExpiry)
	} else if daysUntilExpiry < 30 {
		fmt.Printf("  Status: ⚠️  EXPIRING SOON (%d days remaining)\n", daysUntilExpiry)
	} else {
		fmt.Printf("  Status: ✅ Valid (%d days remaining)\n", daysUntilExpiry)
	}

	// Subject Alternative Names
	if len(cert.DNSNames) > 0 {
		fmt.Printf("\nSubject Alternative Names:\n")
		for _, name := range cert.DNSNames {
			fmt.Printf("  - %s\n", name)
		}
	}

	// Key information
	fmt.Printf("\nKey Information:\n")
	fmt.Printf("  Algorithm: %s\n", cert.PublicKeyAlgorithm)
	fmt.Printf("  Signature: %s\n", cert.SignatureAlgorithm)

	// Extensions
	fmt.Printf("\nExtensions:\n")
	fmt.Printf("  Key Usage: ")
	if cert.KeyUsage&x509.KeyUsageDigitalSignature != 0 {
		fmt.Print("Digital Signature ")
	}
	if cert.KeyUsage&x509.KeyUsageKeyEncipherment != 0 {
		fmt.Print("Key Encipherment ")
	}
	fmt.Println()

	if len(cert.ExtKeyUsage) > 0 {
		fmt.Printf("  Extended Key Usage: ")
		for _, usage := range cert.ExtKeyUsage {
			switch usage {
			case x509.ExtKeyUsageServerAuth:
				fmt.Print("Server Authentication ")
			case x509.ExtKeyUsageClientAuth:
				fmt.Print("Client Authentication ")
			}
		}
		fmt.Println()
	}
}

func displayRecommendations(cert *x509.Certificate, state tls.ConnectionState) {
	fmt.Println("\n📋 SECURITY ANALYSIS")
	fmt.Println(strings.Repeat("=", 80))

	var issues []string
	var warnings []string
	var good []string

	// Check TLS version
	tlsVersion := getTLSVersion(state.Version)
	if state.Version < tls.VersionTLS12 {
		issues = append(issues, fmt.Sprintf("Using outdated TLS version: %s (use TLS 1.2 or higher)", tlsVersion))
	} else if state.Version == tls.VersionTLS12 {
		warnings = append(warnings, fmt.Sprintf("Using TLS 1.2 (consider upgrading to TLS 1.3)"))
	} else {
		good = append(good, fmt.Sprintf("Using modern TLS version: %s", tlsVersion))
	}

	// Check cipher suite
	cipherSuite := tls.CipherSuiteName(state.CipherSuite)
	good = append(good, fmt.Sprintf("Cipher Suite: %s", cipherSuite))

	// Check certificate expiration
	daysUntilExpiry := int(time.Until(cert.NotAfter).Hours() / 24)
	if daysUntilExpiry < 0 {
		issues = append(issues, "Certificate has expired")
	} else if daysUntilExpiry < 30 {
		warnings = append(warnings, fmt.Sprintf("Certificate expires in %d days", daysUntilExpiry))
	} else {
		good = append(good, fmt.Sprintf("Certificate valid for %d more days", daysUntilExpiry))
	}

	// Check signature algorithm
	if strings.Contains(cert.SignatureAlgorithm.String(), "SHA1") {
		issues = append(issues, "Using weak SHA-1 signature algorithm")
	} else {
		good = append(good, fmt.Sprintf("Using strong signature algorithm: %s", cert.SignatureAlgorithm))
	}

	// Display results
	if len(issues) > 0 {
		fmt.Println("\n❌ CRITICAL ISSUES:")
		for _, issue := range issues {
			fmt.Printf("  • %s\n", issue)
		}
	}

	if len(warnings) > 0 {
		fmt.Println("\n⚠️  WARNINGS:")
		for _, warning := range warnings {
			fmt.Printf("  • %s\n", warning)
		}
	}

	if len(good) > 0 {
		fmt.Println("\n✅ GOOD PRACTICES:")
		for _, item := range good {
			fmt.Printf("  • %s\n", item)
		}
	}

	// Overall score
	score := 100 - (len(issues) * 30) - (len(warnings) * 10)
	if score < 0 {
		score = 0
	}

	fmt.Printf("\n📊 SECURITY SCORE: %d/100\n", score)

	if score >= 80 {
		fmt.Println("   Rating: ✅ Excellent")
	} else if score >= 60 {
		fmt.Println("   Rating: ⚠️  Good (improvements recommended)")
	} else {
		fmt.Println("   Rating: ❌ Poor (immediate action required)")
	}

	fmt.Println(strings.Repeat("=", 80))
}

func getTLSVersion(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║                  🔒 SSL/TLS Certificate Checker             ║
║                                                              ║
║              SSL/TLS Security Analysis Tool                 ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

