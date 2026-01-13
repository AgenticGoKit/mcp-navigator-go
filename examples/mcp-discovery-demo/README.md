# MCP Tool Discovery Demo

This example demonstrates how to use MCP tool discovery in AgenticGoKit with practical, working implementations.

## What This Example Shows

✅ **MCP Tool Discovery** - How AgenticGoKit discovers tools from MCP servers  
✅ **Stdio Transport** - Running MCP servers as child processes  
✅ **TCP Transport** - Running MCP servers on network ports  
✅ **Automatic Discovery** - Auto-discovering MCP servers on the network  
✅ **Tool Execution** - Using discovered tools with the agent  

## Prerequisites

### Required
- Go 1.21+
- Ollama running locally with llama2 model available
  ```bash
  ollama pull llama2
  ```

### Optional (for advanced examples)
- Python (if you want to implement custom MCP servers)

## Quick Start (Stdio Mode - Recommended)

This is the easiest way to get started:

```bash
# Build everything
go build -o mcp-demo .
go build -o math-server math-server.go

# Run the example with stdio transport
./mcp-demo -mode stdio
```

**Expected Output:**
```
=== STDIO MCP Example ===

This example uses a simple math server via stdio transport.
...
✓ Found 4 tools:
  - add: Adds two numbers together
  - subtract: Subtracts second number from first number
  - multiply: Multiplies two numbers together
  - divide: Divides first number by second number

🤖 Using agent to perform calculation...
Input: "Calculate 42 plus 8 and tell me the result"

✓ Agent response:
The result of 42 plus 8 is 50.

=== Success! ===
```

## Running MCP Servers Outside Docker

### Option 1: Stdio Transport (Recommended for Testing)

**Pros:**
- No network exposure
- Simple to start/stop
- Perfect for testing
- No port conflicts

**Cons:**
- Only works with local processes
- One server per agent configuration

**How it works:**
The MCP server runs as a child process and communicates via stdin/stdout with JSON-RPC.

**Example Configuration:**
```go
MCPConfig{
    Servers: []MCPServer{
        {
            Name:    "math-tools",
            Type:    "stdio",
            Command: "./math-server",        // Path to executable
            Enabled: true,
        },
    },
}
```

**Starting a Stdio MCP Server:**
```bash
# Build the server
go build -o math-server math-server.go

# Agent automatically starts it when building
# (No manual start needed!)
```

### Option 2: TCP Transport (Network-Based)

**Pros:**
- Network accessible
- Multiple agents can use same server
- Can be on different machine
- Good for production

**Cons:**
- Requires manual server startup
- Network latency
- Port management

**Example Configuration:**
```go
MCPConfig{
    Servers: []MCPServer{
        {
            Name:    "tcp-math",
            Type:    "tcp",
            Address: "localhost",
            Port:    9999,
            Enabled: true,
        },
    },
}
```

**Starting a TCP MCP Server:**
```bash
# Terminal 1: Start the server
./math-server -mode tcp -port 9999

# Terminal 2: Run your agent
go build -o mcp-demo .
./mcp-demo -mode tcp
```

### Option 3: HTTP/SSE Transport (Web-Based)

**Pros:**
- Works over HTTP
- Standard web protocol
- Firewall friendly
- Good for cloud deployment

**Cons:**
- More complex setup
- Requires HTTP server

**Example Configuration:**
```go
MCPConfig{
    Servers: []MCPServer{
        {
            Name:    "web-tools",
            Type:    "http_sse",
            Address: "http://localhost:8080/mcp",
            Enabled: true,
        },
    },
}
```

**Starting an HTTP/SSE Server:**
```bash
# You would need to implement HTTP server
# For now, see the Python example below
```

### Option 4: Python MCP Servers (Without Docker)

Python has a rich MCP ecosystem. Here's how to use Python MCP servers without Docker:

**Install MCP SDK:**
```bash
pip install mcp
```

**Create a simple Python MCP server (`python-weather.py`):**
```python
#!/usr/bin/env python3
import json
import sys
from mcp.server import Server
from mcp.types import Tool

app = Server("weather-tools")

@app.call_tool
async def handle_tool(name: str, arguments: dict) -> str:
    if name == "get_temperature":
        location = arguments.get("location", "Unknown")
        # In real world, call weather API
        return f"The temperature in {location} is 72°F"
    return "Unknown tool"

@app.list_tools
async def list_tools():
    return [
        Tool(
            name="get_temperature",
            description="Get current temperature for a location",
            inputSchema={
                "type": "object",
                "properties": {
                    "location": {"type": "string", "description": "City name"}
                },
                "required": ["location"]
            }
        )
    ]

if __name__ == "__main__":
    app.run()
```

**Use in AgenticGoKit:**
```go
MCPConfig{
    Servers: []MCPServer{
        {
            Name:    "python-weather",
            Type:    "stdio",
            Command: "python3 python-weather.py",
            Enabled: true,
        },
    },
}
```

## Different Transport Types

| Transport | How It Works | Best For | Setup |
|-----------|-------------|----------|-------|
| **Stdio** | Child process, stdin/stdout | Testing, local tools | easiest |
| **TCP** | Network socket | Production, network | medium |
| **WebSocket** | WebSocket connection | Real-time, web apps | medium |
| **HTTP/SSE** | HTTP with Server-Sent Events | Cloud, REST APIs | complex |

## Example: Full Setup with Multiple Servers

**Directory structure:**
```
mcp-discovery-demo/
├── main.go              # Agent code
├── math-server.go       # Math tools
├── weather-server.go    # (Optional) Weather tools
└── README.md
```

**Configuration with multiple servers:**
```go
config := &vnext.Config{
    Tools: &vnext.ToolsConfig{
        Enabled: true,
        MCP: &vnext.MCPConfig{
            Enabled: true,
            Servers: []vnext.MCPServer{
                // Stdio server (child process)
                {
                    Name:    "math-tools",
                    Type:    "stdio",
                    Command: "./math-server",
                    Enabled: true,
                },
                // TCP server (running on port 9999)
                {
                    Name:    "weather-tools",
                    Type:    "tcp",
                    Address: "localhost",
                    Port:    9999,
                    Enabled: true,
                },
                // Remote HTTP server
                {
                    Name:    "remote-api",
                    Type:    "http_sse",
                    Address: "https://api.example.com/mcp",
                    Enabled: true,
                },
            },
        },
    },
}
```

## Testing Different Modes

### Test 1: Stdio Transport
```bash
# This will start math server automatically
./mcp-demo -mode stdio
```

### Test 2: TCP Transport
```bash
# Terminal 1: Start server manually
./math-server -mode tcp -port 9999

# Terminal 2: Run agent
./mcp-demo -mode tcp
```

### Test 3: Auto-Discovery
```bash
# Make sure at least one server is running on scanned ports
# Then run:
./mcp-demo -mode discover
```

## Running MCP Servers in Background

### Option A: Use `nohup` (Linux/macOS)
```bash
nohup ./math-server -mode tcp -port 9999 > math-server.log 2>&1 &
# Server continues running even after terminal closes
```

### Option B: Use PowerShell (Windows)
```powershell
Start-Process -FilePath "math-server.exe" -ArgumentList "-mode tcp -port 9999" -NoNewWindow
```

### Option C: Use `tmux` (Linux/macOS)
```bash
tmux new-session -d -s math-server './math-server -mode tcp -port 9999'
# Later, attach to it:
tmux attach-session -t math-server
```

### Option D: Use `screen` (Linux/macOS)
```bash
screen -d -m -S math-server ./math-server -mode tcp -port 9999
# Later, attach to it:
screen -r math-server
```

## Creating Custom MCP Servers

### Go Example
See `math-server.go` for a complete Go implementation.

### Python Example
```python
from mcp.server import Server

app = Server("my-tools")

@app.list_tools
async def list_tools():
    return [
        {
            "name": "my_tool",
            "description": "Does something useful",
            "inputSchema": {...}
        }
    ]

@app.call_tool
async def handle_tool(name: str, arguments: dict):
    if name == "my_tool":
        return "Result here"
    return "Unknown tool"

if __name__ == "__main__":
    app.run()
```

### Node.js Example
```javascript
const { Server } = require('@modelcontextprotocol/sdk/server');

const server = new Server('my-tools');

server.setRequestHandler(ListToolsRequest, async (request) => {
    return {
        tools: [
            {
                name: "my_tool",
                description: "Does something useful",
                inputSchema: {...}
            }
        ]
    };
});

server.setRequestHandler(CallToolRequest, async (request) => {
    if (request.params.name === "my_tool") {
        return {
            content: [{
                type: "text",
                text: "Result here"
            }]
        };
    }
    throw new Error("Unknown tool");
});

server.start();
```

## Troubleshooting

### "Tool not found" Error
1. Verify server is running
2. Check server output for errors
3. Ensure server name matches configuration
4. Try without -mode stdio (manual server)

### "Connection refused"
1. Check if TCP server is running: `lsof -i :9999`
2. Try different port if 9999 is in use
3. Verify firewall allows connection

### Server crashes silently
1. Run server manually to see errors: `./math-server -mode tcp -port 9999`
2. Check logs for parse errors
3. Verify JSON-RPC format is correct

### Agent doesn't find tools
1. Verify DebugMode = true to see detailed logs
2. Check tool registry: `vnext.DiscoverTools()`
3. Ensure MCP plugin is imported: `_ "github.com/agenticgokit/agenticgokit/plugins/mcp/stdio"`
4. Check agent config has Tools.MCP.Enabled = true

## Performance Notes

### Startup Time
- Stdio: 100-500ms (includes server startup)
- TCP: 10-100ms (already running)
- HTTP: 50-200ms (network latency)

### Optimization Tips
1. **Reuse agents** - Don't create new agents for every request
2. **Persistent servers** - Use TCP for long-running apps
3. **Caching** - Enable result caching for frequently used tools
4. **Parallel execution** - Run multiple tool calls concurrently

## What's Next?

1. ✅ Run the stdio example (fastest to get started)
2. ✅ Try the TCP example with manual server
3. ✅ Create your own MCP server (Python/Go/Node.js)
4. ✅ Add error handling and retry logic
5. ✅ Deploy to production (use TCP/HTTP)

## Additional Resources

- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [AgenticGoKit MCP Docs](../docs/MCP_DISCOVERY_QUICKREF.md)
- [Tool Discovery Guide](../docs/MCP_TOOL_DISCOVERY_GUIDE.md)
- [Troubleshooting Guide](../docs/MCP_DISCOVERY_TROUBLESHOOTING.md)

## Running Without Ollama

If you don't have Ollama installed, you can:

1. Use OpenAI (set LLM.Provider to "openai")
2. Use OpenRouter (set LLM.Provider to "openrouter")
3. Use a mock LLM (for testing):
   ```go
   LLM: vnext.LLMConfig{
       Provider: "mock",  // For testing only
       Model:    "test",
   }
   ```

## Getting Help

See the comprehensive documentation:
- `docs/MCP_DISCOVERY_QUICKREF.md` - Quick answers
- `docs/MCP_TOOL_DISCOVERY_GUIDE.md` - Complete guide
- `docs/MCP_DISCOVERY_TROUBLESHOOTING.md` - Debugging
