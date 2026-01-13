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

// Test client for STDIO transport
// Usage: go run test-client-stdio.go [-debug]
// This launches the math-server in stdio mode and communicates with it

func main() {
	testStdioProtocol()
}

func testStdioProtocol() {
	debugMode := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	fmt.Println("=== MCP Navigator - STDIO Protocol Test ===")
	fmt.Println("Testing math-server via STDIO transport")
	if *debugMode {
		fmt.Println("🐛 Debug mode: ENABLED")
	}
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create STDIO transport - launches math-server in stdio mode
	// The command and args should point to the math-server executable
	trans := transport.NewStdioTransport("./math-server.exe", []string{"-mode", "stdio"})

	// Create and connect client
	config := client.ClientConfig{
		Name:    "test-client-stdio",
		Version: "1.0.0",
		Debug:   *debugMode,
	}

	c := client.NewClient(trans, config)

	// Connect
	fmt.Println("📡 Connecting via STDIO...")
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer c.Disconnect()
	fmt.Println("✅ Connected successfully")

	// Initialize
	fmt.Println("\n🔧 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "test-client-stdio",
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
			if tool.Name == "multiply" {
				fmt.Println("\n   🧪 Testing multiply tool...")
				result, err := c.CallTool(ctx, "multiply", map[string]interface{}{
					"a": 4,
					"b": 6,
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
	fmt.Println("✅ STDIO Test Completed Successfully")
}
