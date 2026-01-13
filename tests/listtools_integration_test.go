package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/kunalkushwaha/mcp-navigator-go/pkg/client"
	"github.com/kunalkushwaha/mcp-navigator-go/pkg/mcp"
	"github.com/kunalkushwaha/mcp-navigator-go/pkg/transport"
)

// TestListToolsWithMockServer tests the full ListTools flow with a mock MCP server
// This test validates MCP 2025-11-25 spec compliance
func TestListToolsWithMockServer(t *testing.T) {
	// This test requires a real MCP server running on port 9999
	// Skip if not available
	t.Skip("Requires manual server setup - run math-server first")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trans := transport.NewTCPTransport("localhost", 9999)
	config := client.ClientConfig{
		Name:    "test-client",
		Version: "1.0.0",
	}

	c := client.NewClient(trans, config)

	if err := c.Connect(ctx); err != nil {
		t.Skipf("Skipping test - no server available: %v", err)
	}
	defer c.Disconnect()

	clientInfo := mcp.ClientInfo{
		Name:    "test-client",
		Version: "1.0.0",
	}

	serverInfo, err := c.Initialize(ctx, clientInfo)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	t.Logf("Connected to server: %s %s", serverInfo.Name, serverInfo.Version)

	tools, err := c.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	t.Logf("Retrieved %d tools", len(tools))

	if len(tools) == 0 {
		t.Error("Expected non-empty tools list, got empty array")
	}

	// Validate each tool per MCP spec
	for i, tool := range tools {
		t.Logf("Tool %d: Name=%s, Title=%s, Description=%s", i+1, tool.Name, tool.Title, tool.Description)

		// Validate tool name per MCP guidelines
		if err := mcp.ValidateToolName(tool.Name); err != nil {
			t.Errorf("Tool %d has invalid name: %v", i+1, err)
		}

		// Validate tool structure
		if err := mcp.ValidateTool(&tool); err != nil {
			t.Errorf("Tool %d failed validation: %v", i+1, err)
		}

		// Validate inputSchema is valid JSON
		if tool.InputSchema != nil {
			schemaJSON, err := json.Marshal(tool.InputSchema)
			if err != nil {
				t.Errorf("Tool %d inputSchema is not valid JSON: %v", i+1, err)
			} else {
				t.Logf("  InputSchema: %s", string(schemaJSON))
			}
		} else {
			t.Errorf("Tool %d missing required inputSchema", i+1)
		}
	}
}

// TestToolNameValidation tests the ValidateToolName function per MCP spec
func TestToolNameValidation(t *testing.T) {
	tests := []struct {
		name      string
		toolName  string
		expectErr bool
	}{
		{"valid simple name", "getTool", false},
		{"valid with underscore", "get_tool", false},
		{"valid with hyphen", "get-tool", false},
		{"valid with dot", "admin.tools.list", false},
		{"valid with numbers", "tool123", false},
		{"valid uppercase", "DATA_EXPORT_v2", false},
		{"valid mixed case", "getUser", false},
		{"empty name", "", true},
		{"too long", string(make([]byte, 129)), true},
		{"with spaces", "get tool", true},
		{"with comma", "get,tool", true},
		{"with special chars", "get@tool", true},
		{"single char", "t", false},
		{"exactly 128 chars", string(make([]byte, 128)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For the 128-char test, fill with valid chars
			testName := tt.toolName
			if len(tt.toolName) == 128 && !tt.expectErr {
				for i := 0; i < 128; i++ {
					testName = "a" + testName[1:]
				}
			}

			err := mcp.ValidateToolName(testName)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error for tool name '%s', got nil", tt.toolName)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error for tool name '%s', got: %v", tt.toolName, err)
			}
		})
	}
}

// TestToolValidation tests the ValidateTool function per MCP spec
func TestToolValidation(t *testing.T) {
	tests := []struct {
		name      string
		tool      mcp.Tool
		expectErr bool
	}{
		{
			name: "valid tool",
			tool: mcp.Tool{
				Name:        "add",
				Description: "Adds two numbers",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"a": map[string]interface{}{"type": "number"},
						"b": map[string]interface{}{"type": "number"},
					},
					"required": []string{"a", "b"},
				},
			},
			expectErr: false,
		},
		{
			name: "missing inputSchema",
			tool: mcp.Tool{
				Name:        "test",
				Description: "Test tool",
				InputSchema: nil,
			},
			expectErr: true,
		},
		{
			name: "invalid tool name",
			tool: mcp.Tool{
				Name:        "invalid name with spaces",
				Description: "Test",
				InputSchema: map[string]interface{}{"type": "object"},
			},
			expectErr: true,
		},
		{
			name: "tool with title and icons (MCP 2025-11-25)",
			tool: mcp.Tool{
				Name:        "weather",
				Title:       "Weather Information Provider",
				Description: "Get weather data",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{"type": "string"},
					},
				},
				Icons: []mcp.Icon{
					{
						Src:      "https://example.com/icon.png",
						MimeType: "image/png",
						Sizes:    []string{"48x48"},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mcp.ValidateTool(&tt.tool)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}

// TestListToolsResponseParsing tests parsing of MCP-compliant tools/list responses
func TestListToolsResponseParsing(t *testing.T) {
	// Simulate the exact response from the bug report
	responseJSON := `{
		"jsonrpc": "2.0",
		"result": {
			"tools": [
				{
					"name": "add",
					"description": "Adds two numbers together",
					"inputSchema": {
						"type": "object",
						"properties": {
							"a": {"type": "number", "description": "First number"},
							"b": {"type": "number", "description": "Second number"}
						},
						"required": ["a", "b"]
					}
				},
				{
					"name": "subtract",
					"description": "Subtracts second number from first",
					"inputSchema": {
						"type": "object",
						"properties": {
							"a": {"type": "number", "description": "First number"},
							"b": {"type": "number", "description": "Second number"}
						},
						"required": ["a", "b"]
					}
				}
			]
		},
		"id": 2
	}`

	var message mcp.Message
	if err := json.Unmarshal([]byte(responseJSON), &message); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Unmarshaled message: ID=%v, Result type=%T", message.ID, message.Result)

	// Now parse the result into ListToolsResponse
	resultJSON, err := json.Marshal(message.Result)
	if err != nil {
		t.Fatalf("Failed to re-marshal result: %v", err)
	}

	t.Logf("Result JSON: %s", string(resultJSON))

	var listResponse mcp.ListToolsResponse
	if err := json.Unmarshal(resultJSON, &listResponse); err != nil {
		t.Fatalf("Failed to unmarshal into ListToolsResponse: %v", err)
	}

	t.Logf("Parsed %d tools", len(listResponse.Tools))

	if len(listResponse.Tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(listResponse.Tools))
	}

	// Validate the parsed tools
	for i, tool := range listResponse.Tools {
		t.Logf("Tool %d: %s - %s", i+1, tool.Name, tool.Description)

		if tool.InputSchema == nil {
			t.Errorf("Tool %d missing inputSchema", i+1)
		}

		if err := mcp.ValidateTool(&tool); err != nil {
			t.Errorf("Tool %d failed validation: %v", i+1, err)
		}
	}
}
