//go:build example
// +build example

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kunalkushwaha/mcp-navigator-go/pkg/discovery"
)

// Example demonstrates MCP server discovery functionality
func main() {
	fmt.Println("🔍 MCP Server Discovery Example")
	fmt.Println("===============================")

	// 1. Discover all servers
	if err := exampleDiscoverAll(); err != nil {
		log.Printf("Discover all example failed: %v", err)
	}

	fmt.Println()

	// 2. Discover TCP servers only
	if err := exampleDiscoverTCP(); err != nil {
		log.Printf("Discover TCP example failed: %v", err)
	}

	fmt.Println()

	// 3. Discover HTTP servers only
	if err := exampleDiscoverHTTP(); err != nil {
		log.Printf("Discover HTTP example failed: %v", err)
	}

	fmt.Println()

	// 4. Discover Docker servers only
	if err := exampleDiscoverDocker(); err != nil {
		log.Printf("Discover Docker example failed: %v", err)
	}
}

// exampleDiscoverAll demonstrates comprehensive server discovery
func exampleDiscoverAll() error {
	fmt.Println("1. Discover All Servers")
	fmt.Println("----------------------")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fmt.Println("🔍 Discovering all MCP servers...")

	discoverer := discovery.NewDiscoverer()
	servers, err := discoverer.DiscoverAll(ctx)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}

	fmt.Printf("✅ Found %d server(s):\n", len(servers))
	for i, server := range servers {
		fmt.Printf("   %d. %s (%s)\n", i+1, server.Name, server.Type)
		fmt.Printf("      Address: %s\n", server.Address)
		fmt.Printf("      Description: %s\n", server.Description)
		fmt.Println()
	}

	return nil
}

// exampleDiscoverTCP demonstrates TCP-specific server discovery
func exampleDiscoverTCP() error {
	fmt.Println("2. Discover TCP Servers")
	fmt.Println("----------------------")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("🔍 Scanning for TCP MCP servers...")

	discoverer := discovery.NewDiscoverer()

	// Scan common MCP ports
	servers, err := discoverer.DiscoverTCP(ctx, "localhost", 8810, 8820)
	if err != nil {
		return fmt.Errorf("TCP discovery failed: %w", err)
	}

	fmt.Printf("✅ Found %d TCP server(s):\n", len(servers))
	for i, server := range servers {
		fmt.Printf("   %d. %s\n", i+1, server.Name)
		fmt.Printf("      Address: %s\n", server.Address)
		fmt.Printf("      Port: %s\n", extractPort(server.Address))
	}

	return nil
}

// exampleDiscoverHTTP demonstrates HTTP-specific server discovery
func exampleDiscoverHTTP() error {
	fmt.Println("3. Discover HTTP Servers")
	fmt.Println("-----------------------")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("🔍 Scanning for HTTP MCP servers...")

	discoverer := discovery.NewDiscoverer()
	servers, err := discoverer.DiscoverHTTP(ctx, "localhost")
	if err != nil {
		return fmt.Errorf("HTTP discovery failed: %w", err)
	}

	fmt.Printf("✅ Found %d HTTP server(s):\n", len(servers))
	for i, server := range servers {
		fmt.Printf("   %d. %s\n", i+1, server.Name)
		fmt.Printf("      Address: %s\n", server.Address)
		fmt.Printf("      Type: %s\n", getHTTPType(server.Name))
	}

	return nil
}

// exampleDiscoverDocker demonstrates Docker-specific server discovery
func exampleDiscoverDocker() error {
	fmt.Println("4. Discover Docker Servers")
	fmt.Println("-------------------------")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("🔍 Scanning for Docker MCP servers...")

	discoverer := discovery.NewDiscoverer()
	servers, err := discoverer.DiscoverDocker(ctx)
	if err != nil {
		return fmt.Errorf("Docker discovery failed: %w", err)
	}

	fmt.Printf("✅ Found %d Docker server(s):\n", len(servers))
	for i, server := range servers {
		fmt.Printf("   %d. %s\n", i+1, server.Name)
		fmt.Printf("      Address: %s\n", server.Address)
		fmt.Printf("      Description: %s\n", server.Description)
	}

	if len(servers) == 0 {
		fmt.Println("   Note: Docker MCP servers require Docker to be running")
		fmt.Println("   and MCP servers configured in Docker containers.")
	}

	return nil
}

// Helper function to extract port from address
func extractPort(address string) string {
	// Simple extraction for display purposes
	if len(address) > 0 {
		parts := address[len(address)-4:]
		return parts
	}
	return "unknown"
}

// Helper function to determine HTTP server type
func getHTTPType(name string) string {
	if contains(name, "SSE") {
		return "Server-Sent Events"
	} else if contains(name, "Streaming") {
		return "HTTP Streaming"
	}
	return "HTTP"
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && findInString(s, substr))
}

// Simple substring search
func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
