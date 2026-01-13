#!/usr/bin/env bash

# MCP Discovery Demo - Quick Start Script
# This script sets up and runs the MCP discovery example

set -e

DEMO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DEMO_DIR"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

show_usage() {
    echo "MCP Discovery Demo - Quick Start"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  build              Build both demo and math-server"
    echo "  run-stdio          Run demo with stdio transport"
    echo "  run-tcp            Run demo with TCP transport (requires manual server start)"
    echo "  run-discover       Run demo with auto-discovery"
    echo "  start-tcp-server   Start TCP math server (port 9999)"
    echo "  test-stdio         Test stdio server directly"
    echo "  test-tcp           Test TCP server directly"
    echo "  clean              Remove built binaries"
    echo "  help               Show this message"
    echo ""
}

build() {
    print_header "Building MCP Demo"
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        return 1
    fi
    
    print_success "Building mcp-demo..."
    go build -o mcp-demo .
    
    print_success "Building math-server..."
    go build -o math-server math-server.go
    
    print_success "Build complete!"
}

check_ollama() {
    if ! command -v ollama &> /dev/null; then
        print_error "Ollama is not installed or not in PATH"
        echo ""
        echo "Please install Ollama from: https://ollama.ai"
        echo "Then pull the llama2 model:"
        echo "  ollama pull llama2"
        return 1
    fi
    
    # Check if llama2 model is available
    if ! ollama list | grep -q llama2; then
        print_error "llama2 model not found"
        echo "Please pull it with: ollama pull llama2"
        return 1
    fi
    
    print_success "Ollama and llama2 model are ready"
    return 0
}

run_stdio() {
    print_header "Running MCP Demo - Stdio Transport"
    
    if [ ! -f "mcp-demo" ] || [ ! -f "math-server" ]; then
        print_error "Binaries not found. Running build first..."
        build
    fi
    
    if ! check_ollama; then
        return 1
    fi
    
    print_success "Starting agent with stdio MCP server..."
    echo ""
    
    ./mcp-demo -mode stdio
}

run_tcp() {
    print_header "Running MCP Demo - TCP Transport"
    
    if [ ! -f "mcp-demo" ]; then
        print_error "mcp-demo binary not found. Running build first..."
        build
    fi
    
    if ! check_ollama; then
        return 1
    fi
    
    echo "Checking if TCP server is running on localhost:9999..."
    
    if ! nc -z localhost 9999 2>/dev/null; then
        print_error "TCP server is not running on localhost:9999"
        echo ""
        echo "Please start the server in another terminal:"
        echo "  $0 start-tcp-server"
        return 1
    fi
    
    print_success "TCP server is running"
    print_success "Starting agent..."
    echo ""
    
    ./mcp-demo -mode tcp
}

run_discover() {
    print_header "Running MCP Demo - Auto-Discovery"
    
    if [ ! -f "mcp-demo" ]; then
        print_error "mcp-demo binary not found. Running build first..."
        build
    fi
    
    if ! check_ollama; then
        return 1
    fi
    
    print_success "Starting agent with auto-discovery..."
    echo ""
    
    ./mcp-demo -mode discover
}

start_tcp_server() {
    print_header "Starting TCP Math Server"
    
    if [ ! -f "math-server" ]; then
        print_error "math-server binary not found. Running build first..."
        build
    fi
    
    PORT="${1:-9999}"
    
    echo "Starting math-server on port $PORT..."
    echo "Press Ctrl+C to stop"
    echo ""
    
    ./math-server -mode tcp -port "$PORT"
}

test_stdio() {
    print_header "Testing Stdio Server"
    
    if [ ! -f "math-server" ]; then
        print_error "math-server binary not found. Running build first..."
        build
    fi
    
    print_success "Testing tools/list request..."
    echo ""
    
    echo '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}' | \
        ./math-server
    
    echo ""
    echo ""
    print_success "Testing tools/call request (add 5 + 3)..."
    echo ""
    
    echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"add","arguments":{"a":5,"b":3}},"id":2}' | \
        ./math-server
    
    echo ""
}

test_tcp() {
    print_header "Testing TCP Server"
    
    if ! nc -z localhost 9999 2>/dev/null; then
        print_error "TCP server is not running on localhost:9999"
        echo "Start it with: $0 start-tcp-server"
        return 1
    fi
    
    print_success "TCP server is running on localhost:9999"
    echo ""
    print_success "Testing with telnet..."
    echo "Send: {\"jsonrpc\":\"2.0\",\"method\":\"tools/list\",\"params\":{},\"id\":1}"
    echo "Then press Enter twice"
    echo ""
    
    telnet localhost 9999
}

clean() {
    print_header "Cleaning up"
    
    rm -f mcp-demo math-server
    
    print_success "Cleaned up binaries"
}

# Main script
case "${1:-help}" in
    build)
        build
        ;;
    run-stdio)
        run_stdio
        ;;
    run-tcp)
        run_tcp
        ;;
    run-discover)
        run_discover
        ;;
    start-tcp-server)
        start_tcp_server "$2"
        ;;
    test-stdio)
        test_stdio
        ;;
    test-tcp)
        test_tcp
        ;;
    clean)
        clean
        ;;
    help|--help|-h|"")
        show_usage
        ;;
    *)
        print_error "Unknown command: $1"
        show_usage
        exit 1
        ;;
esac
