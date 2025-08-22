# MCP Navigator Go - API Reference

This document provides a concise API reference for using MCP Navigator as a library in your Go applications.

## Quick Start

```go
import (
    "context"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
    "github.com/kunalkushwaha/mcp-navigator-go/pkg/discovery"
)
```

## Core Types

### Client

The library exposes a concrete `Client` type in package `client`. The current constructor and core methods are:

```go
// Create new client
func NewClient(transport transport.Transport, config client.ClientConfig) *client.Client

// ClientConfig (passed by value)
type ClientConfig struct {
    Name    string
    Version string
    Logger  *log.Logger
    Timeout time.Duration
}

// Connection management
func (c *client.Client) Connect(ctx context.Context) error
func (c *client.Client) Disconnect() error

// Protocol initialization
// Initialize identifies the client to the server and performs the MCP handshake.
func (c *client.Client) Initialize(ctx context.Context, clientInfo mcp.ClientInfo) error

// After Initialize, use GetServerInfo() to inspect server name/version
func (c *client.Client) GetServerInfo() *mcp.ServerInfo
func (c *client.Client) GetServerCapabilities() *mcp.ServerCapabilities

// Tools
func (c *client.Client) ListTools(ctx context.Context) ([]mcp.Tool, error)
func (c *client.Client) CallTool(ctx context.Context, name string, args map[string]interface{}) (*mcp.CallToolResponse, error)

// Resources
func (c *client.Client) ListResources(ctx context.Context) ([]mcp.Resource, error)
func (c *client.Client) ReadResource(ctx context.Context, uri string) (*mcp.ReadResourceResponse, error)

// Prompts
func (c *client.Client) ListPrompts(ctx context.Context) ([]mcp.Prompt, error)
func (c *client.Client) GetPrompt(ctx context.Context, name string, args map[string]interface{}) (*mcp.GetPromptResponse, error)
```

### Transport Types

```go
// TCP Transport
func NewTCPTransport(host string, port int) *transport.TCPTransport

// HTTP (generic) Transport factory
func NewHTTPTransport(baseURL, endpoint string) *transport.HTTPTransport

// HTTP SSE Transport (real-time)
func NewSSETransport(baseURL, endpoint string) *transport.SSETransport

// HTTP Streaming Transport (request-response)
func NewStreamingHTTPTransport(baseURL, endpoint string) *transport.StreamingHTTPTransport

// STDIO Transport (process-based)
// Note: factory name is `NewStdioTransport` (capitalization: Stdio)
func NewStdioTransport(command string, args []string) *transport.StdioTransport

// WebSocket Transport
func NewWebSocketTransport(wsURL string) *transport.WebSocketTransport
```

### Discovery

```go
type Discoverer struct{}

func NewDiscoverer() *Discoverer

// Discover all server types
func (d *Discoverer) DiscoverAll(ctx context.Context) ([]ServerInfo, error)

// Discover specific types
func (d *Discoverer) DiscoverTCP(ctx context.Context, host string, startPort, endPort int) ([]ServerInfo, error)
func (d *Discoverer) DiscoverHTTP(ctx context.Context, host string) ([]ServerInfo, error)
func (d *Discoverer) DiscoverDocker(ctx context.Context) ([]ServerInfo, error)
```

## Data Types

### ServerInfo

```go
type ServerInfo struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    Type        string `json:"type"`        // "tcp", "http", "docker"
    Address     string `json:"address"`
    Description string `json:"description"`
}
```

### Tool

```go
type Tool struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    InputSchema map[string]interface{} `json:"inputSchema"`
}
```

### ToolResult

```go
type Tool struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    InputSchema map[string]interface{} `json:"inputSchema"`
}

// CallToolResponse is returned by client.CallTool
type CallToolResponse struct {
    Content []mcp.Content `json:"content"`
    IsError bool          `json:"isError,omitempty"`
}

// Content (from package mcp) - fields are plain strings (not pointers)
type Content struct {
    Type        string                 `json:"type"`
    Text        string                 `json:"text,omitempty"`
    Data        string                 `json:"data,omitempty"`
    MimeType    string                 `json:"mimeType,omitempty"`
    Name        string                 `json:"name,omitempty"`
    URI         string                 `json:"uri,omitempty"`
    Annotations map[string]interface{} `json:"annotations,omitempty"`
}
```
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    MimeType    string `json:"mimeType,omitempty"`
}
```go
type Resource struct {
    URI         string                 `json:"uri"`
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    MimeType    string                 `json:"mimeType,omitempty"`
    Annotations map[string]interface{} `json:"annotations,omitempty"`
}
```
    Arguments   []PromptArgument       `json:"arguments,omitempty"`
}

type PromptArgument struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    Required    bool   `json:"required,omitempty"`
}
```

## Error Types

```text
Error helpers: The repository does not export generic helper functions like `IsConnectionError` or `IsProtocolError` by default. The client returns wrapped errors; users should check error strings or implement their own helpers (or use `errors.Is`/`errors.As` if sentinel errors are introduced). The client does use internal sentinel errors (for example: ErrNotConnected) in places — consult the code for specific error values.
```

## Usage Patterns

### Basic Connection

```go
// Create transport and client
transport := transport.NewTCPTransport("localhost", 8811)
config := client.ClientConfig{Name: "my-app", Version: "1.0.0"}
client := client.NewClient(transport, config)

// Connect with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := client.Connect(ctx)
if err != nil {
    return err
}
defer client.Close()

// Initialize protocol
serverInfo, err := client.Initialize(ctx, "my-app", "1.0.0")
if err != nil {
    return err
}
```

### Server Discovery

```go
discoverer := discovery.NewDiscoverer()

// Find all servers
servers, err := discoverer.DiscoverAll(ctx)
if err != nil {
    return err
}

for _, server := range servers {
    fmt.Printf("Found: %s (%s) at %s\n", server.Name, server.Type, server.Address)
}
```

### Tool Execution

```go
// List available tools
tools, err := client.ListTools(ctx)
if err != nil {
    return err
}

// Execute a tool
if len(tools) > 0 {
    args := map[string]interface{}{
        "query": "search term",
        "limit": 10,
    }

    result, err := client.CallTool(ctx, tools[0].Name, args)
    if err != nil {
        return err
    }

    // Process result (note Content.Text is a string)
    for _, content := range result.Content {
        if content.Text != "" {
            fmt.Println(content.Text)
        }
    }
}
```

### HTTP Transports

```go
// SSE (Server-Sent Events) for real-time communication
// SSE (Server-Sent Events) for real-time communication
sseTransport := transport.NewSSETransport("http://localhost:8812", "/sse/")
config := client.ClientConfig{Name: "sse-client", Version: "1.0.0"}
sseClient := client.NewClient(sseTransport, config)

// Streaming for simple request-response
streamingTransport := transport.NewStreamingHTTPTransport("http://localhost:8813", "/mcp")
streamingClient := client.NewClient(streamingTransport, config)
```

### Error Handling

```go
err := client.Connect(ctx)
if err != nil {
    switch {
    case client.IsConnectionError(err):
        // Handle connection issues
        log.Printf("Connection failed: %v", err)
    case client.IsTimeoutError(err):
        // Handle timeouts
        log.Printf("Operation timed out: %v", err)
    case client.IsProtocolError(err):
        // Handle protocol issues
        log.Printf("Protocol error: %v", err)
    default:
        log.Printf("Unknown error: %v", err)
    }
}
```

### Resource Management

```go
// List resources
resources, err := client.ListResources(ctx)
if err != nil {
    return err
}

// Read resource content
for _, resource := range resources {
    content, err := client.ReadResource(ctx, resource.URI)
    if err != nil {
        log.Printf("Failed to read %s: %v", resource.URI, err)
        continue
    }
    
    // Process content
    fmt.Printf("Resource %s: %s\n", resource.URI, string(content))
}
```

### Concurrent Operations

```go
import "sync"

var wg sync.WaitGroup
results := make(chan ToolResult, len(tools))

for _, tool := range tools {
    wg.Add(1)
    go func(toolName string) {
        defer wg.Done()
        
        result, err := client.CallTool(ctx, toolName, map[string]interface{}{})
        if err != nil {
            log.Printf("Tool %s failed: %v", toolName, err)
            return
        }
        
        results <- *result
    }(tool.Name)
}

wg.Wait()
close(results)

for result := range results {
    // Process results
}
```

## Configuration Options

### Environment Variables

```go
import "os"

host := os.Getenv("MCP_CLIENT_HOST")
if host == "" {
    host = "localhost"
}

port := os.Getenv("MCP_CLIENT_PORT")
if port == "" {
    port = "8811"
}
```

### Timeouts

```go
// Connection timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Operation timeout
opCtx, opCancel := context.WithTimeout(context.Background(), 10*time.Second)
defer opCancel()
```

## Testing Support

### Mock Transport

```go
type MockTransport struct {
    // Mock implementation
}

func (m *MockTransport) Connect(ctx context.Context) error { return nil }
func (m *MockTransport) Send(ctx context.Context, data []byte) error { return nil }
func (m *MockTransport) Receive(ctx context.Context) ([]byte, error) { return nil, nil }
func (m *MockTransport) Close() error { return nil }

// Use in tests
mockTransport := &MockTransport{}
client := client.NewMCPClient(mockTransport)
```

### Test Helpers

```go
// Create test client
func NewTestClient(t *testing.T) *client.MCPClient {
    transport := &MockTransport{}
    return client.NewMCPClient(transport)
}

// Test with timeout
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    client := NewTestClient(t)
    err := client.Connect(ctx)
    require.NoError(t, err)
}
```

## Performance Considerations

### Connection Pooling

```go
// Reuse connections when possible
type ConnectionPool struct {
    clients map[string]*client.MCPClient
    mutex   sync.RWMutex
}

func (p *ConnectionPool) GetClient(address string) (*client.MCPClient, error) {
    p.mutex.RLock()
    client, exists := p.clients[address]
    p.mutex.RUnlock()
    
    if exists {
        return client, nil
    }
    
    // Create new client
    // ... implementation
}
```

### Batch Operations

```go
// Batch tool calls
type BatchRequest struct {
    Tools []ToolCall `json:"tools"`
}

type ToolCall struct {
    Name string                 `json:"name"`
    Args map[string]interface{} `json:"arguments"`
}
```

## Best Practices

1. **Always use context with timeout** for operations
2. **Close clients** when done to free resources
3. **Handle errors gracefully** with appropriate error types
4. **Reuse clients** for multiple operations when possible
5. **Use discovery** to find available servers dynamically
6. **Test with mock transports** for unit testing
7. **Configure timeouts** based on your use case
8. **Log errors** with structured logging for debugging

For complete examples, see:
- [Basic Examples](examples/basic/) - Simple client, transports, discovery
- [Advanced Examples](examples/advanced/) - Production patterns, service wrappers
- [Complete Examples](examples/complete/) - Full protocol demonstration
