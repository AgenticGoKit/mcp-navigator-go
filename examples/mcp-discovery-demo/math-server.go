package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

// Simple MCP server that provides basic math tools
// This demonstrates how to create an MCP server using stdio or TCP transport

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

func main() {
	mode := flag.String("mode", "stdio", "Mode: stdio, tcp, or websocket")
	port := flag.Int("port", 9999, "TCP/WebSocket port")
	flag.Parse()

	switch *mode {
	case "stdio":
		runStdioServer()
	case "tcp":
		runTCPServer(*port)
	case "websocket":
		runWebSocketServer(*port)
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}

// runStdioServer runs MCP server via stdio (stdin/stdout)
func runStdioServer() {
	log.Println("Math MCP Server started (stdio mode)")
	log.Println("Waiting for JSON-RPC requests on stdin...")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Debug: Log incoming request
		log.Printf("📥 Received request: %s", line)

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("❌ Parse error: %v", err)
			sendError("Parse error", -32700, nil)
			continue
		}

		// Debug: Log method being called
		log.Printf("🔧 Handling method: %s", req.Method)

		response := handleRequest(&req)
		if data, err := json.Marshal(response); err == nil {
			log.Printf("📤 Sending response: %s", string(data))
			fmt.Println(string(data))
		}
	}
}

// runTCPServer runs MCP server via TCP
func runTCPServer(port int) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer listener.Close()

	log.Printf("Math MCP Server started (TCP mode on %s)", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

// handleConnection handles a single TCP connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendErrorToConn(conn, "Parse error", -32700, nil)
			continue
		}

		response := handleRequest(&req)
		if data, err := json.Marshal(response); err == nil {
			fmt.Fprintln(conn, string(data))
		}
	}
}

// handleRequest processes a JSON-RPC request
func handleRequest(req *JSONRPCRequest) JSONRPCResponse {
	log.Printf("📨 Received request: %s", req.Method)

	switch req.Method {
	case "initialize":
		return handleInitialize(req)

	case "tools/list":
		return handleToolsList(req)

	case "tools/call":
		return handleToolCall(req)

	case "notifications/initialized":
		// Client is notifying that initialization is complete (MCP protocol)
		// Notifications don't require a response
		log.Printf("✅ Client initialization complete")
		return JSONRPCResponse{JSONRPC: "2.0"} // Minimal response

	default:
		return JSONRPCResponse{
			JSONRPC: "2.0",
			Error: map[string]interface{}{
				"code":    -32601,
				"message": fmt.Sprintf("Method not found: %s", req.Method),
			},
			ID: req.ID,
		}
	}
}

// handleInitialize handles MCP protocol initialization
func handleInitialize(req *JSONRPCRequest) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "math-server",
				"version": "1.0.0",
			},
		},
		ID: req.ID,
	}
}

// handleToolsList returns list of available tools
func handleToolsList(req *JSONRPCRequest) JSONRPCResponse {
	log.Printf("🔧 Handling tools/list request")

	tools := []map[string]interface{}{
		{
			"name":        "add",
			"description": "Adds two numbers together",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "number",
						"description": "First number",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "Second number",
					},
				},
				"required": []string{"a", "b"},
			},
		},
		{
			"name":        "subtract",
			"description": "Subtracts second number from first number",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "number",
						"description": "First number",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "Second number",
					},
				},
				"required": []string{"a", "b"},
			},
		},
		{
			"name":        "multiply",
			"description": "Multiplies two numbers together",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "number",
						"description": "First number",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "Second number",
					},
				},
				"required": []string{"a", "b"},
			},
		},
		{
			"name":        "divide",
			"description": "Divides first number by second number",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "number",
						"description": "Dividend",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "Divisor",
					},
				},
				"required": []string{"a", "b"},
			},
		},
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"tools": tools,
		},
		ID: req.ID,
	}

	log.Printf("📤 Returning %d tools", len(tools))
	return response
}

// handleToolCall executes a tool
func handleToolCall(req *JSONRPCRequest) JSONRPCResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   "Invalid params",
			ID:      req.ID,
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   "Missing tool name",
			ID:      req.ID,
		}
	}

	toolParams, ok := params["arguments"].(map[string]interface{})
	if !ok {
		toolParams = make(map[string]interface{})
	}

	result := executeTool(toolName, toolParams)

	return JSONRPCResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
		ID: req.ID,
	}
}

// executeTool executes a specific math tool
func executeTool(name string, params map[string]interface{}) string {
	a, aOk := getNumber(params, "a")
	b, bOk := getNumber(params, "b")

	if !aOk || !bOk {
		return "Error: Missing required parameters 'a' and 'b'"
	}

	switch name {
	case "add":
		return fmt.Sprintf("%g + %g = %g", a, b, a+b)
	case "subtract":
		return fmt.Sprintf("%g - %g = %g", a, b, a-b)
	case "multiply":
		return fmt.Sprintf("%g × %g = %g", a, b, a*b)
	case "divide":
		if b == 0 {
			return "Error: Cannot divide by zero"
		}
		return fmt.Sprintf("%g ÷ %g = %g", a, b, a/b)
	default:
		return fmt.Sprintf("Unknown tool: %s", name)
	}
}

// getNumber extracts a number from params
func getNumber(params map[string]interface{}, key string) (float64, bool) {
	val, exists := params[key]
	if !exists {
		return 0, false
	}

	switch v := val.(type) {
	case float64:
		return v, true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	case int:
		return float64(v), true
	}

	return 0, false
}

// sendError sends a JSON-RPC error to stdout
func sendError(message string, code int, id interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
		},
		ID: id,
	}
	if data, err := json.Marshal(response); err == nil {
		fmt.Println(string(data))
	}
}

// sendErrorToConn sends a JSON-RPC error to a connection
func sendErrorToConn(conn net.Conn, message string, code int, id interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
		},
		ID: id,
	}
	if data, err := json.Marshal(response); err == nil {
		fmt.Fprintln(conn, string(data))
	}
}

// runWebSocketServer runs MCP server via WebSocket
func runWebSocketServer(port int) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for testing
		},
	}

	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade WebSocket: %v", err)
			return
		}
		defer conn.Close()

		log.Printf("🔗 WebSocket client connected from %s", conn.RemoteAddr())

		// Handle messages from WebSocket client
		for {
			var req JSONRPCRequest
			err := conn.ReadJSON(&req)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("❌ WebSocket error: %v", err)
				}
				break
			}

			log.Printf("📥 Received WebSocket request: %s", req.Method)

			response := handleRequest(&req)
			if err := conn.WriteJSON(response); err != nil {
				log.Printf("❌ Failed to send WebSocket response: %v", err)
				break
			}
		}

		log.Printf("🔌 WebSocket client disconnected")
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Math MCP Server started (WebSocket mode on ws://localhost:%d/mcp)", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
