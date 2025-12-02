# 🔐 arcadio

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**arcadio** (short: `arc`) is a lightweight, secure CLI tool for creating encrypted document arcs with fast search capabilities. Built in Go for speed and portability.

```bash
# Create an encrypted arc
$ arc create work-docs

# Add documents with tags
$ arc add work-docs contract.pdf --tags legal,important

# Search instantly
$ arc search work-docs invoice

# Export when needed
$ arc export work-docs <doc-id> output.pdf
```

## Features

- **Encryption**: AES-256-GCM with Argon2id key derivation
- **Fast & lightweight**: Built in Go, single binary, no dependencies
- **Tag-based organization**: Organize documents with flexible tagging
- **Fuzzy search**: Find documents quickly by filename
- **Portable**: Export entire arcs as encrypted archives
- **Remote sync**: Share arcs securely over mTLS (coming soon)
- **Simple CLI**: Intuitive commands, memorable syntax

## Quick Start

### Installation

```bash
# Install from source
go install github.com/ViniTamanhao/arcadio@latest

# Or build locally
git clone https://github.com/ViniTamanhao/arcadio.git
cd arcadio
go build -o arc
sudo mv arc /usr/local/bin/
```

### Basic Usage

```bash
# Create a new encrypted arc
arc create my-arc

# Add documents
arc add my-arc document.pdf --tags work,finance
arc add my-arc photos/ --recursive --tags vacation

# List documents
arc docs my-arc

# Search
arc search my-arc invoice

# Export a document
arc export my-arc <doc-id> output.pdf

# Get arc info
arc info my-arc

# List all arcs
arc list
```

### Commands

#### Arc Management

| Command | Description | Example |
|---------|-------------|---------|
| `arc create <name>` | Create a new encrypted arc | `arc create work-docs` |
| `arc list` | List all arcs | `arc list` |
| `arc info <arc>` | Show arc information | `arc info work-docs` |
| `arc delete <arc>` | Delete a arc permanently | `arc delete old-project` |

#### Document Operations

| Command | Description | Example |
|---------|-------------|---------|
| `arc add <arc> <file>` | Add document(s) to arc | `arc add work-docs file.pdf` |
| `arc add <arc> <dir> -r` | Add directory recursively | `arc add work-docs docs/ -r` |
| `arc docs <arc>` | List all documents | `arc docs work-docs` |
| `arc remove <arc> <doc-id>` | Remove a document | `arc remove work-docs abc123...` |
| `arc export <arc> <doc-id> <out>` | Export a document | `arc export work-docs abc123 file.pdf` |
| `arc search <arc> <query>` | Search documents | `arc search work-docs invoice` |
| `arc tag <arc> <doc-id> <tags>` | Add tags to document | `arc tag work-docs abc123,urgent` |


#### Adding Documents with Tags

```bash
# Single file with tags
arc add work-docs contract.pdf --tags legal,2024,important

# Entire directory with tags
arc add photos vacation-2024/ --recursive --tags travel,family,2024
```

#### Searching and Filtering

```bash
# Search by filename
arc search work-docs invoice

# List all documents (then filter by tags manually)
arc docs work-docs | grep "legal"
```

#### Arc Management

```bash
# Get detailed arc info
arc info work-docs

# Delete with confirmation
arc delete old-project

# Force delete (skip confirmation)
arc delete temp-arc --force
```

## 🏗️ Architecture

### Storage Structure

```
~/.arcadio/
├── registry.json          # Arc name → ID mappings
└── arcs/
    └── <arc-uuid>/
        ├── arc.sec        # Security config (salt, hashes)
        ├── arc.meta       # Encrypted arc metadata
        └── documents/
            ├── <doc-uuid-1>.bin
            ├── <doc-uuid-2>.bin
            └── ...
```

### Encryption Flow

```
Password → Argon2id → Encryption Key → AES-256-GCM → Encrypted Documents
```

1. **Key Derivation**: Argon2id transforms password into 256-bit key
2. **Encryption**: AES-256-GCM encrypts each document individually
3. **Integrity**: SHA-256 hashes verify document integrity
4. **Authentication**: GCM provides authenticated encryption

### Security Features

- **AES-256-GCM**: Industry-standard authenticated encryption
- **Argon2id**: Memory-hard key derivation (resistant to GPU attacks)
- **SHA-256**: Document integrity verification
- **Random nonces**: Unique nonce per encryption operation
- **Secure key storage**: Keys never written to disk

## 🔒 Security

### What's Encrypted

- ✅ All document content
- ✅ Arc metadata (document names, tags, etc.)
- ✅ Document filenames and metadata

### What's Not Encrypted

- ❌ Arc IDs (random UUIDs - not sensitive)
- ❌ Arc names (organizational labels)
- ❌ Security configuration (only contains hashes and salt)

### Best Practices

1. **Use strong passwords**: Minimum 12 characters, mixed case, numbers, symbols
2. **Don't reuse passwords**: Each arc should have a unique password
3. **Backup regularly**: Export arc to encrypted archives
4. **Store passwords securely**: Use a password manager
5. **Verify exports**: Check file hashes after export

### Threat Model

**arcadio protects against:**
- Unauthorized access to local files
- Data breaches (all data encrypted at rest)
- Password guessing (Argon2id is memory-hard)
- Data tampering (GCM authentication)

**arcadio does NOT protect against:**
- Keyloggers or malware on your system
- Physical access to unlocked system
- Weak passwords or password reuse
- Loss of password (encryption is irrecoverable)

## 🛣️ Roadmap

### Current Status: v0.1.0-alpha

- [x] Core encryption and arc management
- [x] Document operations (add, remove, export)
- [x] Tag-based organization
- [x] Fuzzy search
- [x] Arc registry for name management

### Upcoming Features

#### v0.2.0 - Export/Import
- [ ] Export arc to single encrypted file (.arcx)
- [ ] Import arc from archive
- [ ] Arc compression
- [ ] Backup and restore functionality

#### v0.3.0 - Remote Sync
- [ ] HTTP server mode (`arc serve`)
- [ ] mTLS authentication
- [ ] Remote arc access
- [ ] Push/pull synchronization
- [ ] Multi-user support

#### v0.4.0 - Enhanced Features
- [ ] Document compression (before encryption)
- [ ] Better fuzzy search (using library)
- [ ] File deduplication
- [ ] Batch operations
- [ ] Progress bars for large files

#### v1.0.0 - Production Ready
- [ ] TUI (Terminal UI) mode
- [ ] Mount arc as filesystem (FUSE)
- [ ] Plugin system
- [ ] Cloud storage backends (S3, etc.)
- [ ] Comprehensive test suite
- [ ] Performance benchmarks

## 🧪 Development

### Prerequisites

- Go 1.21 or higher
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/ViniTamanhao/arcadio.git
cd arcadio

# Install dependencies
go mod download

# Build
go build -o arc

# Run tests
go test ./...

# Install locally
go install
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow standard Go formatting (`gofmt`)
- Write tests for new features
- Update documentation as needed
- Keep commits atomic and well-described

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Argon2](https://github.com/P-H-C/phc-winner-argon2) - Password hashing

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/ViniTamanhao/arcadio/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/ViniTamanhao/arcadio/discussions)
- 📧 **Email**: vtamanhao@gmail.com

## ⚠️ Disclaimer

This software is provided "as is", without warranty of any kind. Always maintain backups of important data. The encryption is strong, but if you lose your password, your data is irrecoverable.

---

**Made with ❤️ by ViniTamanhao**

⭐ Star this repo if you find it useful!
