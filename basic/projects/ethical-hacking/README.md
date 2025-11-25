# 🔒 Ethical Hacking with Go

A comprehensive collection of security tools and penetration testing utilities built with Go for educational purposes.

## ⚠️ LEGAL DISCLAIMER

**READ THIS CAREFULLY BEFORE USING ANY TOOLS**

This project is for **EDUCATIONAL PURPOSES ONLY**. These tools are designed to help you:
- Learn about network security
- Understand common vulnerabilities
- Practice ethical hacking skills
- Secure your own systems

### Legal Requirements

✅ **YOU MAY:**
- Use these tools on systems you own
- Use these tools on systems where you have explicit written permission
- Use these tools in authorized penetration testing engagements
- Use these tools in controlled lab environments

❌ **YOU MAY NOT:**
- Scan, test, or attack systems without authorization
- Use these tools for malicious purposes
- Violate computer fraud and abuse laws
- Access systems or data without permission

**Unauthorized use of these tools may be illegal and could result in criminal prosecution, civil liability, and other legal consequences.**

By using these tools, you agree to use them responsibly, ethically, and legally.

## 📋 Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Tools](#tools)
- [Usage Examples](#usage-examples)
- [Learning Resources](#learning-resources)

## ✨ Features

### Network Security Tools
- ✅ **Port Scanner**: Fast concurrent port scanning
- ✅ **Network Mapper**: Network discovery and host enumeration
- ✅ **Packet Sniffer**: Network traffic analysis (requires root)

### Web Security Tools
- ✅ **Web Vulnerability Scanner**: XSS, SQL injection, security headers
- ✅ **SSL/TLS Checker**: Certificate analysis and security assessment

### Cryptography Tools
- ✅ **Hash Cracker**: MD5, SHA1, SHA256, bcrypt password cracking
- ✅ **Password Analyzer**: Password strength analysis and recommendations

## 📦 Prerequisites

- **Go** 1.21 or higher
- **libpcap** (for packet sniffer)
  - Ubuntu/Debian: `sudo apt-get install libpcap-dev`
  - macOS: `brew install libpcap`
  - Windows: Install WinPcap or Npcap
- **Root/Administrator privileges** (for packet sniffer only)

## 🚀 Installation

### 1. Clone and Navigate

```bash
cd basic/projects/ethical-hacking
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Build All Tools

```bash
make build
```

This creates binaries in the `bin/` directory.

## 🛠️ Tools

### 1. Port Scanner

Scan network ports to discover open services.

**Features:**
- Concurrent scanning with configurable workers
- Custom port ranges
- Service identification
- Fast and efficient

**Usage:**

```bash
# Scan localhost ports 1-1000
./bin/port-scanner -host localhost -ports 1-1000

# Scan specific ports
./bin/port-scanner -host example.com -ports 80,443,8080

# Custom timeout and workers
./bin/port-scanner -host 192.168.1.1 -ports 1-65535 -timeout 500 -workers 200
```

**Example Output:**

```
✅ Port 22 (SSH) - OPEN
✅ Port 80 (HTTP) - OPEN
✅ Port 443 (HTTPS) - OPEN

📋 SCAN RESULTS
============================================================
✅ Found 3 open port(s):

PORT       STATE      SERVICE
------------------------------------------------------------
22         open       SSH
80         open       HTTP
443        open       HTTPS
```

### 2. Network Mapper

Discover active hosts on a network.

**Features:**
- Subnet scanning (CIDR notation)
- Hostname resolution
- Latency measurement
- Network statistics

**Usage:**

```bash
# Scan local network
./bin/network-mapper -subnet 192.168.1.0/24

# Custom timeout and workers
./bin/network-mapper -subnet 10.0.0.0/24 -timeout 2000 -workers 100
```

**Example Output:**

```
✅ 192.168.1.1 - router.local (2.34ms)
✅ 192.168.1.100 - desktop.local (1.56ms)
✅ 192.168.1.101 - laptop.local (3.21ms)

📋 NETWORK MAP
================================================================================
IP ADDRESS           STATUS     HOSTNAME                       LATENCY
--------------------------------------------------------------------------------
192.168.1.1          up         router.local                   2.34ms
192.168.1.100        up         desktop.local                  1.56ms
192.168.1.101        up         laptop.local                   3.21ms
```

### 3. Packet Sniffer

Capture and analyze network traffic.

**Features:**
- Real-time packet capture
- Protocol analysis (TCP, UDP, ICMP, HTTP, HTTPS, DNS)
- BPF filtering
- Statistics tracking

**Usage:**

```bash
# Capture all traffic (requires root)
sudo ./bin/packet-sniffer -interface eth0

# Capture HTTP traffic only
sudo ./bin/packet-sniffer -interface eth0 -filter "tcp port 80"

# Capture DNS queries
sudo ./bin/packet-sniffer -interface eth0 -filter "udp port 53"

# Capture limited packets
sudo ./bin/packet-sniffer -interface eth0 -count 100
```

**Example Output:**

```
[14:23:45.123] HTTP: 192.168.1.100:54321 → 93.184.216.34:80 (512 bytes)
[14:23:45.234] HTTPS: 192.168.1.100:54322 → 172.217.14.206:443 (1024 bytes)
[14:23:45.345] DNS: 192.168.1.100:53241 → 8.8.8.8:53 (64 bytes)

📊 CAPTURE STATISTICS
============================================================
Total packets: 150
  TCP:   85 (56.7%)
  UDP:   45 (30.0%)
  ICMP:  10 (6.7%)
  HTTP:  25 (16.7%)
  HTTPS: 60 (40.0%)
  DNS:   20 (13.3%)
```

### 4. Web Vulnerability Scanner

Scan websites for common security vulnerabilities.

**Features:**
- HTTPS redirect check
- Security headers analysis
- XSS vulnerability detection
- SQL injection testing
- Directory listing check

**Usage:**

```bash
# Scan a website
./bin/web-scanner -url https://example.com

# Custom timeout
./bin/web-scanner -url https://example.com -timeout 30
```

**Example Output:**

```
📋 SCAN RESULTS
================================================================================
⚠️  Found 3 potential vulnerabilities:
   Critical: 0 | High: 1 | Medium: 1 | Low: 1

1. [High] Potential XSS
   Description: User input reflected in response without sanitization
   URL: https://example.com?test=<script>alert('XSS')</script>

2. [Medium] Missing Security Header
   Description: Missing X-Frame-Options (Clickjacking protection)
   URL: https://example.com

3. [Low] Missing Security Header
   Description: Missing Content-Security-Policy (XSS protection)
   URL: https://example.com
```

### 5. SSL/TLS Checker

Analyze SSL/TLS certificates and security configuration.

**Features:**
- Certificate information
- Expiration checking
- TLS version analysis
- Cipher suite evaluation
- Security scoring

**Usage:**

```bash
# Check SSL certificate
./bin/ssl-checker -host example.com

# Custom port
./bin/ssl-checker -host example.com -port 8443
```

**Example Output:**

```
📋 CERTIFICATE INFORMATION
================================================================================
Common Name: example.com
Organization: Example Inc
Issuer: Let's Encrypt Authority X3

Validity:
  Not Before: 2024-01-01 00:00:00 UTC
  Not After:  2024-04-01 00:00:00 UTC
  Status: ✅ Valid (45 days remaining)

📋 SECURITY ANALYSIS
================================================================================
✅ GOOD PRACTICES:
  • Using modern TLS version: TLS 1.3
  • Cipher Suite: TLS_AES_128_GCM_SHA256
  • Certificate valid for 45 more days
  • Using strong signature algorithm: SHA256-RSA

📊 SECURITY SCORE: 95/100
   Rating: ✅ Excellent
```

### 6. Hash Cracker

Crack password hashes using dictionary attacks.

**Features:**
- Multiple hash types (MD5, SHA1, SHA256, bcrypt)
- Custom wordlists
- Performance metrics
- Common password database

**Usage:**

```bash
# Crack MD5 hash
./bin/hash-cracker -hash 5f4dcc3b5aa765d61d8327deb882cf99 -type md5

# Use custom wordlist
./bin/hash-cracker -hash <hash> -type sha256 -wordlist wordlist.txt

# Crack bcrypt hash
./bin/hash-cracker -hash '$2a$10$...' -type bcrypt -wordlist passwords.txt
```

**Example Output:**

```
🔨 Progress: 5000/10000 (50.0%)

✅ HASH CRACKED!
   Password: password
   Time: 234ms
   Speed: 42735 hashes/sec
```

### 7. Password Analyzer

Analyze password strength and provide recommendations.

**Features:**
- Strength scoring
- Entropy calculation
- Time-to-crack estimation
- Pattern detection
- Security recommendations

**Usage:**

```bash
# Analyze single password
./bin/password-analyzer -password "MyP@ssw0rd123"

# Analyze passwords from file
./bin/password-analyzer -file passwords.txt
```

**Example Output:**

```
📋 PASSWORD ANALYSIS
================================================================================
Password: ************* (length: 13)

🔒 Strength: ✅ Strong
📊 Score: 75/100
🔢 Entropy: 65.23 bits
⏱️  Time to crack: 2.3 years

✅ GOOD PRACTICES:
  • Uses uppercase and lowercase letters
  • Contains numbers
  • Contains special characters
  • Good length (13 characters)

💡 SUGGESTIONS:
  • Consider using 16+ characters for maximum security
```

## 📚 Learning Resources

### Recommended Reading

- **Books:**
  - "The Web Application Hacker's Handbook" by Dafydd Stuttard
  - "Metasploit: The Penetration Tester's Guide" by David Kennedy
  - "Black Hat Go" by Tom Steele, Chris Patten, Dan Kottmann

- **Online Resources:**
  - OWASP Top 10: https://owasp.org/www-project-top-ten/
  - HackTheBox: https://www.hackthebox.com/
  - TryHackMe: https://tryhackme.com/

### Practice Environments

- **Legal Practice Platforms:**
  - HackTheBox
  - TryHackMe
  - PentesterLab
  - OverTheWire
  - VulnHub

## 🎓 Learning Outcomes

After completing this tutorial, you'll understand:

### Network Security
- ✅ Port scanning techniques
- ✅ Network reconnaissance
- ✅ Packet analysis
- ✅ Protocol understanding

### Web Security
- ✅ Common web vulnerabilities (XSS, SQLi)
- ✅ Security headers
- ✅ SSL/TLS configuration
- ✅ Certificate validation

### Cryptography
- ✅ Hash functions
- ✅ Password security
- ✅ Dictionary attacks
- ✅ Password strength metrics

### General
- ✅ Ethical hacking principles
- ✅ Legal considerations
- ✅ Responsible disclosure
- ✅ Security best practices

## 🔧 Development

### Run Tests

```bash
make test
```

### Code Quality

```bash
make quality  # Run all quality checks
make lint     # Run security linter
make vet      # Run go vet
make fmt      # Format code
```

## 📝 License

MIT License - for educational purposes only.

## 🤝 Contributing

Contributions are welcome! Please ensure all tools include:
- Legal warnings
- Educational documentation
- Responsible use guidelines
- Proper error handling

## ⚖️ Responsible Disclosure

If you discover vulnerabilities using these tools:
1. Do NOT exploit them
2. Report to the system owner immediately
3. Follow responsible disclosure practices
4. Allow time for fixes before public disclosure

---

**Remember: With great power comes great responsibility. Use these tools ethically and legally.**

