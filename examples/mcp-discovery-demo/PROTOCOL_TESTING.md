# Protocol Testing Guide

This directory contains test clients and servers that demonstrate the mcp-navigator-go library working with different transport protocols.

## Supported Protocols

The math-server and test clients support three transport protocols:

1. **TCP** - Direct TCP socket connection
2. **STDIO** - Standard input/output (subprocess communication)
3. **WebSocket** - WebSocket over HTTP

## Building

First, build the math-server and test clients:

```powershell
# Build math-server
go build -o math-server.exe math-server.go

# Build test clients
go build -o test-client-tcp.go
go build -o test-client-stdio.go
go build -o test-client-websocket.go
```

Or compile all at once:

```powershell
go build -o math-server.exe math-server.go;
go build -o test-client-tcp.exe test-client-tcp.go;
go build -o test-client-stdio.exe test-client-stdio.go;
go build -o test-client-websocket.exe test-client-websocket.go
```

## Testing Each Protocol

### 1. TCP Protocol Test

Terminal 1 - Start the math-server in TCP mode:
```powershell
.\math-server.exe -mode tcp -port 9999
```

Expected output:
```
Math MCP Server started (TCP mode on :9999)
```

Terminal 2 - Run the TCP test client:
```powershell
.\test-client-tcp.exe
```

or with debug logging:
```powershell
.\test-client-tcp.exe -debug
```

Expected output:
```
=== MCP Navigator - TCP Protocol Test ===
Testing math-server on TCP port 9999
📡 Connecting via TCP...
✅ Connected successfully

🔧 Initializing MCP protocol...
✅ Initialized with server: math-server v1.0.0

🔍 Listing tools...

📊 Results: Found 4 tools
============================================================
🔧 Tool #1: add
   Description: Adds two numbers together
   InputSchema: {
     "properties": {...}
   }

   🧪 Testing add tool...
   ✅ Result: {Contents:[{Type:text Text:Result: 5 + 3 = 8}]}

============================================================
✅ TCP Test Completed Successfully
```

### 2. STDIO Protocol Test

Terminal 1 - Start the test client (it launches the server automatically):
```powershell
.\test-client-stdio.exe
```

or with debug logging:
```powershell
.\test-client-stdio.exe -debug
```

Expected output:
```
=== MCP Navigator - STDIO Protocol Test ===
Testing math-server via STDIO transport
📡 Connecting via STDIO...
✅ Connected successfully

🔧 Initializing MCP protocol...
✅ Initialized with server: math-server v1.0.0

🔍 Listing tools...

📊 Results: Found 4 tools
============================================================
🔧 Tool #1: add
   Description: Adds two numbers together
   ...

   🧪 Testing multiply tool...
   ✅ Result: {Contents:[{Type:text Text:Result: 4 * 6 = 24}]}

============================================================
✅ STDIO Test Completed Successfully
```

### 3. WebSocket Protocol Test

Terminal 1 - Start the math-server in WebSocket mode:
```powershell
.\math-server.exe -mode websocket -port 9999
```

Expected output:
```
Math MCP Server started (WebSocket mode on ws://localhost:9999/mcp)
```

Terminal 2 - Run the WebSocket test client:
```powershell
.\test-client-websocket.exe
```

or with debug logging:
```powershell
.\test-client-websocket.exe -debug
```

Expected output:
```
=== MCP Navigator - WebSocket Protocol Test ===
Testing math-server on WebSocket port 9999
📡 Connecting via WebSocket...
✅ Connected successfully

🔧 Initializing MCP protocol...
✅ Initialized with server: math-server v1.0.0

🔍 Listing tools...

📊 Results: Found 4 tools
============================================================
🔧 Tool #1: add
   Description: Adds two numbers together
   ...

   🧪 Testing divide tool...
   ✅ Result: {Contents:[{Type:text Text:Result: 20 / 4 = 5}]}

============================================================
✅ WebSocket Test Completed Successfully
```

## Batch Testing Script

Run all protocol tests with a single PowerShell script:

```powershell
# Test all protocols sequentially
.\run-all-protocol-tests.ps1
```

This will:
1. Build all executables
2. Start math-server in TCP mode and run TCP test
3. Run STDIO test (launches its own server)
4. Start math-server in WebSocket mode and run WebSocket test
5. Display a summary of results

## Debug Mode

All test clients support a `-debug` flag for detailed logging:

```powershell
# See detailed protocol handshake and message details
.\test-client-tcp.exe -debug

# Also works with other clients
.\test-client-stdio.exe -debug
.\test-client-websocket.exe -debug
```

Debug output includes:
- Component-based logging with `[CLIENT]` and `[TRANSPORT]` prefixes
- Detailed request/response details
- Message parsing information
- Connection state changes

## Protocol Comparison

| Feature | TCP | STDIO | WebSocket |
|---------|-----|-------|-----------|
| Connection Setup | Direct socket | Subprocess launch | HTTP upgrade |
| Performance | Lowest latency | Good (process overhead) | Good (HTTP overhead) |
| Firewall Friendly | Port required | N/A (local) | Port required |
| Scalability | Many connections | One process | Many connections |
| Use Case | Local server | Local tool integration | Remote server |

## Tools Available

The math-server provides 4 MCP tools:

1. **add** - Adds two numbers
   - Parameters: a (number), b (number)
   - Returns: Result message

2. **subtract** - Subtracts two numbers
   - Parameters: a (number), b (number)
   - Returns: Result message

3. **multiply** - Multiplies two numbers
   - Parameters: a (number), b (number)
   - Returns: Result message

4. **divide** - Divides two numbers
   - Parameters: a (number), b (number)
   - Returns: Result message (or error if b=0)

## Troubleshooting

### Port Already in Use
If you get "address already in use" error:
```powershell
# Find process using port 9999
Get-NetTCPConnection -LocalPort 9999 | Stop-Process -Force
```

### STDIO Test Fails
Make sure `math-server.exe` is in the same directory as the test client, or provide the full path:
```go
trans := transport.NewStdioTransport("C:\\path\\to\\math-server.exe", []string{"-mode", "stdio"})
```

### WebSocket Connection Refused
Make sure the math-server is running in WebSocket mode:
```powershell
.\math-server.exe -mode websocket -port 9999
```

## Expected Results

After running all tests, you should see:
- ✅ All 4 tools returned in each protocol
- ✅ Tool schemas validated (MCP spec compliant)
- ✅ Tool execution works (add, subtract, multiply, divide)
- ✅ Protocol handshakes complete successfully
- ✅ No errors in any output

This confirms that the mcp-navigator-go library correctly handles:
- Multiple transport protocols
- Protocol initialization
- Tool discovery (ListTools)
- Tool execution (CallTool)
- Full MCP specification compliance
- Proper error handling

## Next Steps

To integrate these tests into your CI/CD pipeline or develop your own MCP server:

1. Study the [test-client-tcp.go](test-client-tcp.go) for client implementation examples
2. Review [math-server.go](math-server.go) to see how to implement MCP servers
3. Adapt the patterns for your own use cases
4. Use the `-debug` flag to troubleshoot any integration issues

