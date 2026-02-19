# Anongo 👻

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-Donate-orange.svg?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/teodev1611)

**Anongo** (Go Ghost) is a high-performance, security-focused port of [**AnonGT**](https://github.com/gt0day/AnonGT) (originally in Python by **gt0day**) rewritten in **Go**. It creates a transparent anonymity layer by forcing all system network traffic through the Tor network, implementing advanced anti-leak protections and forensic cleaning.


## 🚀 How it Works

Anongo creates a "Ghost Tunnel" using a combination of **Tor's Transparent Proxy** and **Surgical Iptables Redirection**:
1. **Network Shielding:** It creates custom `iptables` chains to redirect all TCP traffic to Tor's `TransPort` (9040) and DNS traffic to `DNSPort` (5353).
2. **Anti-Leak Engine:** It automatically disables IPv6 (a common source of leaks) and blocks non-Tor UDP traffic (preventing protocols like QUIC from bypassing the proxy).
3. **Privilege Dropping:** Tor is executed under a specific system user (`tor` or `debian-tor`). This allows Anongo to tell `iptables`: "Redirect everything EXCEPT the traffic coming from the Tor user," preventing infinite traffic loops.
4. **Watchdog Monitoring:** A background goroutine verifies connection integrity every 15 seconds. If rules are deleted or the IP leaks, it instantly reapplies the shield.

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

Anongo requires **root** privileges to manage network interfaces.

```bash
# Basic run (Spanish by default)
sudo ./anongo

# Run in English
sudo ./anongo -lang en

# Run with logging enabled
sudo ./anongo -logs
```

### Menu Options
- **Activate Anonymity:** Secures the connection and starts the tunnel.
- **Stop Tunnel:** Safely restores original network settings without flushing your personal rules.
- **Change Identity:** Restarts Tor circuits to obtain a new public IP.
- **Detailed Check:** Displays a table comparing your public vs. encrypted interface.
- **Anti-Forensics:** Wipes system traces and clears RAM caches.
- **Emergency Cleanup:** Failsafe option to force-restore all settings.

## 🔒 Security Features
- **Surgical Iptables:** Uses `ANONGO_NAT` and `ANONGO_FILTER` chains. It won't break your existing firewall.
- **Memory Safety:** Written in Go, providing better memory management than the original Python version.
- **Forensic Cleaning:** Clears `drop_caches` and session logs to minimize the footprint left on the machine.

## ☕ Support

If you find **Anongo** useful and want to support its development, you can buy me a coffee!

[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20A%20Coffee-Donate-orange.svg?logo=buy-me-a-coffee&logoColor=white)](https://buymeacoffee.com/teodev1611)

## 📄 License
This project is licensed under the **GPL v3 License**.
