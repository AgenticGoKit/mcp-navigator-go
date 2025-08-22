# Basic Examples

These examples demonstrate fundamental MCP Navigator library usage patterns. Perfect for getting started with MCP integration.

## 📁 Examples

### [simple-client](simple-client/)
**Minimal TCP client example**
- Basic connection setup
- Protocol initialization  
- Simple tool execution
- Resource cleanup

**When to use:** First-time integration, learning the basics

### [transport-examples](transport-examples/)
**Different transport types**
- TCP transport
- HTTP SSE (Server-Sent Events)
- HTTP Streaming
- WebSocket transport

**When to use:** Need to support multiple connection types

### [discovery](discovery/)
**Server discovery functionality**
- Automatic server detection
- TCP port scanning
- HTTP endpoint discovery
- Docker container discovery

**When to use:** Dynamic server discovery, multi-environment deployments

## 🚀 Quick Start

1. **Start here:** [simple-client](simple-client/)
2. **Learn transports:** [transport-examples](transport-examples/)
3. **Add discovery:** [discovery](discovery/)

## ⚙️ Prerequisites

- Go 1.21+
- MCP server running (any of: TCP:8811, HTTP:8812/8813)
- Docker (for Docker examples)

## 🏃 Running Examples

```bash
# Navigate to any example directory
cd simple-client

# Run the example
go run -tags example main.go
```

## 📖 Next Steps

After mastering these basics, explore:
- [Advanced Examples](../advanced/) - Production patterns
- [Complete Examples](../complete/) - Full protocol features
