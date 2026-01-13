# Changelog

All notable changes to MCP Navigator will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [v2.0.0] - 2026-01-14

### Added
- **WebSocket Transport Support** - New WebSocket transport for web-compatible MCP servers
- **Configurable Debug Logging** - Optional debug mode with component-based logging (disabled by default for clean production output)
- **Component Prefixes in Logs** - [CLIENT], [INIT], [TRANSPORT], [PARSE] prefixes for easier troubleshooting
- **MCP 2025-11-25 Compliance** - Full specification compliance with pagination support and optional fields
- **Protocol Testing Suite** - Comprehensive test clients for all supported protocols (TCP, STDIO, WebSocket)
- **Tool Icons Support** - Support for icon specifications in tool definitions
- **Tool Pagination** - Cursor-based pagination for tool listings (ListToolsWithCursor method)
- **OutputSchema Support** - Optional output schema field for tools per MCP spec

### Fixed
- **ListTools() Returning Empty Array** - Fixed critical bug where ListTools() would return empty array despite server sending tools
  - Root cause: Improved JSON parsing and validation of tool responses
  - Added comprehensive DEBUG logging to trace data flow
  - Added MCP validation functions (ValidateToolName, ValidateTool)
- **Verbose Logging in Production** - Wrapped all verbose client logs with debug checks to prevent log pollution in production
- **Protocol Compliance** - Fixed server implementation issues in example math-server:
  - Proper JSON-RPC error format with {code, message} structure
  - Correct notification handling (notifications/initialized)
- **Transport Initialization** - Improved connection state management across all transports

### Changed
- **Client.ListTools()** - Now uses ListToolsWithCursor internally for pagination support
- **Debug Logging** - All logging is now conditional on `Debug: true` in ClientConfig
- **Logger Helper** - Added `Client.logf()` method for structured component-based logging in library
- **Test Client Examples** - Updated with modern MCP client initialization patterns

### Improved
- **Performance** - No logging overhead in production (debug=false)
- **Observability** - Detailed logs available on-demand for troubleshooting
- **Documentation** - Added logging guides and protocol testing documentation
- **Test Coverage** - Added test clients for all transport protocols with working examples

### Documentation
- [LOG_IMPROVEMENTS.md](pkg/client/LOG_IMPROVEMENTS.md) - Guide to configurable debug logging
- [PROTOCOL_TESTING.md](examples/mcp-discovery-demo/PROTOCOL_TESTING.md) - How to test all protocols
- [PROTOCOL_TEST_RESULTS.md](examples/mcp-discovery-demo/PROTOCOL_TEST_RESULTS.md) - Complete test results
- [QUICK_REFERENCE.md](examples/mcp-discovery-demo/QUICK_REFERENCE.md) - Quick reference for testing

### Validation
- ✅ TCP Protocol: Tested and working
- ✅ STDIO Protocol: Tested and working
- ✅ WebSocket Protocol: Tested and working
- ✅ Tool Discovery: All 4 tools returned with complete schemas
- ✅ Tool Execution: CallTool works across all protocols
- ✅ MCP Spec Compliance: Validated against 2025-11-25 specification

## [v1.0.0] - 2025-06-08

### Added
- Initial stable release
- Complete MCP protocol implementation
- Production-ready library and CLI
- Comprehensive test suite
- Full documentation

---

## Release Notes

### v1.0.0 - Complete MCP Implementation

This is the first stable release of MCP Navigator, providing complete Model Context Protocol support with both CLI and library interfaces.

**Key Highlights:**
- 🎯 **Production Ready**: 95% library readiness score
- 🔧 **Complete Protocol**: Full MCP support (Tools, Resources, Prompts)
- 🔍 **Unique Discovery**: Server discovery capabilities not found in other implementations
- 📚 **Library + CLI**: Dual-purpose design for both applications and automation
- ⚡ **High Performance**: Native Go implementation with excellent performance

**Library Features:**
- Thread-safe concurrent usage
- Fluent builder pattern for easy configuration
- Comprehensive error handling
- Extensive examples and documentation
- Production-tested with real MCP servers

**CLI Features:**
- Interactive mode with auto-completion
- Server discovery and health testing
- Direct tool execution from command line
- Docker container support
- Configurable timeouts and retries

**Getting Started:**
- Library: `go get github.com/kunalkushwaha/mcp-navigator-go`
- CLI: `go install github.com/kunalkushwaha/mcp-navigator-go@latest`
- Source: `git clone https://github.com/kunalkushwaha/mcp-navigator-go.git`

This release marks the completion of the core MCP implementation with all major features in place for production use.
