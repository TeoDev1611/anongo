# Anongo 👻

**Anongo** (Go Ghost) is a high-performance, security-focused port of
[**AnonGT**](https://github.com/gt0day/AnonGT) (originally in Python) rewritten
in **Go**. It creates a transparent anonymity layer by forcing all system
network traffic through the Tor network, implementing advanced anti-leak
protections and forensic cleaning.

## 🌟 Why the Go Port? (Advantages over Python)

- **Parallel Execution:** Leverages Go's Goroutines to run the network watchdog, ARP monitor, and Tor manager concurrently with minimal overhead.
- **Enhanced Cryptography:** Upgraded from AES-CBC to **AES-256-GCM**, providing authenticated encryption and better resistance against data tampering.
- **Static Binary:** No need for a Python interpreter or complex `pip` dependencies; just a single, fast-executing binary.
- **Lower Resource Footprint:** Significantly more efficient memory and CPU usage compared to the original Python implementation.
- **Robust Network Handling:** Uses `gopacket` for low-level network monitoring, which is faster and more reliable than Scapy.

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-Donate-orange.svg?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/teodev1611)

## 🚀 How it Works

Anongo creates a "Ghost Tunnel" using a combination of **Tor's Transparent
Proxy** and **Surgical Iptables Redirection**:

1. **Network Shielding:** It creates custom `iptables` chains (`ANONGO_NAT` and `ANONGO_FILTER`) to redirect all
   TCP traffic to Tor's `TransPort` (9040) and DNS traffic to `DNSPort` (5353).
2. **Anti-Leak Engine:** 
   - **IPv6 Killswitch:** Automatically drops all IPv6 traffic to prevent common leaks.
   - **UDP Filtering:** Blocks non-Tor UDP traffic (e.g., QUIC, STUN) while allowing DNS redirection.
   - **LAN Exemption:** Automatically detects and excludes local networks (127.0.0.1, 192.168.x.x, etc.) so you don't lose access to your local router or devices.
3. **Privilege Dropping:** Tor is executed under a specific system user (`tor` or `debian-tor`). This allows Anongo to tell `iptables`: "Redirect everything EXCEPT the traffic coming from the Tor user," preventing infinite traffic loops.
4. **Watchdog Monitoring:** A background goroutine verifies connection integrity
   every 15 seconds. If rules are deleted or the IP leaks, it instantly
   reapplies the shield.

## 🏗️ Technical Architecture

The project is organized into modular packages to ensure maintainability and high performance:

- **`pkg/network`**: Manages `iptables` rules, IPv6 disabling, and the network Watchdog.
- **`pkg/tor`**: Handles the lifecycle of the Tor process, identity switching, and circuit health.
- **`pkg/security`**: Implements anti-forensics measures (RAM wiping, history truncation, log cleaning).
- **`pkg/crypto`**: Provides high-grade encryption for local data (AES-256-GCM).
- **`pkg/i18n`**: Multi-language support (English/Spanish).

## 🛠️ Requirements & Dependencies

### System Binaries

You must have the following installed on your Linux system:

- **tor**: The core anonymity service.
- **iptables / ip6tables**: For network redirection.
- **procps (pkill)**: To manage process cleanup.
- **kmod**: To ensure iptables modules are loaded.

### Go Dependencies

The project leverages these libraries:

- `github.com/pterm/pterm`: For the interactive CLI and dashboard.
- `github.com/coreos/go-iptables`: For surgical firewall management.
- `golang.org/x/net/proxy`: For secure SOCKS5 verification.

## 📥 Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/TeoDev1611/anongo.git
   cd anongo
   ```

2. **Build the binary:**
   ```bash
   go build -o anongo ./cmd/anongo/main.go
   ```

## 🎮 Usage

Anongo requires **root** privileges to manage network interfaces and process state.

### Command Line Flags

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-lang` | Set the interface language (`en` or `es`). | `es` |
| `-logs` | Enable logging the session output to `anongo_session.log`. | `false` |

### Basic Execution

```bash
# Run in English with logs enabled
sudo ./anongo -lang en -logs
```

### Menu Options

- **Activate Anonymity:** Secures the connection and starts the tunnel.
- **Stop Tunnel:** Safely restores original network settings without flushing
  your personal rules.
- **Change Identity:** Restarts Tor circuits to obtain a new public IP.
- **Detailed Check:** Displays a table comparing your public vs. encrypted
  interface.
- **Anti-Forensics:** Wipes system traces, clears RAM caches (`drop_caches`), and truncates bash/zsh history.
- **Emergency Cleanup:** Failsafe option to force-restore all settings.

## 🔒 Security Features

- **Surgical Iptables:** Uses isolated chains. It won't interfere with your custom firewall rules.
- **Memory Safety:** Written in Go, eliminating buffer overflow risks present in C-based alternatives.
- **Anti-Forensics:** Clears `/proc/sys/vm/drop_caches`, system logs (`auth.log`, `syslog`), and session histories to minimize the forensic footprint.

## ⚠️ Security Disclaimer

**Anongo is a tool for security research and privacy.** While it provides strong anonymity, no tool is 100% foolproof. 
- Using this tool does not make you immune to fingerprinting or application-level leaks (e.g., browser plugins).
- Always use a privacy-hardened browser (like Tor Browser) even when the tunnel is active.
- **The developers are not responsible for any misuse or damages caused by this tool.**

## ☕ Support

If you find **Anongo** useful and want to support its development, you can buy
me a coffee!

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-Donate-orange.svg?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/teodev1611)

## 📄 License

This project is licensed under the **GPL v3 License**.
