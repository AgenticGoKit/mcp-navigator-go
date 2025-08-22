# Simple Client Example

This example demonstrates the most basic usage of MCP Navigator as a library.

## What it does

- Creates a TCP connection to an MCP server
- Initializes the MCP protocol
- Lists available tools
- Executes the first tool with simple arguments

## Prerequisites

- MCP server running on `localhost:8811`

## Running the example

```bash
go run -tags example main.go
```

## Code walkthrough

1. **Create Transport**: `transport.NewTCPTransport("localhost", 8811)`
2. **Create Client**: `client.NewMCPClient(tcpTransport)`
3. **Connect**: `mcpClient.Connect(ctx)`
4. **Initialize**: `mcpClient.Initialize(ctx, "client-name", "version")`
5. **Use Tools**: `mcpClient.ListTools(ctx)` and `mcpClient.CallTool(ctx, name, args)`
6. **Cleanup**: `mcpClient.Close()`

## Expected output

```
🚀 Simple MCP Client Example
============================
📡 Connecting to MCP server...
✅ Connected successfully!
🚀 Initializing MCP protocol...
🖥️  Server: docker-mcp-server v0.0.1
📋 Listing available tools...
📝 Found 4 tools:
   1. docker - use the docker cli
   2. fetch - Fetches a URL from the internet...
   3. fetch_content - Fetch and parse content from a webpage URL
   4. search - Search DuckDuckGo and return formatted results
🔧 Calling tool: docker
✅ Tool executed successfully!
   Result contains 1 content item(s)
   Preview: Docker command executed successfully...
✅ Example completed successfully!
```

## Next steps

- Try [transport-examples](../transport-examples/) to see different connection types
- Explore [discovery](../discovery/) to automatically find servers
