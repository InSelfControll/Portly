# Portly - Unified Firewall Orchestrator

**Portly** is a unified management tool for configuring NAT (Network Address Translation) rules, firewall services, and security policies across multiple operating systems. It provides both an interactive Terminal User Interface (TUI) and a traditional Command Line Interface (CLI).

## What is Portly Used For?

Portly simplifies network configuration for developers and system administrators by:

- **Port Forwarding (NAT)**: Expose internal services to external networks
- **Firewall Rules**: Open ports for incoming connections with IP-based access control
- **Container Networking**: Easily configure NAT for Docker, Podman, and other container platforms
- **VPN Integration**: Set up port forwarding for Tailscale, Headscale, Twingate
- **Gaming Servers**: Configure ports for Steam, Minecraft, and other game servers
- **Database Access**: Expose PostgreSQL, Redis, and other databases securely with IP restrictions
- **Firewall Management**: Start, stop, and install firewall services
- **Security Policies**: Manage SELinux (RHEL) and AppArmor (Ubuntu)
- **Cross-Platform**: Works identically on RHEL, Ubuntu, Debian, Fedora, and macOS

### Common Use Cases

1. **Expose a Podman container** running on `10.88.0.1:8080` to port `80` on the host
2. **Open port 5432** for PostgreSQL, but only allow access from `192.168.1.100`
3. **Set up Tailscale** exit node with proper port forwarding (port 41641)
4. **Host a Minecraft server** with automatic port configuration (port 25565)
5. **Manage firewall service** - start, stop, or install if missing
6. **Configure SELinux/AppArmor** for proper security policies

## Features

- **Interactive TUI**: Beautiful terminal interface with Bubble Tea
- **Firewall Rules**: Open ports with or without IP-based restrictions
- **NAT Port Forwarding**: Forward external ports to internal destinations
- **Firewall Management**: Start, stop, and auto-install firewall services
- **Security Management**: Control SELinux (RHEL) and AppArmor (Ubuntu)
- **Smart Auto-Fill**: Selecting a product auto-populates suggested ports
- **Product Database**: Pre-configured defaults for 10+ popular services
- **Input Validation**: Port fields accept only numbers, proper IP validation
- **Custom Products**: Use defaults or define your own service names
- **State Management**: Persistent tracking of rules with rollback support
- **Cross-Platform**: Native support for firewalld, nftables, and pfctl

## Supported Platforms

| Platform | Firewall Tool | Security | Package Manager |
|----------|--------------|----------|-----------------|
| RHEL / Rocky / Alma / Fedora | firewalld | SELinux | dnf/yum |
| Ubuntu / Debian | nftables | AppArmor | apt |
| macOS | pfctl | SIP | brew |

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/portly.git
cd portly

# Build the binary
go build -o portly cmd/orchestrator/*.go

# Install system-wide (optional)
sudo cp portly /usr/local/bin/
sudo chmod +x /usr/local/bin/portly
```

### Requirements

- Go 1.21 or later
- Root privileges (for firewall configuration)
- Supported firewall tool installed (firewalld, nftables, or pf)

## Usage

Portly provides two interfaces:
1. **TUI** (Terminal User Interface) - Interactive, menu-driven
2. **CLI** (Command Line Interface) - Scriptable, traditional commands

### TUI Mode (Interactive)

Launch the TUI by running without arguments:

```bash
sudo portly
```

Or explicitly:

```bash
sudo portly tui
```

#### TUI Navigation

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Navigate up/down |
| `Enter` | Select / Confirm |
| `Tab` | Next field |
| `Esc` | Go back |
| `Ctrl+C` | Quit |

#### Main Menu Options

1. **Add NAT Rule** - Create port forwarding rules
2. **List Rules** - View and manage NAT rules
3. **Firewall** - Start/stop or install firewall service
4. **Security** - Manage SELinux/AppArmor policies
5. **System Status** - View system and provider status
6. **Check Configuration** - Verify system configuration
7. **Quit** - Exit Portly

#### Adding a NAT Rule (TUI)

1. Launch: `sudo portly`
2. Select **"Add NAT Rule"**
3. **Product Field**: Press `Ctrl+D` to show dropdown
4. Select a product (e.g., `podman`) - ports auto-fill!
5. Modify ports if needed
6. Press `Enter` to submit

#### Opening a Port (TUI)

1. Launch: `sudo portly`
2. Select **"Add NAT Rule"** (or use CLI for direct port opening)
3. For firewall-only rules, use the CLI:
   ```bash
   sudo portly open-port --port 8080
   ```

#### Firewall Management (TUI)

1. Launch: `sudo portly`
2. Select **"Firewall"**
3. View current status
4. Press `1` to start, `2` to stop, or `i` to install

#### Security Management (TUI)

1. Launch: `sudo portly`
2. Select **"Security"**
3. On RHEL: Toggle SELinux enforcing/permissive
4. On Ubuntu: Enable/disable AppArmor

### CLI Mode (Scriptable)

Portly supports traditional command-line operations:

#### NAT Rules (Port Forwarding)

```bash
# Show help
portly --help

# Add a NAT rule (forward external port to internal destination)
sudo portly add-nat --product podman --port 8080 --to 10.88.0.1:80

# Add with all options
sudo portly add-nat \
  --product nginx \
  --port 443 \
  --to 192.168.1.100:8443 \
  --protocol tcp \
  --description "HTTPS to nginx container"

# List NAT rules
portly list

# Remove a NAT rule by ID
sudo portly remove-nat --id abc123

# Remove by product and port
sudo portly remove-nat --product podman --port 8080
```

#### Firewall Rules (Opening Ports)

```bash
# Open a port for all IPs
sudo portly open-port --port 8080

# Open a port with specific protocol
sudo portly open-port --port 8080 --protocol tcp

# Open port only for specific IP (IP-based access control)
sudo portly open-port --port 5432 --source-ip 192.168.1.100 --product postgres

# Open UDP port for specific IP
sudo portly open-port --port 53 --protocol udp --source-ip 10.0.0.0/24

# List all open ports
portly list-ports

# Close a port by rule ID
sudo portly close-port abc123
```

#### Firewall Service Management

```bash
# Check firewall status
portly firewall status

# Start firewall
sudo portly firewall start

# Stop firewall
sudo portly firewall stop

# Install firewall (auto-detects OS)
sudo portly firewall install
```

#### Security Management

```bash
# Check security status (SELinux/AppArmor)
portly security status

# SELinux commands (RHEL)
portly security selinux status
sudo portly security selinux enforcing
sudo portly security selinux permissive

# AppArmor commands (Ubuntu)
portly security apparmor status
sudo portly security apparmor enable
sudo portly security apparmor disable
```

#### Other Commands

```bash
# Check configuration
portly check

# Check specific port
portly check --port 8080

# Clean up failed rules
sudo portly rollback

# Show version
portly --version
```

#### CLI Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `add-nat` | Add NAT/port forwarding rule | `sudo portly add-nat --product docker --port 8080 --to 172.17.0.2:80` |
| `remove-nat` | Remove NAT rule | `sudo portly remove-nat --id abc123` |
| `list` | List NAT rules | `portly list --product podman` |
| `open-port` | Open firewall port | `sudo portly open-port --port 8080 --source-ip 192.168.1.100` |
| `close-port` | Close firewall port | `sudo portly close-port abc123` |
| `list-ports` | List open ports | `portly list-ports` |
| `firewall` | Firewall service management | `sudo portly firewall start` |
| `security` | Security management | `sudo portly security selinux enforcing` |
| `check` | Verify config | `portly check --port 8080` |
| `rollback` | Clean failed rules | `sudo portly rollback` |
| `status` | System status | `portly status` |
| `tui` | Launch TUI | `sudo portly tui` |

#### add-nat Flags

| Flag | Required | Description | Example |
|------|----------|-------------|---------|
| `--product` | Yes | Service name | `--product podman` |
| `--port` | Yes | External port | `--port 8080` |
| `--to` | Yes | Target (IP:port) | `--to 10.88.0.1:80` |
| `--internal-port` | Alternative | Internal port only | `--internal-port 80` |
| `--protocol` | No | tcp or udp (default: tcp) | `--protocol udp` |
| `--description` | No | Rule description | `--description "Web server"` |
| `--auto-install` | No | Auto-install missing products | `--auto-install` |
| `--no-security` | No | Skip security policies | `--no-security` |

#### open-port Flags

| Flag | Required | Description | Example |
|------|----------|-------------|---------|
| `--port` | Yes | Port to open | `--port 8080` |
| `--protocol` | No | tcp or udp (default: tcp) | `--protocol tcp` |
| `--source-ip` | No | Limit to specific IP/CIDR | `--source-ip 192.168.1.100` |
| `--product` | No | Product name (default: custom) | `--product nginx` |
| `--description` | No | Rule description | `--description "API server"` |

## Product Database

Portly includes pre-configured settings for popular services:

| Product | Description | Suggested Ports |
|---------|-------------|-----------------|
| `podman` | Container engine | 8080, 8443 |
| `docker` | Container platform | 8080, 443 |
| `tailscale` | VPN mesh network | 41641 |
| `headscale` | Self-hosted Tailscale | 8080 |
| `twingate` | Zero trust network | 443 |
| `steam` | Gaming platform | 27015 |
| `minecraft` | Minecraft server | 25565 |
| `nginx` | Web server | 80, 443 |
| `postgres` | PostgreSQL database | 5432 |
| `redis` | Redis cache | 6379 |
| `custom` | Type your own | any |

## Examples

### Example 1: Expose Podman Container (TUI)

```bash
sudo portly
# Select "Add NAT Rule"
# Product: Select "podman" (ports auto-fill to 8080)
# External Port: 80 (change from 8080)
# Internal IP: 10.88.0.1
# Internal Port: 8080
# Protocol: tcp
# Press Enter to submit
```

### Example 2: Open Port with IP Restriction (CLI)

```bash
# Allow PostgreSQL access only from specific IP
sudo portly open-port \
  --port 5432 \
  --source-ip 192.168.1.100 \
  --product postgres \
  --description "DB access for app server"
```

### Example 3: Manage Firewall (CLI)

```bash
# Check status
portly firewall status

# Start firewall if stopped
sudo portly firewall start

# Or install if missing
sudo portly firewall install
```

### Example 4: Configure SELinux (RHEL)

```bash
# Check current status
portly security selinux status

# Set to permissive mode for testing
sudo portly security selinux permissive

# Set back to enforcing
sudo portly security selinux enforcing
```

### Example 5: Full Workflow - Secure Database Access

```bash
# 1. Check system
portly status

# 2. Start firewall if needed
sudo portly firewall start

# 3. Set SELinux to permissive (if on RHEL)
sudo portly security selinux permissive

# 4. Open PostgreSQL port with IP restriction
sudo portly open-port \
  --port 5432 \
  --source-ip 10.0.0.0/24 \
  --product postgres

# 5. Verify port is open
portly list-ports

# 6. Set SELinux back to enforcing
sudo portly security selinux enforcing
```

### Example 6: NAT + Firewall Rule

```bash
# Scenario: Expose nginx container on port 80, 
# and open port 443 directly on firewall

# Add NAT rule for port 80
sudo portly add-nat \
  --product nginx \
  --port 80 \
  --to 10.88.0.1:8080

# Open port 443 directly on firewall
sudo portly open-port \
  --port 443 \
  --product nginx

# List all rules
portly list
portly list-ports
```

## Architecture

Portly is organized into modular components:

```
cmd/orchestrator/          # CLI entry point
├── main.go               # TUI/CLI dispatcher
├── tui.go                # TUI command
├── firewall.go           # Firewall service commands
├── security.go           # Security commands
├── open_port.go          # Open/Close port commands
├── add_nat*.go           # NAT rule commands
└── *.go                  # Other commands

internal/
├── tui/                  # TUI implementation
│   ├── styles/           # Theme and colors
│   ├── model.go          # State management
│   ├── init.go           # Initialization
│   ├── update.go         # Event handlers
│   ├── view.go           # Rendering
│   ├── menu.go           # Main menu
│   ├── form*.go          # Form handling
│   ├── firewall_screen.go    # Firewall UI
│   ├── firewall_ops.go       # Firewall operations
│   ├── security_screen.go    # Security UI
│   ├── security_ops.go       # Security operations
│   ├── list_rules.go         # Rules table
│   └── *_screen.go       # Other screens
├── drivers/              # OS-specific drivers
│   ├── firewalld/        # RHEL/Fedora driver
│   │   ├── driver.go
│   │   ├── nat.go
│   │   ├── list.go
│   │   ├── utils.go
│   │   └── firewall.go   # Port opening
│   ├── nftables/         # Ubuntu/Debian driver
│   │   ├── driver.go
│   │   ├── nat.go
│   │   ├── list.go
│   │   ├── utils.go
│   │   └── firewall.go   # Port opening
│   └── pf/               # macOS driver
│       ├── driver.go
│       ├── nat.go
│       ├── list.go
│       ├── utils.go
│       └── firewall.go   # Port opening
├── security/             # Security policy management
├── platform/             # OS detection
├── state/                # State persistence
└── installer/            # Package installation

pkg/models/               # Data models
├── protocol.go
├── nat_rule.go
├── firewall_rule.go      # Firewall rule model
├── security_policy.go
├── os_info.go
├── product_info.go
├── rule_status.go
└── state.go
```

## Development

### Building

```bash
# Build for current platform
go build -o portly cmd/orchestrator/*.go

# Build for Linux AMD64
GOOS=linux GOARCH=amd64 go build -o portly-linux-amd64 cmd/orchestrator/*.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o portly-darwin-amd64 cmd/orchestrator/*.go
```

### Testing

```bash
go test ./...
```

### Project Structure

All files are kept under 150 lines for maintainability.

## Troubleshooting

### "root privileges required"
Portly needs root to modify firewall rules and security policies. Use `sudo`.

### "no firewall provider available"
Your OS may not have a supported firewall installed. Use `sudo portly firewall install` to install one.

### Port already in use
Run `portly check --port <port>` to verify availability.

### Rules not persisting
Ensure the firewall service is enabled:
- RHEL: `sudo systemctl enable firewalld`
- Ubuntu: `sudo systemctl enable nftables`

### SELinux blocking connections
Temporarily set to permissive mode for testing:
```bash
sudo portly security selinux permissive
```

Remember to set back to enforcing after testing:
```bash
sudo portly security selinux enforcing
```

### IP-based rules not working
Make sure the source IP is correct and reachable. Use CIDR notation for subnets:
```bash
sudo portly open-port --port 8080 --source-ip 192.168.0.0/24
```

## License

MIT License - See LICENSE file for details.

## Contributing

Contributions welcome! Please ensure:
- Files stay under 150 lines
- Follow existing code patterns
- Add tests for new functionality
- Update documentation

## Support

For issues and feature requests, please use GitHub Issues.
# Portly
