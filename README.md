# MCP Navigator (Go)

A comprehensive Model Context Protocol (MCP) client implementation in Go, providing both a powerful CLI and a library for integration into other applications.

## Features

- 🔍 **Server Discovery**: Automatically discover MCP servers on TCP ports, HTTP endpoints, and Docker containers
- 🔌 **Multiple Transports**: Support for TCP, HTTP/SSE, HTTP Streaming, STDIO, and Docker-based connections
- 💬 **Interactive CLI**: Full-featured command-line interface for server interaction
- 🛠️ **Complete MCP Protocol**: Full support for Tools, Resources, and Prompts
- 🐳 **Docker Support**: Direct support for Docker-based MCP servers
- 🌐 **HTTP Support**: Support for both Server-Sent Events (SSE) and streaming HTTP modes
- 📚 **Library Integration**: Use as a library in your Go applications ([Library Documentation](LIBRARY.md))
- ⚡ **High Performance**: Written in Go for speed and efficiency

## Installation

### From Source

```bash
git clone https://github.com/kunalkushwaha/mcp-navigator-go.git
cd mcp-navigator-go
go mod download
go build -o mcp-navigator main.go
```

### Using Go Install

```bash
go install github.com/kunalkushwaha/mcp-navigator-go@latest
```

### As a Library

```bash
go get github.com/kunalkushwaha/mcp-navigator-go
```

## Quick Start

### 1. Discover Available MCP Servers

```bash
./mcp-navigator discover
```

This will show available servers including:
- TCP servers on common MCP ports (8811, 8812, 8813, etc.)
- HTTP MCP servers (both SSE and streaming modes)
- Docker-based MCP servers
- Standard Docker MCP configuration (alpine/socat)

### 2. Interactive Mode

```bash
./mcp-navigator interactive
```

This starts the interactive CLI where you can:
- `help` - Show available commands
- `discover` - Find MCP servers
- `connect <server-name-or-index>` - Connect to a server
- `list-tools` - List available tools on connected server
- `list-resources` - List available resources on connected server
- `call-tool <tool-name> [json-args]` - Execute a tool
- `status` - Show connection status
- `exit` - Exit the client

### 3. Direct Commands

Connect to a TCP server:
```bash
./mcp-navigator connect --tcp --host localhost --port 8811
```

Connect to HTTP MCP server (SSE mode):
```bash
./mcp-navigator connect --http --url http://localhost:8812 --endpoint /sse/
```

Connect to HTTP MCP server (streaming mode):
```bash
./mcp-navigator connect --http --url http://localhost:8813 --endpoint /mcp
```

Connect to Docker MCP server:
```bash
./mcp-navigator connect --docker
```

Execute a tool directly:
```bash
./mcp-navigator tool --name search --arguments '{"query": "golang"}' --docker
```

## Library Usage

MCP Navigator can be used as a library in your Go applications. 

📚 **Documentation:**
- [**Library Documentation**](LIBRARY.md) - Comprehensive guide with examples and best practices
- [**API Reference**](API.md) - Concise API documentation and usage patterns

### Quick Library Example

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
)

func main() {
    // Create client with TCP transport
    mcpClient := client.NewMCPClient(transport.NewTCPTransport("localhost", 8811))
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Connect and initialize
    err := mcpClient.Connect(ctx)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer mcpClient.Close()
    
    serverInfo, err := mcpClient.Initialize(ctx, "my-app", "1.0.0")
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    
    log.Printf("Connected to: %s %s", serverInfo.Name, serverInfo.Version)
    
    // List and execute tools
    tools, err := mcpClient.ListTools(ctx)
    if err != nil {
        log.Fatalf("Failed to list tools: %v", err)
    }
    
    for _, tool := range tools {
        log.Printf("Available tool: %s", tool.Name)
    }
}
```

📖 **See [LIBRARY.md](LIBRARY.md) for:**
- Complete transport examples (TCP, HTTP SSE, HTTP Streaming, STDIO, WebSocket)
- Server discovery integration
- Advanced usage patterns (builder pattern, service wrappers, error handling)
- Testing strategies and mock implementations
- Performance optimization techniques

📋 **See [API.md](API.md) for:**
- Complete API reference with all methods and types
- Quick usage patterns and code snippets
- Error handling strategies
- Configuration options

## CLI Usage

The following examples show CLI usage. For library integration, see [Library Documentation](LIBRARY.md).

### Server Discovery

```bash
# Discover all servers
./mcp-navigator discover

# Scan specific host
./mcp-navigator discover --host 192.168.1.100

# Only scan TCP ports
./mcp-navigator discover --tcp-only

# Only check Docker containers
./mcp-navigator discover --docker-only

# Custom port range
./mcp-navigator discover --start-port 8000 --end-port 9000
```

### Connecting to Servers

```bash
# TCP connection
./mcp-navigator connect --tcp --host localhost --port 8811

# HTTP connection (SSE mode)
./mcp-navigator connect --http --url http://localhost:8812 --endpoint /sse/

# HTTP connection (streaming mode) 
./mcp-navigator connect --http --url http://localhost:8813 --endpoint /mcp

# STDIO connection
./mcp-navigator connect --stdio --command "node" --args "server.js"

# Docker connection (uses alpine/socat bridge)
./mcp-navigator connect --docker

# With custom timeout
./mcp-navigator connect --tcp --host localhost --port 8811 --timeout 45s
```

### Tool Execution

```bash
# Execute tool with JSON arguments
./mcp-navigator tool --name search --arguments '{"query": "golang", "limit": 10}' --tcp

# Execute tool via Docker MCP server
./mcp-navigator tool --name docker --arguments '{"command": "ps -a"}' --docker

# Execute tool with no arguments
./mcp-navigator tool --name list-files --docker
```

## Docker MCP Server Support

The client automatically supports the standard Docker-based MCP server configuration used by Claude Desktop:

```json
{
  "mcpServers": {
    "MCP_DOCKER": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "alpine/socat",
        "STDIO",
        "TCP:host.docker.internal:8811"
      ]
    }
  }
}
```

This configuration allows MCP servers running in Docker containers to communicate with external TCP services.

## Interactive Mode Example

```bash
$ ./mcp-navigator interactive

🚀 MCP Navigator Interactive Mode
Type 'help' for available commands.
🔍 Discovering MCP servers...
✅ Found 6 server(s):
  1. TCP Server localhost:8811 (tcp)
     Address: localhost:8811
  2. TCP Server localhost:8812 (tcp)
     Address: localhost:8812
  3. TCP Server localhost:8813 (tcp)
     Address: localhost:8813
  4. HTTP MCP Server (SSE) localhost:8812/sse/ (http)
     Address: localhost:8812
  5. HTTP MCP Server (Streaming) localhost:8813/mcp (http)
     Address: localhost:8813
  6. Docker MCP (Direct TCP) (docker)
     Address: localhost:8811

mcp-client> help

📋 Available Commands:
  help              - Show this help message
  discover          - Discover available MCP servers
  connect <n>       - Connect to a server by name or index
  disconnect        - Disconnect from current server
  list-tools        - List tools available on current server
  list-resources    - List resources available on current server
  call-tool <n> [args] - Call a tool with optional JSON arguments
  status            - Show connection status
  exit/quit         - Exit the client

mcp-client> connect 5
🔌 Connecting to HTTP MCP Server (Streaming) localhost:8813/mcp...
✅ Connected to HTTP MCP Server (Streaming) localhost:8813/mcp
🚀 Server: Docker AI MCP Gateway 2.0.1

mcp-client> list-tools
📝 Available tools (4):
  1. docker
     Description: use the docker cli
  2. fetch
     Description: Fetches a URL from the internet and optionally extracts its contents as markdown
  3. fetch_content
     Description: Fetch and parse content from a webpage URL
  4. search
     Description: Search DuckDuckGo and return formatted results

mcp-client> call-tool search {"query": "Model Context Protocol"}
🔧 Calling tool: search
📤 Tool result:
The Model Context Protocol (MCP) is an open standard that enables secure connections between AI assistants and data sources...

mcp-client> status

📊 Status:
  Available servers: 6
  Current connection: HTTP MCP Server (Streaming) localhost:8813/mcp ✅
  Server info: Docker AI MCP Gateway 2.0.1

mcp-client> exit

👋 Shutting down MCP client...
✅ Goodbye!
```

## Configuration

### Environment Variables

- `MCP_CLIENT_TIMEOUT`: Default timeout for operations (default: 30s)
- `MCP_CLIENT_HOST`: Default host for TCP connections (default: localhost)
- `MCP_CLIENT_PORT`: Default port for TCP connections (default: 8811)
- `MCP_CLIENT_URL`: Default URL for HTTP connections (default: http://localhost:8812)
- `MCP_CLIENT_ENDPOINT`: Default endpoint for HTTP connections (default: /mcp)
- `MCP_CLIENT_VERBOSE`: Enable verbose logging (default: false)

### Config File

Create `~/.mcp-client.yaml`:

```yaml
# Default connection settings
host: localhost
port: 8811
timeout: 30s

# Discovery settings
discovery:
  startPort: 8810
  endPort: 8820
  timeout: 5s

# Logging
verbose: false
```

## Command Reference

### Global Flags

- `--config`: Config file path
- `--verbose, -v`: Enable verbose output

### Commands

#### `discover`
Discover available MCP servers

**Flags:**
- `--host`: Host to scan (default: localhost)
- `--start-port`: Start port for scanning (default: 8810)
- `--end-port`: End port for scanning (default: 8820)
- `--timeout`: Discovery timeout (default: 5s)
- `--tcp-only`: Only scan TCP ports
- `--docker-only`: Only check Docker containers

#### `connect`
Connect to an MCP server and show available tools/resources

**Flags:**
- `--type`: Connection type (tcp, stdio, docker)
- `--tcp, -t`: Use TCP transport
- `--stdio, -s`: Use STDIO transport
- `--docker, -d`: Use Docker transport
- `--host`: TCP host (default: localhost)
- `--port`: TCP port (default: 8811)
- `--command`: Command for STDIO transport
- `--args`: Arguments for STDIO command
- `--timeout`: Connection timeout (default: 30s)

#### `tool`
Execute a specific tool on an MCP server

**Flags:**
- `--name`: Tool name (required)
- `--arguments`: JSON arguments for the tool (default: "{}")
- All connection flags from `connect` command

#### `interactive`
Start interactive mode

**Aliases:** `i`, `shell`

## Development

### Project Structure

```
├── main.go                 # Entry point
├── internal/
│   └── cli/               # CLI commands
│       ├── root.go        # Root command and configuration
│       ├── discover.go    # Server discovery command
│       ├── connect.go     # Connection command
│       ├── tool.go        # Tool execution command
│       └── interactive.go # Interactive mode
├── pkg/
│   ├── client/           # MCP client implementation
│   ├── discovery/        # Server discovery logic
│   ├── mcp/             # MCP protocol types and utilities
│   └── transport/       # Transport implementations (TCP, STDIO, WebSocket, HTTP)
├── go.mod
└── README.md
```

### Building

```bash
go build -o mcp-client main.go
```

### Testing

```bash
go test ./...
```

### Running in Development

```bash
go run main.go interactive
```

## Requirements

- Go 1.21 or later
- Docker (for Docker-based MCP servers)
- MCP Server running on TCP port or Docker

## Related

- [Library Documentation](LIBRARY.md) - Comprehensive guide for using MCP Navigator as a library
- [API Reference](API.md) - Concise API documentation and examples
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [TypeScript MCP SDK](https://github.com/modelcontextprotocol/typescript-sdk)

## License

MIT License - see LICENSE file for details.
