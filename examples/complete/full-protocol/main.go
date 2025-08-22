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

// Example demonstrating complete MCP protocol support including Tools, Resources, and Prompts
func main() {
	fmt.Println("🚀 MCP Client - Complete Protocol Features")
	fmt.Println("==========================================")

	// Create client with TCP transport
	tcpTransport := transport.NewTCPTransport("localhost", 8811)
	config := client.ClientConfig{
		Name:    "complete-features-demo",
		Version: "1.0.0",
		Timeout: 30 * time.Second,
	}
	mcpClient := client.NewClient(tcpTransport, config)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect and initialize
	fmt.Println("📡 Connecting to MCP server...")
	if err := mcpClient.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer mcpClient.Disconnect()

	fmt.Println("🚀 Initializing MCP protocol...")
	clientInfo := mcp.ClientInfo{
		Name:    "complete-features-example",
		Version: "1.0.0",
	}
	err := mcpClient.Initialize(ctx, clientInfo)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	fmt.Printf("🖥️  Connected successfully!\n\n")

	// Demonstrate all MCP features
	demonstrateTools(ctx, mcpClient)
	fmt.Println()
	demonstrateResources(ctx, mcpClient)
	fmt.Println()
	demonstratePrompts(ctx, mcpClient)
	fmt.Println()
	demonstrateAdvancedUsage(ctx, mcpClient)
}

// demonstrateTools shows comprehensive tool usage
func demonstrateTools(ctx context.Context, mcpClient *client.Client) {
	fmt.Println("🛠️  1. Tools Support")
	fmt.Println("===================")

	// List available tools
	fmt.Println("📋 Listing available tools...")
	tools, err := mcpClient.ListTools(ctx)
	if err != nil {
		log.Printf("❌ Failed to list tools: %v", err)
		return
	}

	fmt.Printf("📝 Found %d tools:\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("   %d. %s\n", i+1, tool.Name)
		fmt.Printf("      Description: %s\n", tool.Description)
		if tool.InputSchema != nil {
			fmt.Printf("      Schema available: Yes\n")
		}
	}

	// Execute tools with different argument patterns
	if len(tools) > 0 {
		// Execute first tool
		toolName := tools[0].Name
		fmt.Printf("\n🔧 Executing tool '%s'...\n", toolName)

		arguments := map[string]interface{}{
			"query": "Complete features example test",
		}

		result, err := mcpClient.CallTool(ctx, toolName, arguments)
		if err != nil {
			log.Printf("❌ Failed to execute tool: %v", err)
		} else {
			fmt.Printf("✅ Tool execution successful!\n")
			fmt.Printf("   Content items: %d\n", len(result.Content))
			fmt.Printf("   Is error: %t\n", result.IsError)

			// Show content details
			for i, content := range result.Content {
				fmt.Printf("   Content %d: Type=%s\n", i+1, content.Type)
				if content.Text != "" {
					text := content.Text
					if len(text) > 100 {
						text = text[:100] + "..."
					}
					fmt.Printf("   Preview: %s\n", text)
				}
			}
		}

		// Execute multiple tools if available
		if len(tools) > 1 {
			fmt.Printf("\n🔧 Executing multiple tools sequentially...\n")
			for i := 1; i < len(tools) && i < 3; i++ { // Execute up to 3 tools
				toolName := tools[i].Name
				fmt.Printf("   Executing %s...", toolName)

				result, err := mcpClient.CallTool(ctx, toolName, map[string]interface{}{
					"query": fmt.Sprintf("test from tool %d", i),
				})
				if err != nil {
					fmt.Printf(" ❌ Failed: %v\n", err)
				} else {
					fmt.Printf(" ✅ Success (%d items)\n", len(result.Content))
				}
			}
		}
	}
}

// demonstrateResources shows resource management
func demonstrateResources(ctx context.Context, mcpClient *client.Client) {
	fmt.Println("📂 2. Resources Support")
	fmt.Println("======================")

	// List available resources
	fmt.Println("📋 Listing available resources...")
	resources, err := mcpClient.ListResources(ctx)
	if err != nil {
		log.Printf("❌ Failed to list resources: %v", err)
		return
	}

	fmt.Printf("📁 Found %d resources:\n", len(resources))
	for i, resource := range resources {
		fmt.Printf("   %d. %s\n", i+1, resource.Name)
		fmt.Printf("      URI: %s\n", resource.URI)
		fmt.Printf("      Description: %s\n", resource.Description)
		if resource.MimeType != "" {
			fmt.Printf("      MIME Type: %s\n", resource.MimeType)
		}
	}

	// Read resources if available
	if len(resources) > 0 {
		fmt.Printf("\n📖 Reading resources...\n")
		for i, resource := range resources {
			if i >= 3 { // Limit to first 3 resources
				break
			}

			fmt.Printf("   Reading '%s'...", resource.Name)
			content, err := mcpClient.ReadResource(ctx, resource.URI)
			if err != nil {
				fmt.Printf(" ❌ Failed: %v\n", err)
			} else {
				fmt.Printf(" ✅ Success (%d content items)\n", len(content.Contents))
			}
		}
	}

	if len(resources) == 0 {
		fmt.Println("ℹ️  No resources available on this server")
	}
}

// demonstratePrompts shows prompt functionality
func demonstratePrompts(ctx context.Context, mcpClient *client.Client) {
	fmt.Println("💬 3. Prompts Support")
	fmt.Println("====================")

	// List available prompts
	fmt.Println("📋 Listing available prompts...")
	prompts, err := mcpClient.ListPrompts(ctx)
	if err != nil {
		log.Printf("❌ Failed to list prompts: %v", err)
		return
	}

	fmt.Printf("📝 Found %d prompts:\n", len(prompts))
	for i, prompt := range prompts {
		fmt.Printf("   %d. %s\n", i+1, prompt.Name)
		fmt.Printf("      Description: %s\n", prompt.Description)

		if len(prompt.Arguments) > 0 {
			fmt.Printf("      Arguments:\n")
			for _, arg := range prompt.Arguments {
				required := ""
				if arg.Required {
					required = " (required)"
				}
				fmt.Printf("        - %s%s: %s\n", arg.Name, required, arg.Description)
			}
		}
	}

	// Get prompts if available
	if len(prompts) > 0 {
		fmt.Printf("\n💭 Getting prompts...\n")
		for i, prompt := range prompts {
			if i >= 2 { // Limit to first 2 prompts
				break
			}

			fmt.Printf("   Getting '%s'...", prompt.Name)

			// Prepare arguments
			arguments := make(map[string]interface{})
			for _, arg := range prompt.Arguments {
				if arg.Required {
					// Provide sample values for required arguments
					arguments[arg.Name] = fmt.Sprintf("sample_%s_value", arg.Name)
				}
			}

			result, err := mcpClient.GetPrompt(ctx, prompt.Name, arguments)
			if err != nil {
				fmt.Printf(" ❌ Failed: %v\n", err)
			} else {
				fmt.Printf(" ✅ Success (%d messages)\n", len(result.Messages))

				// Show message details
				for j, message := range result.Messages {
					if j >= 2 { // Limit output
						break
					}
					fmt.Printf("      Message %d: Role=%s\n", j+1, message.Role)
				}
			}
		}
	}

	if len(prompts) == 0 {
		fmt.Println("ℹ️  No prompts available on this server")
	}
}

// demonstrateAdvancedUsage shows advanced protocol features
func demonstrateAdvancedUsage(ctx context.Context, mcpClient *client.Client) {
	fmt.Println("⚡ 4. Advanced Usage Patterns")
	fmt.Println("============================")

	// Demonstrate error handling
	fmt.Println("🔍 Testing error handling...")
	_, err := mcpClient.CallTool(ctx, "non-existent-tool", map[string]interface{}{})
	if err != nil {
		fmt.Printf("✅ Error handling working: %v\n", err)
	}

	// Demonstrate concurrent operations
	fmt.Println("\n🚀 Testing concurrent tool execution...")
	tools, err := mcpClient.ListTools(ctx)
	if err == nil && len(tools) > 0 {
		// Execute the same tool multiple times concurrently
		done := make(chan bool, 3)

		for i := 0; i < 3; i++ {
			go func(id int) {
				result, err := mcpClient.CallTool(ctx, tools[0].Name, map[string]interface{}{
					"query": fmt.Sprintf("concurrent test %d", id),
				})
				if err != nil {
					fmt.Printf("   Concurrent call %d: ❌ %v\n", id, err)
				} else {
					fmt.Printf("   Concurrent call %d: ✅ %d items\n", id, len(result.Content))
				}
				done <- true
			}(i)
		}

		// Wait for all to complete
		for i := 0; i < 3; i++ {
			<-done
		}
	}

	// Demonstrate timeout handling
	fmt.Println("\n⏱️  Testing timeout handling...")
	shortCtx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err = mcpClient.ListTools(shortCtx)
	if err != nil {
		fmt.Printf("✅ Timeout handling working: context deadline exceeded\n")
	}

	fmt.Println("\n✅ Complete protocol demonstration finished!")
}
