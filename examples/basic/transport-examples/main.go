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

// Example demonstrates different transport types for MCP connections
func main() {
	fmt.Println("🔌 MCP Transport Examples")
	fmt.Println("=========================")

	// 1. TCP Transport
	if err := exampleTCPClient(); err != nil {
		log.Printf("TCP client example failed: %v", err)
	}

	fmt.Println()

	// 2. HTTP SSE Transport
	if err := exampleHTTPSSEClient(); err != nil {
		log.Printf("HTTP SSE client example failed: %v", err)
	}

	fmt.Println()

	// 3. HTTP Streaming Transport
	if err := exampleHTTPStreamingClient(); err != nil {
		log.Printf("HTTP Streaming client example failed: %v", err)
	}

	fmt.Println()

	// 4. WebSocket Transport (if available)
	if err := exampleWebSocketClient(); err != nil {
		log.Printf("WebSocket client example failed: %v", err)
	}
}

// exampleTCPClient demonstrates TCP transport usage
func exampleTCPClient() error {
	fmt.Println("1. TCP Transport Example")
	fmt.Println("------------------------")

	// Create TCP transport
	tcpTransport := transport.NewTCPTransport("localhost", 8811)

	// Create client configuration
	config := client.ClientConfig{
		Name:    "transport-example-tcp",
		Version: "1.0.0",
		Logger:  log.Default(),
		Timeout: 30 * time.Second,
	}

	// Create client
	mcpClient := client.NewClient(tcpTransport, config)

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("📡 Connecting to TCP MCP server...")
	if err := mcpClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		fmt.Println("🔌 Disconnecting...")
		mcpClient.Disconnect()
	}()

	fmt.Println("✅ Connected successfully!")

	// Initialize protocol
	fmt.Println("🚀 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "transport-example-tcp",
		Version: "1.0.0",
	}

	if err := mcpClient.Initialize(ctx, clientInfo); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Get server info
	if serverInfo := mcpClient.GetServerInfo(); serverInfo != nil {
		fmt.Printf("🖥️  Server: %s v%s\n", serverInfo.Name, serverInfo.Version)
	}

	// List available tools
	return listAndCallTools(ctx, mcpClient, "TCP")
}

// exampleHTTPSSEClient demonstrates HTTP SSE transport usage
func exampleHTTPSSEClient() error {
	fmt.Println("2. HTTP SSE Transport Example")
	fmt.Println("-----------------------------")

	// Create HTTP SSE transport - for real-time communication
	sseTransport := transport.NewSSETransport("http://localhost:8812", "/sse/")

	// Create client configuration
	config := client.ClientConfig{
		Name:    "transport-example-sse",
		Version: "1.0.0",
		Logger:  log.Default(),
		Timeout: 30 * time.Second,
	}

	// Create client
	mcpClient := client.NewClient(sseTransport, config)

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("📡 Connecting to HTTP SSE MCP server...")
	if err := mcpClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		fmt.Println("🔌 Disconnecting...")
		mcpClient.Disconnect()
	}()

	fmt.Println("✅ Connected successfully!")

	// Initialize protocol
	fmt.Println("🚀 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "transport-example-sse",
		Version: "1.0.0",
	}

	if err := mcpClient.Initialize(ctx, clientInfo); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Get server info
	if serverInfo := mcpClient.GetServerInfo(); serverInfo != nil {
		fmt.Printf("🖥️  Server: %s v%s\n", serverInfo.Name, serverInfo.Version)
	}

	// List available tools
	return listAndCallTools(ctx, mcpClient, "HTTP SSE")
}

// exampleHTTPStreamingClient demonstrates HTTP Streaming transport usage
func exampleHTTPStreamingClient() error {
	fmt.Println("3. HTTP Streaming Transport Example")
	fmt.Println("-----------------------------------")

	// Create HTTP Streaming transport - for request-response
	streamingTransport := transport.NewStreamingHTTPTransport("http://localhost:8813", "/mcp")

	// Create client configuration
	config := client.ClientConfig{
		Name:    "transport-example-streaming",
		Version: "1.0.0",
		Logger:  log.Default(),
		Timeout: 30 * time.Second,
	}

	// Create client
	mcpClient := client.NewClient(streamingTransport, config)

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("📡 Connecting to HTTP Streaming MCP server...")
	if err := mcpClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		fmt.Println("🔌 Disconnecting...")
		mcpClient.Disconnect()
	}()

	fmt.Println("✅ Connected successfully!")

	// Initialize protocol
	fmt.Println("🚀 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "transport-example-streaming",
		Version: "1.0.0",
	}

	if err := mcpClient.Initialize(ctx, clientInfo); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Get server info
	if serverInfo := mcpClient.GetServerInfo(); serverInfo != nil {
		fmt.Printf("🖥️  Server: %s v%s\n", serverInfo.Name, serverInfo.Version)
	}

	// List available tools
	return listAndCallTools(ctx, mcpClient, "HTTP Streaming")
}

// exampleWebSocketClient demonstrates WebSocket transport usage
func exampleWebSocketClient() error {
	fmt.Println("4. WebSocket Transport Example")
	fmt.Println("------------------------------")

	// Create WebSocket transport
	wsTransport := transport.NewWebSocketTransport("ws://localhost:8814/mcp")

	// Create client configuration
	config := client.ClientConfig{
		Name:    "transport-example-ws",
		Version: "1.0.0",
		Logger:  log.Default(),
		Timeout: 30 * time.Second,
	}

	// Create client
	mcpClient := client.NewClient(wsTransport, config)

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("📡 Connecting to WebSocket MCP server...")
	if err := mcpClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect (WebSocket server may not be available): %w", err)
	}
	defer func() {
		fmt.Println("🔌 Disconnecting...")
		mcpClient.Disconnect()
	}()

	fmt.Println("✅ Connected successfully!")

	// Initialize protocol
	fmt.Println("🚀 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "transport-example-ws",
		Version: "1.0.0",
	}

	if err := mcpClient.Initialize(ctx, clientInfo); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Get server info
	if serverInfo := mcpClient.GetServerInfo(); serverInfo != nil {
		fmt.Printf("🖥️  Server: %s v%s\n", serverInfo.Name, serverInfo.Version)
	}

	// List available tools
	return listAndCallTools(ctx, mcpClient, "WebSocket")
}

// listAndCallTools is a helper function to list tools and call one if available
func listAndCallTools(ctx context.Context, mcpClient *client.Client, transportType string) error {
	fmt.Printf("📋 Listing available tools on %s server...\n", transportType)
	tools, err := mcpClient.ListTools(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}

	fmt.Printf("📝 Found %d tools:\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("   %d. %s - %s\n", i+1, tool.Name, tool.Description)
	}

	// Call a tool if available
	if len(tools) > 0 {
		toolName := tools[0].Name
		fmt.Printf("🔧 Calling tool: %s\n", toolName)

		// Example arguments for search tool
		arguments := map[string]interface{}{
			"query": fmt.Sprintf("example search from %s transport", transportType),
		}

		result, err := mcpClient.CallTool(ctx, toolName, arguments)
		if err != nil {
			fmt.Printf("⚠️  Tool execution failed: %v\n", err)
			return nil // Don't fail the example for tool errors
		}

		fmt.Printf("✅ Tool result received (%d content items)\n", len(result.Content))
		if len(result.Content) > 0 {
			fmt.Printf("   First item type: %s\n", result.Content[0].Type)
			if result.Content[0].Text != "" {
				// Truncate long responses for example output
				text := result.Content[0].Text
				if len(text) > 200 {
					text = text[:200] + "..."
				}
				fmt.Printf("   Content preview: %s\n", text)
			}
		}
	}

	// List available resources
	fmt.Printf("📂 Listing available resources on %s server...\n", transportType)
	resources, err := mcpClient.ListResources(ctx)
	if err != nil {
		fmt.Printf("⚠️  Failed to list resources: %v\n", err)
		return nil // Don't fail the example for resource errors
	}

	fmt.Printf("📁 Found %d resources:\n", len(resources))
	for i, resource := range resources {
		fmt.Printf("   %d. %s - %s\n", i+1, resource.URI, resource.Name)
	}

	return nil
}
