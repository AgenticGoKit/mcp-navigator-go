# Complete Examples

These examples demonstrate comprehensive usage of the entire MCP protocol specification, including Tools, Resources, and Prompts.

## 📁 Examples

### [full-protocol](full-protocol/)
**Complete MCP protocol implementation**
- **Tools:** List, execute with various argument patterns
- **Resources:** Discover, read different resource types  
- **Prompts:** List, execute with dynamic arguments
- **Advanced Features:** Error handling, concurrent operations, timeout management

**When to use:** Applications requiring full MCP protocol support

## 🎯 Protocol Coverage

| Feature | Coverage | Example Usage |
|---------|----------|---------------|
| **Tools** | ✅ Complete | List tools, execute with args, handle responses |
| **Resources** | ✅ Complete | Discover resources, read content, handle MIME types |
| **Prompts** | ✅ Complete | List prompts, handle arguments, process messages |
| **Error Handling** | ✅ Advanced | Timeout, retry, graceful failure |
| **Concurrency** | ✅ Advanced | Parallel operations, safe execution |

## 🚀 Quick Start

```bash
cd full-protocol
go run -tags example main.go
```

## 📋 Expected Output

```
🚀 MCP Client - Complete Protocol Features
==========================================
📡 Connecting to MCP server...
🚀 Initializing MCP protocol...
🖥️  Connected to: docker-mcp-server v0.0.1

🛠️  1. Tools Support
===================
📋 Listing available tools...
📝 Found 4 tools:
   1. docker
      Description: use the docker cli
   2. fetch  
      Description: Fetches a URL from the internet...
   3. fetch_content
      Description: Fetch and parse content from a webpage URL
   4. search
      Description: Search DuckDuckGo and return formatted results

🔧 Executing tool 'docker'...
✅ Tool execution successful!
   Content items: 1
   Is error: false
   Content 1: Type=text
   Preview: Docker command executed successfully...

📂 2. Resources Support
======================
📋 Listing available resources...
📁 Found 0 resources:
ℹ️  No resources available on this server

💬 3. Prompts Support  
====================
📋 Listing available prompts...
📝 Found 0 prompts:
ℹ️  No prompts available on this server

⚡ 4. Advanced Usage Patterns
============================
🔍 Testing error handling...
✅ Error handling working: tool not found: non-existent-tool

🚀 Testing concurrent tool execution...
   Concurrent call 0: ✅ 1 items
   Concurrent call 1: ✅ 1 items  
   Concurrent call 2: ✅ 1 items

⏱️  Testing timeout handling...
✅ Timeout handling working: context deadline exceeded

✅ Complete protocol demonstration finished!
```

## 🔍 What This Example Demonstrates

### Tool Management
- **Discovery:** Enumerate all available tools
- **Execution:** Call tools with various argument patterns
- **Response Handling:** Process different content types
- **Error Management:** Handle tool execution failures

### Resource Access
- **Listing:** Discover available resources
- **Reading:** Fetch resource content
- **MIME Handling:** Process different resource types
- **URI Management:** Handle resource addressing

### Prompt Processing
- **Discovery:** Find available prompts
- **Argument Handling:** Process required/optional parameters
- **Message Processing:** Handle prompt responses
- **Dynamic Arguments:** Build arguments programmatically

### Advanced Patterns
- **Concurrent Execution:** Safe parallel operations
- **Timeout Management:** Context-based timeouts
- **Error Recovery:** Graceful failure handling
- **Resource Cleanup:** Proper connection management

## 📋 Prerequisites

- Completed [Basic Examples](../basic/)
- Understanding of [Advanced Patterns](../advanced/)
- MCP server with full protocol support

## 🔧 Customization

Modify the example for your use case:

```go
// Change connection parameters
tcpTransport := transport.NewTCPTransport("your-host", your-port)

// Adjust timeouts
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

// Customize tool arguments
arguments := map[string]interface{}{
    "your_param": "your_value",
    "limit": 50,
}
```

## 📖 Integration Guide

### 1. Start Simple
- Begin with basic tool execution
- Add resource access as needed
- Implement prompts for AI workflows

### 2. Add Error Handling
- Implement proper timeout management
- Add retry logic for failed operations
- Handle partial failures gracefully

### 3. Optimize Performance
- Use concurrent operations where appropriate
- Implement connection pooling for high volume
- Cache frequently accessed resources

### 4. Production Readiness
- Add comprehensive logging
- Implement health checks
- Monitor performance metrics
- Plan for disaster recovery

## 📚 Further Reading

- [MCP Specification](https://spec.modelcontextprotocol.io/) - Official protocol documentation
- [Library Documentation](../../LIBRARY.md) - Comprehensive integration guide
- [API Reference](../../API.md) - Complete API documentation
