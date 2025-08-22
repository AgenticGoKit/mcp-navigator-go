//go:build example
// +build example

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
	"github.com/kunalkushwaha/mcp-navigator-go/pkg/mcp"
	"github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
)

// Example showing advanced patterns for MCP client usage
func main() {
	fmt.Println("🏗️  MCP Client - Advanced Patterns")
	fmt.Println("==================================")

	// Example 1: Simple client creation
	simpleExample()

	fmt.Println()

	// Example 2: HTTP client with different transports
	httpExamples()

	fmt.Println()

	// Example 3: Service wrapper pattern
	serviceExample()

	fmt.Println()

	// Example 4: Error handling and retry logic
	robustExample()
}

// simpleExample shows the simplest way to create a client
func simpleExample() {
	fmt.Println("1. Simple Client Creation")
	fmt.Println("------------------------")

	// Create client with TCP transport
	tcpTransport := transport.NewTCPTransport("localhost", 8811)
	config := client.ClientConfig{
		Name:    "simple-app",
		Version: "1.0.0",
		Timeout: 10 * time.Second,
	}
	mcpClient := client.NewClient(tcpTransport, config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("📡 Connecting...")
	if err := mcpClient.Connect(ctx); err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		return
	}
	defer mcpClient.Disconnect()

	clientInfo := mcp.ClientInfo{
		Name:    "simple-app",
		Version: "1.0.0",
	}
	err := mcpClient.Initialize(ctx, clientInfo)
	if err != nil {
		fmt.Printf("❌ Initialization failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Connected successfully!\n")

	// Quick tool listing
	tools, err := mcpClient.ListTools(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to list tools: %v\n", err)
		return
	}

	fmt.Printf("📝 Found %d tools\n", len(tools))
}

// httpExamples shows different HTTP transport options
func httpExamples() {
	fmt.Println("2. HTTP Transport Patterns")
	fmt.Println("--------------------------")

	examples := []struct {
		name      string
		transport transport.Transport
		use_case  string
	}{
		{
			name:      "HTTP SSE",
			transport: transport.NewSSETransport("http://localhost:8812", "/sse/"),
			use_case:  "Real-time communication",
		},
		{
			name:      "HTTP Streaming",
			transport: transport.NewStreamingHTTPTransport("http://localhost:8813", "/mcp"),
			use_case:  "Request-response pattern",
		},
	}

	for _, example := range examples {
		fmt.Printf("📡 Testing %s transport (%s)...\n", example.name, example.use_case)

		config := client.ClientConfig{
			Name:    "http-pattern-app",
			Version: "1.0.0",
			Timeout: 10 * time.Second,
		}
		mcpClient := client.NewClient(example.transport, config)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		if err := mcpClient.Connect(ctx); err != nil {
			fmt.Printf("❌ %s connection failed: %v\n", example.name, err)
			cancel()
			continue
		}

		clientInfo := mcp.ClientInfo{
			Name:    "http-pattern-app",
			Version: "1.0.0",
		}
		err := mcpClient.Initialize(ctx, clientInfo)
		if err != nil {
			fmt.Printf("❌ %s initialization failed: %v\n", example.name, err)
			mcpClient.Disconnect()
			cancel()
			continue
		}

		fmt.Printf("✅ %s connected successfully!\n", example.name)

		mcpClient.Disconnect()
		cancel()
	}
}

// serviceExample shows a service wrapper pattern
func serviceExample() {
	fmt.Println("3. Service Wrapper Pattern")
	fmt.Println("--------------------------")

	service := NewMCPService()

	fmt.Println("📡 Service connecting...")
	if err := service.Connect("localhost", 8811); err != nil {
		fmt.Printf("❌ Service connection failed: %v\n", err)
		return
	}
	defer service.Disconnect()

	fmt.Println("✅ Service connected!")

	// Use the service
	tools, err := service.ListTools()
	if err != nil {
		fmt.Printf("❌ Failed to list tools: %v\n", err)
		return
	}

	fmt.Printf("📝 Service found %d tools\n", len(tools))

	// Execute a tool through the service
	if len(tools) > 0 {
		fmt.Printf("🔧 Executing tool: %s\n", tools[0].Name)
		result, err := service.ExecuteTool(tools[0].Name, map[string]interface{}{
			"query": "test from service wrapper",
		})
		if err != nil {
			fmt.Printf("❌ Tool execution failed: %v\n", err)
		} else {
			fmt.Printf("✅ Tool executed successfully (%d content items)\n", len(result.Content))
		}
	}
}

// robustExample shows error handling and retry logic
func robustExample() {
	fmt.Println("4. Robust Client with Error Handling")
	fmt.Println("------------------------------------")

	// Try multiple connection options with fallback
	connectionOptions := []struct {
		name      string
		transport transport.Transport
		priority  int
	}{
		{"TCP", transport.NewTCPTransport("localhost", 8811), 1},
		{"HTTP SSE", transport.NewSSETransport("http://localhost:8812", "/sse/"), 2},
		{"HTTP Streaming", transport.NewStreamingHTTPTransport("http://localhost:8813", "/mcp"), 3},
	}

	for _, option := range connectionOptions {
		fmt.Printf("🔄 Trying %s connection (priority %d)...\n", option.name, option.priority)

		config := client.ClientConfig{
			Name:    "robust-app",
			Version: "1.0.0",
			Timeout: 5 * time.Second,
		}
		mcpClient := client.NewClient(option.transport, config)

		// Retry logic with exponential backoff
		var err error
		for attempt := 1; attempt <= 3; attempt++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			err = mcpClient.Connect(ctx)
			if err == nil {
				clientInfo := mcp.ClientInfo{
					Name:    "robust-app",
					Version: "1.0.0",
				}
				initErr := mcpClient.Initialize(ctx, clientInfo)
				if initErr == nil {
					fmt.Printf("✅ Successfully connected via %s!\n", option.name)

					// Test functionality
					tools, toolErr := mcpClient.ListTools(ctx)
					if toolErr == nil && len(tools) > 0 {
						fmt.Printf("🔧 Successfully verified %d tools available\n", len(tools))
					}

					mcpClient.Disconnect()
					cancel()
					return
				}
				err = initErr
			}

			cancel()
			if attempt < 3 {
				backoff := time.Duration(attempt*attempt) * time.Second
				fmt.Printf("   Attempt %d failed, retrying in %v...\n", attempt, backoff)
				time.Sleep(backoff)
			}
		}

		fmt.Printf("❌ %s connection failed after 3 attempts: %v\n", option.name, err)
	}

	fmt.Println("❌ All connection options exhausted - check server availability")
}

// MCPService demonstrates a service wrapper pattern for production use
type MCPService struct {
	client    *client.Client
	connected bool
	config    ServiceConfig
}

// ServiceConfig holds configuration for the MCP service
type ServiceConfig struct {
	Host           string
	Port           int
	ClientName     string
	ClientVersion  string
	ConnectTimeout time.Duration
	OperateTimeout time.Duration
}

// NewMCPService creates a new MCP service wrapper with default config
func NewMCPService() *MCPService {
	return &MCPService{
		config: ServiceConfig{
			Host:           "localhost",
			Port:           8811,
			ClientName:     "mcp-service",
			ClientVersion:  "1.0.0",
			ConnectTimeout: 30 * time.Second,
			OperateTimeout: 10 * time.Second,
		},
	}
}

// NewMCPServiceWithConfig creates a service with custom configuration
func NewMCPServiceWithConfig(config ServiceConfig) *MCPService {
	return &MCPService{config: config}
}

// Connect establishes connection to an MCP server
func (s *MCPService) Connect(host string, port int) error {
	s.config.Host = host
	s.config.Port = port

	tcpTransport := transport.NewTCPTransport(host, port)
	config := client.ClientConfig{
		Name:    s.config.ClientName,
		Version: s.config.ClientVersion,
		Timeout: s.config.ConnectTimeout,
	}
	s.client = client.NewClient(tcpTransport, config)

	ctx, cancel := context.WithTimeout(context.Background(), s.config.ConnectTimeout)
	defer cancel()

	if err := s.client.Connect(ctx); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	clientInfo := mcp.ClientInfo{
		Name:    s.config.ClientName,
		Version: s.config.ClientVersion,
	}
	err := s.client.Initialize(ctx, clientInfo)
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	s.connected = true
	return nil
}

// Disconnect closes the connection gracefully
// Disconnect closes the connection
func (s *MCPService) Disconnect() error {
	if s.client != nil {
		s.connected = false
		return s.client.Disconnect()
	}
	return nil
}

// ListTools returns available tools with timeout
func (s *MCPService) ListTools() ([]mcp.Tool, error) {
	if !s.connected {
		return nil, fmt.Errorf("service not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.config.OperateTimeout)
	defer cancel()

	return s.client.ListTools(ctx)
}

// ExecuteTool executes a tool with the given arguments and timeout
func (s *MCPService) ExecuteTool(name string, args map[string]interface{}) (*mcp.CallToolResponse, error) {
	if !s.connected {
		return nil, fmt.Errorf("service not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.config.OperateTimeout*3) // Longer timeout for tool execution
	defer cancel()

	return s.client.CallTool(ctx, name, args)
}

// IsConnected returns the connection status
func (s *MCPService) IsConnected() bool {
	return s.connected && s.client != nil
}

// GetServerInfo returns information about the connected server
func (s *MCPService) GetServerInfo() (string, string, error) {
	if !s.connected {
		return "", "", fmt.Errorf("service not connected")
	}

	// In a real implementation, you'd store server info during initialization
	return "connected-server", "unknown-version", nil
}

// HealthCheck performs a basic health check on the connection
func (s *MCPService) HealthCheck() error {
	if !s.connected {
		return fmt.Errorf("service not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.client.ListTools(ctx)
	return err
}
