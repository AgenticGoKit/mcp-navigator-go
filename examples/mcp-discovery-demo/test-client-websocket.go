package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
	"github.com/kunalkushwaha/mcp-navigator-go/pkg/mcp"
	"github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
)

// Test client for WebSocket transport
// Usage: go run test-client-websocket.go [-debug] [-port 9999]
// Make sure math-server is running with: go run math-server.go -mode websocket -port 9999

func main() {
	testWebSocketProtocol()
}

func testWebSocketProtocol() {
	debugMode := flag.Bool("debug", false, "Enable debug logging")
	port := flag.Int("port", 9999, "WebSocket port")
	flag.Parse()

	fmt.Println("=== MCP Navigator - WebSocket Protocol Test ===")
	fmt.Printf("Testing math-server on WebSocket port %d\n", *port)
	if *debugMode {
		fmt.Println("🐛 Debug mode: ENABLED")
	}
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create WebSocket transport
	wsURL := fmt.Sprintf("ws://localhost:%d/mcp", *port)
	trans := transport.NewWebSocketTransport(wsURL)

	// Create and connect client
	config := client.ClientConfig{
		Name:    "test-client-websocket",
		Version: "1.0.0",
		Debug:   *debugMode,
	}

	c := client.NewClient(trans, config)

	// Connect
	fmt.Println("📡 Connecting via WebSocket...")
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer c.Disconnect()
	fmt.Println("✅ Connected successfully")

	// Initialize
	fmt.Println("\n🔧 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "test-client-websocket",
		Version: "1.0.0",
	}

	if err := c.Initialize(ctx, clientInfo); err != nil {
		log.Fatalf("❌ Initialize failed: %v", err)
	}

	serverInfo := c.GetServerInfo()
	fmt.Printf("✅ Initialized with server: %s v%s\n", serverInfo.Name, serverInfo.Version)

	// List tools
	fmt.Println("\n🔍 Listing tools...")
	tools, err := c.ListTools(ctx)
	if err != nil {
		log.Fatalf("❌ ListTools failed: %v", err)
	}

	// Display results
	fmt.Printf("\n📊 Results: Found %d tools\n", len(tools))
	fmt.Println(strings.Repeat("=", 60))

	if len(tools) == 0 {
		fmt.Println("⚠️  No tools returned")
	} else {
		for i, tool := range tools {
			fmt.Printf("\n🔧 Tool #%d: %s\n", i+1, tool.Name)
			fmt.Printf("   Description: %s\n", tool.Description)

			if tool.InputSchema != nil {
				schemaJSON, _ := json.MarshalIndent(tool.InputSchema, "   ", "  ")
				fmt.Printf("   InputSchema: %s\n", string(schemaJSON))
			}

			// Test the tool
			if tool.Name == "divide" {
				fmt.Println("\n   🧪 Testing divide tool...")
				result, err := c.CallTool(ctx, "divide", map[string]interface{}{
					"a": 20,
					"b": 4,
				})
				if err != nil {
					fmt.Printf("   ❌ Tool call failed: %v\n", err)
				} else {
					fmt.Printf("   ✅ Result: %+v\n", result)
				}
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ WebSocket Test Completed Successfully")
}
