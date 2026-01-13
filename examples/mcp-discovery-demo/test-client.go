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

// Simple test client to verify the ListTools fix
func main() {
	debugMode := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	fmt.Println("=== MCP Navigator Test Client ===")
	fmt.Println("Testing ListTools() with TCP transport on port 9999")
	if *debugMode {
		fmt.Println("🐛 Debug mode: ENABLED")
	} else {
		fmt.Println("📝 Debug mode: DISABLED (use -debug to enable)")
	}
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create TCP transport
	trans := transport.NewTCPTransport("localhost", 9999)

	// Create client configuration
	config := client.ClientConfig{
		Name:    "test-client",
		Version: "1.0.0",
		Debug:   *debugMode, // Enable debug based on flag
	}

	// Create client
	c := client.NewClient(trans, config)

	// Connect to server
	fmt.Println("📡 Connecting to MCP server on localhost:9999...")
	if err := c.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer c.Disconnect()
	fmt.Println("✅ Connected successfully")

	// Initialize the MCP protocol
	fmt.Println("\n🔧 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "test-client",
		Version: "1.0.0",
	}

	if err := c.Initialize(ctx, clientInfo); err != nil {
		log.Fatalf("❌ Initialize failed: %v", err)
	}

	serverInfo := c.GetServerInfo()
	if serverInfo != nil {
		fmt.Printf("✅ Initialized with server: %s v%s\n", serverInfo.Name, serverInfo.Version)
	} else {
		fmt.Println("✅ Initialized successfully")
	}

	// List available tools
	fmt.Println("\n🔍 Listing available tools...")
	tools, err := c.ListTools(ctx)
	if err != nil {
		log.Fatalf("❌ ListTools failed: %v", err)
	}

	// Display results
	fmt.Printf("\n📊 Results: Found %d tools\n", len(tools))
	fmt.Println(strings.Repeat("=", 60))

	if len(tools) == 0 {
		fmt.Println("⚠️  WARNING: No tools returned (this is the bug!)")
	} else {
		fmt.Println("✅ SUCCESS: Tools were returned!")
		for i, tool := range tools {
			fmt.Printf("\n🔧 Tool #%d:\n", i+1)
			fmt.Printf("   Name:        %s\n", tool.Name)
			fmt.Printf("   Description: %s\n", tool.Description)

			// Display input schema
			if tool.InputSchema != nil {
				schemaJSON, err := json.MarshalIndent(tool.InputSchema, "   ", "  ")
				if err == nil {
					fmt.Printf("   InputSchema:\n   %s\n", string(schemaJSON))
				}
			}

			// Validate tool per MCP spec
			if err := mcp.ValidateTool(&tool); err != nil {
				fmt.Printf("   ⚠️  Validation: %v\n", err)
			} else {
				fmt.Printf("   ✅ Validation: Passed MCP spec compliance\n")
			}
		}
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n✅ Test completed successfully!")
}
