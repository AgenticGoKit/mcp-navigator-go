# Quick Reference - Protocol Testing

## 🚀 Quick Start

### Test TCP (in separate terminals)
```bash
# Terminal 1
.\math-server.exe -mode tcp -port 9999

# Terminal 2
.\test-client-tcp.exe
```

### Test STDIO (single command)
```bash
.\test-client-stdio.exe
```

### Test WebSocket (in separate terminals)
```bash
# Terminal 1
.\math-server.exe -mode websocket -port 9999

# Terminal 2
.\test-client-websocket.exe
```

### Test All At Once
```powershell
.\run-all-protocol-tests.ps1
```

---

## 📊 Test Results Summary

| Protocol | Status | Connection | Tools | Sample Call |
|----------|--------|-----------|-------|-------------|
| TCP | ✅ | localhost:9999 | 4 found | add(5,3)=8 |
| STDIO | ✅ | subprocess | 4 found | multiply(4,6)=24 |
| WebSocket | ✅ | ws://localhost:9999/mcp | 4 found | divide(20,4)=5 |

---

## 🛠️ Build Commands

```bash
# Build everything
go build -o math-server.exe math-server.go
go build -o test-client-tcp.exe test-client-tcp.go
go build -o test-client-stdio.exe test-client-stdio.go
go build -o test-client-websocket.exe test-client-websocket.go

# Or use PowerShell script (auto-builds)
.\run-all-protocol-tests.ps1
```

---

## 🐛 Debug Mode

Add `-debug` flag to any test client for detailed logging:

```bash
.\test-client-tcp.exe -debug
.\test-client-stdio.exe -debug
.\test-client-websocket.exe -port 9999 -debug
```

Shows component-based logs:
- `[CLIENT]` - High-level operations
- `[INIT]` - Initialization details
- `[TRANSPORT]` - Low-level transport details
- `[PARSE]` - JSON parsing information

---

## 📁 Key Files

| File | Purpose |
|------|---------|
| math-server.go | MCP server (supports all 3 protocols) |
| test-client-tcp.go | TCP protocol test |
| test-client-stdio.go | STDIO protocol test |
| test-client-websocket.go | WebSocket protocol test |
| PROTOCOL_TESTING.md | Detailed testing guide |
| PROTOCOL_TEST_RESULTS.md | Test results documentation |
| run-all-protocol-tests.ps1 | Automated test suite |

---

## 🎯 Available Tools

All protocols expose these 4 math tools:

```
1. add        → add(5, 3) = 8
2. subtract   → subtract(10, 3) = 7
3. multiply   → multiply(4, 6) = 24
4. divide     → divide(20, 4) = 5
```

Each tool has:
- Name & Description
- InputSchema (JSON schema with required parameters)
- Callable via CallTool

---

## ❓ Troubleshooting

### Port Already in Use
```powershell
Get-NetTCPConnection -LocalPort 9999 | Stop-Process -Force
```

### Build Errors
```bash
# Ensure Go is installed
go version

# Ensure dependencies are available
go mod download

# Then rebuild
go build -o math-server.exe math-server.go
```

### Connection Refused (TCP/WebSocket)
1. Verify server is running in correct mode
2. Check port number matches (default: 9999)
3. Ensure firewall allows port

### STDIO Test Fails
- Ensure `math-server.exe` is in same directory
- Or use full path in test-client-stdio.go

---

## 🔄 Protocol Comparison

| Aspect | TCP | STDIO | WebSocket |
|--------|-----|-------|-----------|
| Setup | Easy | Easiest | Medium |
| Latency | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| Remote | ✅ Yes | ❌ No | ✅ Yes |
| Port Required | ✅ Yes | ❌ No | ✅ Yes |
| Use Case | Fast local | Tool integration | Web services |

---

## ✅ What's Tested

Each protocol test verifies:

1. **Connection** - Can connect to server
2. **Initialization** - MCP handshake succeeds
3. **Discovery** - Gets all 4 tools
4. **Schemas** - Tool schemas are complete
5. **Execution** - Can call tools and get results
6. **Compliance** - Meets MCP spec (2025-11-25)

---

## 📚 Documentation Files

- **PROTOCOL_TESTING.md** - How to test, detailed instructions
- **PROTOCOL_TEST_RESULTS.md** - Results from testing all 3 protocols
- **PROTOCOL_IMPLEMENTATION_SUMMARY.md** - What was implemented and why

---

## 🎓 Learning Path

1. **Start Here** - Read this file (Quick Reference)
2. **Learn Details** - Read PROTOCOL_TESTING.md
3. **See Results** - Review PROTOCOL_TEST_RESULTS.md
4. **Understand Code** - Review test client implementations
5. **Run Tests** - Execute tests yourself

---

## 💡 Key Takeaways

✅ **All protocols working** - TCP, STDIO, WebSocket all functional
✅ **Full spec compliance** - Meets MCP 2025-11-25 requirements  
✅ **Easy to test** - Individual test clients for each protocol
✅ **Automated testing** - PowerShell script tests all at once
✅ **Well documented** - Clear examples and guides
✅ **Production ready** - Ready for deployment with any protocol

---

## 🚀 Next Steps

1. Run `.\run-all-protocol-tests.ps1` to verify all protocols
2. Review results in PROTOCOL_TEST_RESULTS.md
3. Use test clients as templates for your own servers
4. Choose protocol based on your needs (TCP=fast, STDIO=tools, WS=web)
5. Deploy mcp-navigator-go with confidence!

