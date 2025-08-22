# MCP Navigator Go Library

This document provides comprehensive guidance for using MCP Navigator as a library in your Go applications to interact with Model Context Protocol (MCP) servers.

## Installation

Add MCP Navigator to your Go project:

```bash
go get github.com/kunalkushwaha/mcp-navigator-go
```

## Quick Start

Here's a simple example to get started:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
)

func main() {
    // Create a TCP transport
    tcpTransport := transport.NewTCPTransport("localhost", 8811)
    
    // Create MCP client
    mcpClient := client.NewMCPClient(tcpTransport)
    
    // Connect to server
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    err := mcpClient.Connect(ctx)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer mcpClient.Close()
    
    // Initialize MCP protocol
    serverInfo, err := mcpClient.Initialize(ctx, "my-app", "1.0.0")
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    
    fmt.Printf("Connected to: %s %s\n", serverInfo.Name, serverInfo.Version)
    
    // List available tools
    tools, err := mcpClient.ListTools(ctx)
    if err != nil {
        log.Fatalf("Failed to list tools: %v", err)
    }
    
    fmt.Printf("Available tools: %d\n", len(tools))
    for _, tool := range tools {
        fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
    }
}
```

## Core Components

### 1. Transports

MCP Navigator supports multiple transport protocols:

#### TCP Transport

```go
import "github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"

// Create TCP transport
tcpTransport := transport.NewTCPTransport("localhost", 8811)
```

#### HTTP SSE Transport

```go
// Create HTTP SSE transport for real-time communication
sseTransport := transport.NewSSETransport("http://localhost:8812", "/sse/")
```

#### HTTP Streaming Transport

```go
// Create HTTP streaming transport for request-response communication
streamingTransport := transport.NewStreamingHTTPTransport("http://localhost:8813", "/mcp")
```

#### STDIO Transport

```go
// Create STDIO transport for process-based MCP servers
stdioTransport := transport.NewSTDIOTransport("node", []string{"server.js"})
```

#### WebSocket Transport

```go
// Create WebSocket transport
wsTransport := transport.NewWebSocketTransport("ws://localhost:8814/mcp")
```

### 2. Client Builder Pattern

Use the builder pattern for complex client configurations:

```go
import (
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
)

// Build client with custom configuration
mcpClient := client.NewMCPClientBuilder().
    WithTransport(transport.NewTCPTransport("localhost", 8811)).
    WithTimeout(30 * time.Second).
    WithClientInfo("my-app", "1.0.0").
    WithRetryPolicy(3, 5*time.Second).
    Build()
```

### 3. Server Discovery

Automatically discover available MCP servers:

```go
import "github.com/kunalkushwaha/mcp-navigator-go/pkg/discovery"

func discoverServers() {
    discoverer := discovery.NewDiscoverer()
    
    // Discover all types of servers
    servers, err := discoverer.DiscoverAll(context.Background())
    if err != nil {
        log.Fatalf("Discovery failed: %v", err)
    }
    
    for _, server := range servers {
        fmt.Printf("Found server: %s (%s)\n", server.Name, server.Type)
        fmt.Printf("  Address: %s\n", server.Address)
        fmt.Printf("  Description: %s\n", server.Description)
    }
}
```

#### Targeted Discovery

```go
// Discover only TCP servers
tcpServers, err := discoverer.DiscoverTCP(ctx, "localhost", 8810, 8820)

// Discover only HTTP servers
httpServers, err := discoverer.DiscoverHTTP(ctx, "localhost")

// Discover only Docker servers
dockerServers, err := discoverer.DiscoverDocker(ctx)
```

## Advanced Usage Examples

### 1. Tool Execution

```go
func executeTools(mcpClient *client.MCPClient) {
    ctx := context.Background()
    
    // Execute tool with arguments
    args := map[string]interface{}{
        "query": "golang",
        "limit": 10,
    }
    
    result, err := mcpClient.CallTool(ctx, "search", args)
    if err != nil {
        log.Printf("Tool execution failed: %v", err)
        return
    }
    
    fmt.Printf("Tool result: %v\n", result)
}
```

### 2. Resource Management

```go
func manageResources(mcpClient *client.MCPClient) {
    ctx := context.Background()
    
    // List available resources
    resources, err := mcpClient.ListResources(ctx)
    if err != nil {
        log.Printf("Failed to list resources: %v", err)
        return
    }
    
    for _, resource := range resources {
        fmt.Printf("Resource: %s\n", resource.URI)
        
        // Read resource content
        content, err := mcpClient.ReadResource(ctx, resource.URI)
        if err != nil {
            log.Printf("Failed to read resource %s: %v", resource.URI, err)
            continue
        }
        
        fmt.Printf("Content: %s\n", content)
    }
}
```

### 3. Error Handling and Retries

```go
import (
    "time"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
)

func robustConnection() {
    // Create client with retry configuration
    mcpClient := client.NewMCPClientBuilder().
        WithTransport(transport.NewTCPTransport("localhost", 8811)).
        WithRetryPolicy(3, 5*time.Second).
        WithErrorHandler(func(err error) {
            log.Printf("MCP error: %v", err)
        }).
        Build()
    
    // Connection with exponential backoff
    ctx := context.Background()
    for i := 0; i < 3; i++ {
        err := mcpClient.Connect(ctx)
        if err == nil {
            break
        }
        
        log.Printf("Connection attempt %d failed: %v", i+1, err)
        time.Sleep(time.Duration(i+1) * 2 * time.Second)
    }
}
```

### 4. Concurrent Operations

```go
import "sync"

func concurrentToolExecution(mcpClient *client.MCPClient) {
    ctx := context.Background()
    var wg sync.WaitGroup
    
    tools := []string{"search", "docker", "fetch"}
    results := make(chan string, len(tools))
    
    for _, toolName := range tools {
        wg.Add(1)
        go func(tool string) {
            defer wg.Done()
            
            result, err := mcpClient.CallTool(ctx, tool, map[string]interface{}{})
            if err != nil {
                results <- fmt.Sprintf("Error executing %s: %v", tool, err)
                return
            }
            
            results <- fmt.Sprintf("Tool %s result: %v", tool, result)
        }(toolName)
    }
    
    wg.Wait()
    close(results)
    
    for result := range results {
        fmt.Println(result)
    }
}
```

### 5. Custom Transport Implementation

```go
import "github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"

// Implement custom transport
type CustomTransport struct {
    // Your custom fields
}

func (t *CustomTransport) Connect(ctx context.Context) error {
    // Implement connection logic
    return nil
}

func (t *CustomTransport) Send(ctx context.Context, message []byte) error {
    // Implement message sending
    return nil
}

func (t *CustomTransport) Receive(ctx context.Context) ([]byte, error) {
    // Implement message receiving
    return nil, nil
}

func (t *CustomTransport) Close() error {
    // Implement cleanup
    return nil
}

// Use custom transport
func useCustomTransport() {
    customTransport := &CustomTransport{}
    mcpClient := client.NewMCPClient(customTransport)
    // ... use client
}
```

## Configuration

### Environment Variables

Your application can use the same environment variables:

```go
import "os"

func configureFromEnv() *client.MCPClient {
    host := os.Getenv("MCP_CLIENT_HOST")
    if host == "" {
        host = "localhost"
    }
    
    port := os.Getenv("MCP_CLIENT_PORT")
    if port == "" {
        port = "8811"
    }
    
    // Convert port to int and create transport
    // ... implementation
}
```

### Configuration Struct

```go
type MCPConfig struct {
    Host        string        `yaml:"host"`
    Port        int           `yaml:"port"`
    Timeout     time.Duration `yaml:"timeout"`
    RetryCount  int           `yaml:"retry_count"`
    RetryDelay  time.Duration `yaml:"retry_delay"`
    ClientName  string        `yaml:"client_name"`
    Version     string        `yaml:"version"`
}

func loadConfig(path string) (*MCPConfig, error) {
    // Load configuration from YAML file
    // ... implementation
}
```

## Best Practices

### 1. Resource Management

```go
func properResourceManagement() {
    mcpClient := client.NewMCPClient(transport.NewTCPTransport("localhost", 8811))
    
    // Always use context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Always close connections
    defer func() {
        if err := mcpClient.Close(); err != nil {
            log.Printf("Error closing client: %v", err)
        }
    }()
    
    // Connect and use client
    if err := mcpClient.Connect(ctx); err != nil {
        log.Fatalf("Connection failed: %v", err)
    }
}
```

### 2. Error Handling

```go
import "github.com/kunalkushwaha/mcp-navigator-go/pkg/client"

func handleErrors(mcpClient *client.MCPClient) {
    ctx := context.Background()
    
    tools, err := mcpClient.ListTools(ctx)
    if err != nil {
        switch {
        case client.IsConnectionError(err):
            log.Printf("Connection error: %v", err)
            // Attempt reconnection
        case client.IsProtocolError(err):
            log.Printf("Protocol error: %v", err)
            // Handle protocol issues
        case client.IsTimeoutError(err):
            log.Printf("Timeout error: %v", err)
            // Retry with longer timeout
        default:
            log.Printf("Unknown error: %v", err)
        }
        return
    }
    
    // Use tools
    for _, tool := range tools {
        fmt.Printf("Tool: %s\n", tool.Name)
    }
}
```

### 3. Logging Integration

```go
import (
    "log/slog"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
)

func withStructuredLogging() {
    // Configure structured logging
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    
    mcpClient := client.NewMCPClientBuilder().
        WithTransport(transport.NewTCPTransport("localhost", 8811)).
        WithLogger(logger).
        Build()
    
    // Client will now use structured logging
}
```

## Testing

### Unit Testing

```go
import (
    "testing"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
)

func TestMCPClient(t *testing.T) {
    // Create mock transport for testing
    mockTransport := &transport.MockTransport{}
    client := client.NewMCPClient(mockTransport)
    
    // Test client operations
    ctx := context.Background()
    err := client.Connect(ctx)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
}
```

### Integration Testing

```go
func TestIntegration(t *testing.T) {
    // Start test MCP server
    server := startTestMCPServer(t)
    defer server.Stop()
    
    // Test client against real server
    client := client.NewMCPClient(transport.NewTCPTransport("localhost", server.Port()))
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    err := client.Connect(ctx)
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer client.Close()
    
    // Test operations
    tools, err := client.ListTools(ctx)
    if err != nil {
        t.Fatalf("Failed to list tools: %v", err)
    }
    
    if len(tools) == 0 {
        t.Error("Expected at least one tool")
    }
}
```

## API Reference

### Client Interface

```go
type MCPClient interface {
    Connect(ctx context.Context) error
    Close() error
    Initialize(ctx context.Context, clientName, version string) (*ServerInfo, error)
    ListTools(ctx context.Context) ([]Tool, error)
    CallTool(ctx context.Context, name string, args map[string]interface{}) (*ToolResult, error)
    ListResources(ctx context.Context) ([]Resource, error)
    ReadResource(ctx context.Context, uri string) ([]byte, error)
    ListPrompts(ctx context.Context) ([]Prompt, error)
    GetPrompt(ctx context.Context, name string, args map[string]interface{}) (*PromptResult, error)
}
```

### Transport Interface

```go
type Transport interface {
    Connect(ctx context.Context) error
    Send(ctx context.Context, message []byte) error
    Receive(ctx context.Context) ([]byte, error)
    Close() error
}
```

### Discovery Interface

```go
type Discoverer interface {
    DiscoverAll(ctx context.Context) ([]ServerInfo, error)
    DiscoverTCP(ctx context.Context, host string, startPort, endPort int) ([]ServerInfo, error)
    DiscoverHTTP(ctx context.Context, host string) ([]ServerInfo, error)
    DiscoverDocker(ctx context.Context) ([]ServerInfo, error)
}
```

## Examples Repository

For complete working examples, see the `examples/` directory:

- **Basic Examples:**
  - `examples/basic/simple-client/` - Minimal TCP client example
  - `examples/basic/transport-examples/` - Different transport types
  - `examples/basic/discovery/` - Server discovery functionality

- **Advanced Examples:**
  - `examples/advanced/builder-pattern/` - Advanced patterns and service wrappers

- **Complete Examples:**
  - `examples/complete/full-protocol/` - Full MCP protocol demonstration

- **CLI Examples:**
  - `examples/cmd/` - Command-line interface examples

Each example includes detailed README files with setup instructions and expected output.

## Troubleshooting

### Common Issues

1. **Connection Timeouts**
   ```go
   // Increase timeout
   ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
   ```

2. **Protocol Errors**
   ```go
   // Ensure server supports MCP protocol
   serverInfo, err := client.Initialize(ctx, "client", "1.0.0")
   if err != nil {
       // Check server compatibility
   }
   ```

3. **Transport Issues**
   ```go
   // Verify transport configuration
   transport := transport.NewTCPTransport("localhost", 8811)
   // Test connection separately if needed
   ```

### Debug Logging

Enable debug logging to troubleshoot issues:

```go
client := client.NewMCPClientBuilder().
    WithTransport(transport.NewTCPTransport("localhost", 8811)).
    WithDebugLogging(true).
    Build()
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to the library.

## License

MIT License - see [LICENSE](LICENSE) file for details.
