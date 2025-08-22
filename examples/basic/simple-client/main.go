//go:build example
// +build example

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
	"github.com/kunalkushwaha/mcp-navigator-go/pkg/mcp"
	"github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
)

// Simple example demonstrating basic MCP client usage with TCP transport
func main() {
	fmt.Println("🚀 Simple MCP Client Example")
	fmt.Println("============================")

	// Create TCP transport
	tcpTransport := transport.NewTCPTransport("localhost", 8811)

	// Create client configuration
	config := client.ClientConfig{
		Name:    "simple-client-example",
		Version: "1.0.0",
		Logger:  log.Default(),
		Timeout: 30 * time.Second,
	}

	// Create MCP client
	mcpClient := client.NewClient(tcpTransport, config)

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("📡 Connecting to MCP server...")
	if err := mcpClient.Connect(ctx); err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer func() {
		fmt.Println("🔌 Disconnecting...")
		mcpClient.Disconnect()
	}()

	fmt.Println("✅ Connected successfully!")

	// Initialize MCP protocol
	fmt.Println("🚀 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "simple-client-example",
		Version: "1.0.0",
	}

	if err := mcpClient.Initialize(ctx, clientInfo); err != nil {
		log.Fatalf("❌ Failed to initialize: %v", err)
	}

	// Get server information
	if serverInfo := mcpClient.GetServerInfo(); serverInfo != nil {
		fmt.Printf("🖥️  Server: %s v%s\n", serverInfo.Name, serverInfo.Version)
	}

	// List available tools
	fmt.Println("📋 Listing available tools...")
	tools, err := mcpClient.ListTools(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to list tools: %v", err)
	}

	fmt.Printf("📝 Found %d tools:\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("   %d. %s - %s\n", i+1, tool.Name, tool.Description)
	}

	// Execute first tool if available
	if len(tools) > 0 {
		toolName := tools[0].Name
		fmt.Printf("🔧 Calling tool: %s\n", toolName)

		// Simple arguments - adjust based on your server's tools
		arguments := map[string]interface{}{
			"query": "Hello from simple client example",
		}

		result, err := mcpClient.CallTool(ctx, toolName, arguments)
		if err != nil {
			log.Printf("⚠️  Tool execution failed: %v", err)
		} else {
			fmt.Printf("✅ Tool executed successfully!\n")
			fmt.Printf("   Result contains %d content item(s)\n", len(result.Content))

			// Show first content item if available
			if len(result.Content) > 0 && result.Content[0].Text != "" {
				text := result.Content[0].Text
				if len(text) > 100 {
					text = text[:100] + "..."
				}
				fmt.Printf("   Preview: %s\n", text)
			}
		}
	}

	fmt.Println("✅ Example completed successfully!")
}
