# MCP Navigator Examples

This directory contains comprehensive examples demonstrating how to use MCP Navigator as a library in your Go applications.

## 📁 Structure

### 🟢 Basic Examples (`basic/`)
Perfect for getting started with MCP Navigator library integration.

- **[simple-client](basic/simple-client/)** - Minimal TCP client example
- **[transport-examples](basic/transport-examples/)** - Different transport types (TCP, HTTP SSE, HTTP Streaming)
- **[discovery](basic/discovery/)** - Server discovery functionality

### 🟡 Advanced Examples (`advanced/`)
Demonstrates advanced patterns and production-ready techniques.

- **[builder-pattern](advanced/builder-pattern/)** - Builder pattern and service wrapper implementations
- **[error-handling](advanced/error-handling/)** - Robust error handling and retry strategies
- **[concurrent](advanced/concurrent/)** - Concurrent operations and performance optimization

### 🔴 Complete Examples (`complete/`)
Full protocol demonstrations showcasing all MCP features.

- **[full-protocol](complete/full-protocol/)** - Tools, Resources, and Prompts usage

### 🔧 CLI Examples (`cmd/`)
Command-line interface examples and utilities.

- **[complete-protocol](cmd/complete-protocol/)** - Complete protocol CLI
- **[docker-tools](cmd/docker-tools/)** - Docker integration tools
- **[fixed-docker](cmd/fixed-docker/)** - Fixed Docker configuration
- **[tcp-direct](cmd/tcp-direct/)** - Direct TCP connection

## 🚀 Quick Start

### 1. Basic Library Usage
```bash
cd basic/simple-client
go run main.go
```

### 2. Transport Examples
```bash
cd basic/transport-examples
go run main.go
```

### 3. Server Discovery
```bash
cd basic/discovery
go run main.go
```

## 📚 Learning Path

1. **Start with [basic/simple-client](basic/simple-client/)** - Learn the fundamentals
2. **Explore [basic/transport-examples](basic/transport-examples/)** - Understand different connection types
3. **Try [basic/discovery](basic/discovery/)** - Learn server discovery
4. **Move to [advanced/builder-pattern](advanced/builder-pattern/)** - Advanced patterns
5. **Study [advanced/error-handling](advanced/error-handling/)** - Production techniques
6. **Review [complete/full-protocol](complete/full-protocol/)** - Complete MCP protocol

## 🔧 Running Examples

All examples use the `//go:build example` build tag to avoid conflicts. Run them with:

```bash
go run -tags example main.go
```

Or build with:

```bash
go build -tags example -o example main.go
./example
```

## 📋 Prerequisites

- Go 1.21 or later
- MCP server running (for testing connections)
- Docker (for Docker-based examples)

## 🎯 Use Cases

| Example | Use Case | Complexity |
|---------|----------|------------|
| `basic/simple-client` | Quick integration | 🟢 Beginner |
| `basic/transport-examples` | Multi-transport support | 🟢 Beginner |
| `basic/discovery` | Dynamic server finding | 🟢 Beginner |
| `advanced/builder-pattern` | Production apps | 🟡 Intermediate |
| `advanced/error-handling` | Robust applications | 🟡 Intermediate |
| `advanced/concurrent` | High-performance apps | 🟡 Intermediate |
| `complete/full-protocol` | Full MCP features | 🔴 Advanced |

## 📖 Documentation

- [Library Documentation](../LIBRARY.md) - Comprehensive guide
- [API Reference](../API.md) - Quick API reference
- [Main README](../README.md) - Project overview

## 🤝 Contributing

Want to add more examples? See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.
