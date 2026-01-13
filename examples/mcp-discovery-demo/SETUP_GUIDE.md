# MCP Tools Outside Docker - Complete Setup Guide

You now have everything needed to run MCP tools without Docker Desktop. Here's what's been created and how to use it.

## 📦 What You Got

### 1. Working Example Project
**Location:** `examples/mcp-discovery-demo/`

- **main.go** - Example agent with 3 modes (stdio, tcp, discover)
- **math-server.go** - Simple MCP server providing math tools (add, subtract, multiply, divide)
- **README.md** - Comprehensive documentation
- **quick-start.sh** - Bash script for Linux/macOS
- **quick-start.ps1** - PowerShell script for Windows
- **go.mod** - Module file

### 2. Documentation
- **RUNNING_MCP_SERVERS_OUTSIDE_DOCKER.md** - Complete guide with:
  - How each transport works (stdio, TCP, WebSocket, HTTP/SSE)
  - Step-by-step setup for each transport
  - Language-specific implementations (Go, Python, Node.js, Rust)
  - Process manager integration (systemd, supervisor, launchd)
  - Troubleshooting and monitoring

---

## 🚀 Quick Start (5 minutes)

### Option A: Windows (PowerShell)

```powershell
# Navigate to demo directory
cd examples\mcp-discovery-demo

# Build everything
.\quick-start.ps1 build

# Run with stdio (easiest)
.\quick-start.ps1 run-stdio
```

### Option B: Linux/macOS (Bash)

```bash
# Navigate to demo directory
cd examples/mcp-discovery-demo

# Make script executable
chmod +x quick-start.sh

# Build everything
./quick-start.sh build

# Run with stdio (easiest)
./quick-start.sh run-stdio
```

---

## 🎯 Three Transport Options

### 1️⃣ Stdio Transport (Recommended for Testing)

**Best for:** Development, testing, learning

**How it works:**
- Agent automatically starts math-server as child process
- Communication via stdin/stdout (JSON-RPC)
- No manual server management needed

**Run it:**
```bash
# Windows PowerShell
.\quick-start.ps1 run-stdio

# Linux/macOS
./quick-start.sh run-stdio
```

**Advantages:**
✅ No port conflicts  
✅ Automatic startup/shutdown  
✅ Perfect for learning  
✅ See all output directly  

**Disadvantages:**
❌ Local only (no remote servers)  
❌ Process overhead  

**Expected output:**
```
=== STDIO MCP Example ===
...
✓ Found 4 tools:
  - add: Adds two numbers together
  - subtract: Subtracts second number from first number
  - multiply: Multiplies two numbers together
  - divide: Divides first number by second number

✓ Agent response:
The result of 42 plus 8 is 50.

=== Success! ===
```

---

### 2️⃣ TCP Transport (For Production)

**Best for:** Production, multiple agents, remote servers

**How it works:**
- You start math-server on port 9999 (in one terminal)
- Agent connects to it over network socket
- Multiple agents can use same server

**Step 1: Start the server**
```bash
# Windows PowerShell
.\quick-start.ps1 start-tcp-server

# Linux/macOS
./quick-start.sh start-tcp-server

# Output: "Math MCP Server started (TCP mode on :9999)"
```

**Step 2: Run the agent (in another terminal)**
```bash
# Windows PowerShell
.\quick-start.ps1 run-tcp

# Linux/macOS
./quick-start.sh run-tcp
```

**Advantages:**
✅ Multiple agents can connect to same server  
✅ Can be remote (different machine)  
✅ Better performance (no startup overhead)  
✅ Production-ready  

**Disadvantages:**
❌ Must manage server lifecycle  
❌ Port conflicts possible  
❌ Network latency  

---

### 3️⃣ Auto-Discovery Mode

**Best for:** Finding MCP servers on network automatically

**How it works:**
- Agent scans predefined ports (8080, 8081, 8090, 9999, 3000, 3001)
- Automatically connects to any MCP servers it finds
- Discovers tools from all servers

**Run it:**
```bash
# Windows PowerShell
.\quick-start.ps1 run-discover

# Linux/macOS
./quick-start.sh run-discover
```

**When to use:**
- Multiple servers running on network
- Don't want to hardcode server addresses
- Dynamic infrastructure

---

## 🛠️ Manual Setup (Without Quick-Start Script)

### Build the binaries

```bash
cd examples/mcp-discovery-demo

# Build agent
go build -o mcp-demo main.go

# Build math server
go build -o math-server math-server.go
```

### Run with different transports

```bash
# Stdio (automatic server startup)
./mcp-demo -mode stdio

# TCP (manual server startup)
# Terminal 1:
./math-server -mode tcp -port 9999

# Terminal 2:
./mcp-demo -mode tcp

# Auto-discovery
./mcp-demo -mode discover
```

---

## 📚 Creating Your Own MCP Server

### Go Example

Create `my-server.go`:

```go
package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
)

type Request struct {
    Method string      `json:"method"`
    Params interface{} `json:"params"`
    ID     interface{} `json:"id"`
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        var req Request
        json.Unmarshal(scanner.Bytes(), &req)
        
        response := handleRequest(&req)
        data, _ := json.Marshal(response)
        fmt.Println(string(data))
    }
}

func handleRequest(req *Request) map[string]interface{} {
    if req.Method == "tools/list" {
        return map[string]interface{}{
            "jsonrpc": "2.0",
            "result": map[string]interface{}{
                "tools": []map[string]interface{}{
                    {
                        "name":        "my_tool",
                        "description": "My custom tool",
                    },
                },
            },
            "id": req.ID,
        }
    }
    return nil
}
```

### Python Example

Create `my-server.py`:

```python
#!/usr/bin/env python3
import json
import sys

def handle_request(req):
    if req['method'] == 'tools/list':
        return {
            "jsonrpc": "2.0",
            "result": {
                "tools": [
                    {
                        "name": "my_tool",
                        "description": "My custom tool",
                    }
                ]
            },
            "id": req['id']
        }

if __name__ == "__main__":
    for line in sys.stdin:
        if line.strip():
            req = json.loads(line)
            response = handle_request(req)
            print(json.dumps(response))
```

### Use in agent

```go
config := &vnext.Config{
    Tools: &vnext.ToolsConfig{
        Enabled: true,
        MCP: &vnext.MCPConfig{
            Enabled: true,
            Servers: []vnext.MCPServer{
                {
                    Name:    "my-tools",
                    Type:    "stdio",
                    Command: "./my-server",  // or "python3 my-server.py"
                    Enabled: true,
                },
            },
        },
    },
}

agent, _ := vnext.NewBuilder("agent").
    WithConfig(config).
    Build()
```

---

## 🔧 Running Servers in Background

### Linux/macOS

```bash
# Option 1: nohup (simple)
nohup ./math-server -mode tcp -port 9999 > server.log 2>&1 &

# Check it's running
ps aux | grep math-server

# View logs
tail -f server.log

# Stop it
pkill -f "math-server"
```

```bash
# Option 2: screen (better)
screen -d -m -S math-server \
  ./math-server -mode tcp -port 9999

# Later, attach
screen -r math-server

# Detach: Ctrl-A then D
```

```bash
# Option 3: tmux (modern)
tmux new-session -d -s math-server \
  './math-server -mode tcp -port 9999'

# Later, attach
tmux attach-session -t math-server

# Detach: Ctrl-B then D
```

### Windows (PowerShell)

```powershell
# Start in background
Start-Process -FilePath "math-server.exe" `
              -ArgumentList "-mode tcp -port 9999" `
              -NoNewWindow `
              -RedirectStandardOutput "server.log"

# Check running processes
Get-Process | Where-Object { $_.Name -match "math" }

# Stop it
Stop-Process -Name "math-server"
```

### As System Service (Linux - systemd)

```bash
# Create service file
sudo tee /etc/systemd/system/mcp-math.service > /dev/null <<EOF
[Unit]
Description=MCP Math Server
After=network.target

[Service]
Type=simple
ExecStart=/home/user/math-server -mode tcp -port 9999
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable mcp-math
sudo systemctl start mcp-math

# Check status
sudo systemctl status mcp-math

# View logs
sudo journalctl -u mcp-math -f
```

---

## 🧪 Testing Without Agent

### Test Stdio Server

```bash
# Send list tools request
echo '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}' | ./math-server

# Expected: JSON response with tools list
```

### Test TCP Server

```bash
# Terminal 1: Start server
./math-server -mode tcp -port 9999

# Terminal 2: Test connection
telnet localhost 9999

# Type request and press Enter twice:
{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}
```

### Using Quick-Start Script

```bash
# Windows PowerShell
.\quick-start.ps1 test-stdio
.\quick-start.ps1 test-tcp

# Linux/macOS
./quick-start.sh test-stdio
./quick-start.sh test-tcp
```

---

## ⚠️ Common Issues & Fixes

### Issue: "Tool not found" error

**Cause:** MCP plugin not imported

**Fix:** Make sure to import the stdio plugin:
```go
import _ "github.com/agenticgokit/agenticgokit/plugins/mcp/stdio"
```

### Issue: "Connection refused" on TCP

**Cause:** Server not running

**Fix:** Check if server is running:
```bash
# Linux/macOS
netstat -an | grep 9999

# Windows PowerShell
netstat -ano | findstr 9999
```

Start it:
```bash
./math-server -mode tcp -port 9999
```

### Issue: "Port already in use"

**Cause:** Another process using the port

**Fix:** Either stop that process or use different port:
```bash
# Use port 9998 instead
./math-server -mode tcp -port 9998

# Update agent config to match
MCPConfig{
    Servers: []MCPServer{
        {Port: 9998, ...}
    }
}
```

### Issue: Ollama not running

**Cause:** Ollama service not started

**Fix:**
```bash
# Linux/macOS
ollama serve

# Windows: Run Ollama.exe or use OllamA's system tray app

# Pull model if needed
ollama pull llama2
```

### Issue: Agent completes instantly with no output

**Cause:** Missing LLM provider configuration or no tools discovered

**Fix:** Enable debug mode to see what's happening:
```go
config.DebugMode = true  // Shows detailed logs
```

---

## 📊 Performance Comparison

| Metric | Stdio | TCP | HTTP/SSE |
|--------|-------|-----|----------|
| Startup Time | 100-500ms | 10-100ms | 50-200ms |
| Memory (server) | ~5-20MB | ~10-50MB | ~5-15MB |
| Latency | Low | Very low | Medium |
| Network | No | Yes | Yes |
| Multiple agents | No | Yes | Yes |
| Best for | Testing | Production | Cloud |

---

## 📖 Full Documentation

For more detailed information, see:

1. **Quick Reference:** `docs/MCP_DISCOVERY_QUICKREF.md`
2. **Complete Guide:** `docs/MCP_TOOL_DISCOVERY_GUIDE.md`
3. **Running Servers:** `docs/RUNNING_MCP_SERVERS_OUTSIDE_DOCKER.md`
4. **Troubleshooting:** `docs/MCP_DISCOVERY_TROUBLESHOOTING.md`
5. **Code Flow:** `docs/MCP_DISCOVERY_CODE_FLOW.md`
6. **Example Project:** `examples/mcp-discovery-demo/README.md`

---

## ✅ Next Steps

### 1. Try the Stdio Example (Fastest)
```bash
cd examples/mcp-discovery-demo
./quick-start.sh run-stdio    # or quick-start.ps1 run-stdio (Windows)
```

### 2. Try TCP with Manual Server
```bash
# Terminal 1
./quick-start.sh start-tcp-server

# Terminal 2
./quick-start.sh run-tcp
```

### 3. Create Your Own MCP Server
- Copy math-server.go as template
- Add your own tools
- Use with your agent

### 4. Integrate with Your Agent
- Add MCP server to agent config
- Import stdio plugin
- Call agent.Build()

### 5. Deploy to Production
- Use TCP transport for network isolation
- Run as system service (systemd/launchd)
- Enable logging and monitoring

---

## 🎓 Learning Path

**Phase 1: Understanding (30 min)**
- [ ] Read this file completely
- [ ] Look at math-server.go to understand JSON-RPC
- [ ] Review example in main.go

**Phase 2: Hands-On (1 hour)**
- [ ] Run stdio example
- [ ] Test with manual TCP server
- [ ] Run discovery mode

**Phase 3: Creation (2+ hours)**
- [ ] Create your own MCP server
- [ ] Integrate with your agent
- [ ] Test all transport modes

**Phase 4: Production (ongoing)**
- [ ] Set up as system service
- [ ] Enable monitoring/logging
- [ ] Handle error cases
- [ ] Optimize performance

---

## 💡 Pro Tips

1. **Start with Stdio** - Simplest way to test
2. **Use TCP for Production** - Better performance
3. **Enable DebugMode** - See what's happening
4. **Cache Results** - Tools can be slow, cache helps
5. **Run Multiple Servers** - Combine different tools
6. **Monitor Logs** - Find issues faster

---

## 🆘 Getting Help

### Check Documentation
1. Example README - `examples/mcp-discovery-demo/README.md`
2. Running Servers - `docs/RUNNING_MCP_SERVERS_OUTSIDE_DOCKER.md`
3. Troubleshooting - `docs/MCP_DISCOVERY_TROUBLESHOOTING.md`
4. Code Flow - `docs/MCP_DISCOVERY_CODE_FLOW.md`

### Debug Steps
1. Enable DebugMode in config
2. Check logs for errors
3. Verify server is responding
4. Test directly with telnet/netcat
5. Check if port is in use
6. Verify Ollama is running

---

## 📝 Summary

You now have:

✅ Working example project with multiple transport modes  
✅ Math server implementation (reference for your servers)  
✅ Quick-start scripts for easy setup  
✅ Comprehensive documentation  
✅ Multiple ways to run MCP servers outside Docker  
✅ Production-ready patterns and examples  

**Start with:** `./quick-start.sh run-stdio` or `.\quick-start.ps1 run-stdio`

**Explore:** Different transports (TCP, discovery)

**Create:** Your own MCP servers (Go, Python, Node.js)

**Deploy:** Use TCP transport with system service

---

Last Updated: January 14, 2026  
AgenticGoKit v1beta
