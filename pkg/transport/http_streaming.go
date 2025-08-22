package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/kunalkushwaha/mcp-navigator-go/pkg/mcp"
)

// StreamingHTTPTransport implements Transport for streaming HTTP MCP servers
type StreamingHTTPTransport struct {
	baseURL      string
	endpoint     string
	client       *http.Client
	sessionID    string
	connected    bool
	mu           sync.RWMutex
	timeout      time.Duration
	initialized  bool
	lastResponse *mcp.Message
}

// NewStreamingHTTPTransport creates a new streaming HTTP transport
func NewStreamingHTTPTransport(baseURL, endpoint string) *StreamingHTTPTransport {
	return &StreamingHTTPTransport{
		baseURL:  baseURL,
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 30 * time.Second,
	}
}

// Connect establishes HTTP connection
func (h *StreamingHTTPTransport) Connect(ctx context.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connected {
		return nil
	}

	h.connected = true
	return nil
}

// Close closes the HTTP connection
func (h *StreamingHTTPTransport) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.connected = false
	h.sessionID = ""
	h.initialized = false
	return nil
}

// Send sends a message over HTTP
func (h *StreamingHTTPTransport) Send(message *mcp.Message) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.connected {
		return fmt.Errorf("transport not connected")
	}

	url := h.baseURL + h.endpoint

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add session ID to subsequent requests
	if h.sessionID != "" {
		req.Header.Set("Mcp-Session-Id", h.sessionID)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Extract session ID from response headers
	sessionID := resp.Header.Get("Mcp-Session-Id")
	if sessionID != "" {
		h.sessionID = sessionID
	}

	// Read response body and store for Receive()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Handle empty responses (for notifications)
	if len(body) == 0 {
		return nil
	}

	var response mcp.Message
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	h.lastResponse = &response
	return nil
}

// Receive receives a message from HTTP
func (h *StreamingHTTPTransport) Receive() (*mcp.Message, error) {
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

	return nil, fmt.Errorf("no response available")
}

// GetReader returns nil for HTTP transport as it's request-response based
func (h *StreamingHTTPTransport) GetReader() io.Reader {
	return nil
}

// GetWriter returns nil for HTTP transport as it's request-response based
func (h *StreamingHTTPTransport) GetWriter() io.Writer {
	return nil
}

// IsConnected returns connection status
func (h *StreamingHTTPTransport) IsConnected() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.connected
}

// SetTimeout sets the HTTP client timeout
func (h *StreamingHTTPTransport) SetTimeout(timeout time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.timeout = timeout
	h.client.Timeout = timeout
}
