package transport

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kunalkushwaha/mcp-navigator-go/pkg/mcp"
)

// SSETransport implements Transport for HTTP/SSE connections
type SSETransport struct {
	baseURL       string
	endpoint      string
	client        *http.Client
	sessionID     string
	sessionURL    string // The actual message endpoint
	connected     bool
	mu            sync.RWMutex
	timeout       time.Duration
	initialized   bool
	lastResponse  *mcp.Message
	sseConnection *http.Response // Keep SSE connection alive
}

// NewSSETransport creates a new SSE transport for SSE-based MCP servers
func NewSSETransport(baseURL, endpoint string) *SSETransport {
	return &SSETransport{
		baseURL:  baseURL,
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 30 * time.Second,
	}
}

// Connect establishes HTTP connection and handles initialization
func (h *SSETransport) Connect(ctx context.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connected {
		return nil
	}

	h.connected = true
	return nil
}

// Close closes the HTTP connection
func (h *SSETransport) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Close the SSE connection if it exists
	if h.sseConnection != nil {
		h.sseConnection.Body.Close()
		h.sseConnection = nil
	}

	h.connected = false
	h.sessionID = ""
	h.sessionURL = ""
	h.initialized = false
	return nil
}

// Send sends a message over HTTP
func (h *SSETransport) Send(message *mcp.Message) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.connected {
		return fmt.Errorf("transport not connected")
	}

	// Special handling for initialize request
	if message.Method == "initialize" && !h.initialized {
		return h.sendInitializeRequest(message)
	}

	// For notifications (no ID), use different handling
	if message.ID == nil {
		return h.sendNotification(message)
	}

	// For other requests, use session-based request
	return h.sendSessionRequest(message)
}

// sendInitializeRequest handles the initial request that establishes session
func (h *SSETransport) sendInitializeRequest(message *mcp.Message) error {
	// First, establish SSE connection to get session endpoint
	sseURL := h.baseURL + h.endpoint

	req, err := http.NewRequest("GET", sseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create SSE request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to establish SSE connection: %w", err)
	}
	// Don't defer close here - we'll close after getting session URL

	// Parse SSE stream to get session endpoint
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data: ") {
			sessionEndpoint := strings.TrimSpace(line[6:]) // Remove "data: " prefix
			h.sessionURL = h.baseURL + sessionEndpoint
			// Keep the SSE connection alive by storing it
			h.sseConnection = resp
			// Immediately send the initialize request while session is valid
			return h.sendMessageToSession(message)
		}
	}

	resp.Body.Close()
	return fmt.Errorf("failed to get session endpoint from SSE")
} // sendMessageToSession sends a message to the established session endpoint
func (h *SSETransport) sendMessageToSession(message *mcp.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", h.sessionURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	// Read response body and store for Receive()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Handle empty response (might be normal for SSE mode)
	if len(body) == 0 {
		// For SSE mode, the response might come through the SSE stream
		// We'll try to read it in the Receive() method
		h.lastResponse = nil
		h.initialized = true
		return nil
	}

	var response mcp.Message
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	h.lastResponse = &response
	h.initialized = true
	return nil
}

// sendNotification sends notification messages (no response expected)
func (h *SSETransport) sendNotification(message *mcp.Message) error {
	// Use the session URL for notifications too
	if h.sessionURL == "" {
		return fmt.Errorf("session not established")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", h.sessionURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	// For notifications, we don't expect a meaningful response
	// Just check that we got a successful HTTP status
	if resp.StatusCode >= 400 {
		return fmt.Errorf("notification failed with status: %d", resp.StatusCode)
	}

	return nil
}

// sendSessionRequest sends requests using session URL after initialization
func (h *SSETransport) sendSessionRequest(message *mcp.Message) error {
	// Use the session URL instead of the base endpoint
	if h.sessionURL == "" {
		return fmt.Errorf("session not established")
	}

	return h.sendMessageToSession(message)
}

// Add fields to store the last response
// Receive returns the stored response from the last Send() call
func (h *SSETransport) Receive() (*mcp.Message, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.connected {
		return nil, fmt.Errorf("transport not connected")
	}

	// Return the stored response from the last Send() call
	if h.lastResponse != nil {
		response := h.lastResponse
		h.lastResponse = nil // Clear after returning
		return response, nil
	}

	// If no stored response and we have an SSE connection, try reading from it
	if h.sseConnection != nil {
		scanner := bufio.NewScanner(h.sseConnection.Body)

		// Read multiple lines to find the data
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "data: ") {
				dataStr := strings.TrimSpace(line[6:]) // Remove "data: " prefix

				var response mcp.Message
				if err := json.Unmarshal([]byte(dataStr), &response); err == nil {
					return &response, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no response available")
}

// GetReader returns nil for HTTP transport as it's request-response based
func (h *SSETransport) GetReader() io.Reader {
	return nil
}

// GetWriter returns nil for HTTP transport as it's request-response based
func (h *SSETransport) GetWriter() io.Writer {
	return nil
}

// IsConnected returns connection status
func (h *SSETransport) IsConnected() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.connected
}

// SetTimeout sets the HTTP client timeout
func (h *SSETransport) SetTimeout(timeout time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.timeout = timeout
	h.client.Timeout = timeout
}

// GetSessionID returns the current session ID
func (h *SSETransport) GetSessionID() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.sessionID
}
