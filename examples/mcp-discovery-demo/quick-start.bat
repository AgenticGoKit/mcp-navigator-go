@echo off
REM MCP Discovery Demo - Quick Start Script (Windows Batch)
REM This script sets up and runs the MCP discovery example on Windows

setlocal enabledelayedexpansion

if "%1"=="" (
    call :show_usage
    exit /b 0
)

if /i "%1"=="build" (
    call :build
) else if /i "%1"=="run-stdio" (
    call :run_stdio
) else if /i "%1"=="run-tcp" (
    call :run_tcp
) else if /i "%1"=="run-discover" (
    call :run_discover
) else if /i "%1"=="start-tcp-server" (
    call :start_tcp_server %2
) else if /i "%1"=="test-stdio" (
    call :test_stdio
) else if /i "%1"=="test-tcp" (
    call :test_tcp
) else if /i "%1"=="clean" (
    call :clean
) else if /i "%1"=="help" (
    call :show_usage
) else (
    echo.
    echo [!] Unknown command: %1
    echo.
    call :show_usage
    exit /b 1
)

exit /b 0

REM ============================================================================
REM FUNCTIONS
REM ============================================================================

:show_usage
echo.
echo MCP Discovery Demo - Quick Start
echo.
echo Usage: quick-start.bat [command]
echo.
echo Commands:
echo   build              Build both demo and math-server
echo   run-stdio          Run demo with stdio transport
echo   run-tcp            Run demo with TCP transport ^(requires manual server start^)
echo   run-discover       Run demo with auto-discovery
echo   start-tcp-server   Start TCP math server ^(port 9999^)
echo   test-stdio         Test stdio server directly
echo   test-tcp           Test TCP server directly
echo   clean              Remove built binaries
echo   help               Show this message
echo.
exit /b 0

:build
echo.
echo === Building MCP Demo ===
echo.

where go >nul 2>&1
if errorlevel 1 (
    echo [!] Go is not installed or not in PATH
    exit /b 1
)

echo Building mcp-demo.exe...
call go build -o mcp-demo.exe .
if errorlevel 1 (
    echo [!] Failed to build mcp-demo
    exit /b 1
)
echo [+] mcp-demo.exe built

echo Building math-server.exe...
call go build -o math-server.exe math-server.go
if errorlevel 1 (
    echo [!] Failed to build math-server
    exit /b 1
)
echo [+] math-server.exe built

echo.
echo [+] Build complete!
echo.
exit /b 0

:check_ollama
where ollama >nul 2>&1
if errorlevel 1 (
    echo [!] Ollama is not installed or not in PATH
    echo.
    echo Please install Ollama from: https://ollama.ai
    echo Then pull the llama2 model:
    echo   ollama pull llama2
    exit /b 1
)

echo [+] Ollama is ready
exit /b 0

:run_stdio
echo.
echo === Running MCP Demo - Stdio Transport ===
echo.

if not exist "mcp-demo.exe" (
    if not exist "math-server.exe" (
        echo [!] Binaries not found. Running build first...
        call :build
    )
)

call :check_ollama
if errorlevel 1 exit /b 1

echo [+] Starting agent with stdio MCP server...
echo.

mcp-demo.exe -mode stdio
exit /b 0

:run_tcp
echo.
echo === Running MCP Demo - TCP Transport ===
echo.

if not exist "mcp-demo.exe" (
    echo [!] mcp-demo.exe not found. Running build first...
    call :build
)

call :check_ollama
if errorlevel 1 exit /b 1

echo Checking if TCP server is running on localhost:9999...

REM Try to connect to TCP server
powershell -Command "$tcpClient = New-Object System.Net.Sockets.TcpClient; $result = $tcpClient.BeginConnect('localhost', 9999, $null, $null); $success = $result.AsyncWaitHandle.WaitOne(1000, $false); $tcpClient.Close(); exit 0"
if errorlevel 1 (
    echo [!] TCP server is not running on localhost:9999
    echo.
    echo Please start the server in another terminal:
    echo   quick-start.bat start-tcp-server
    exit /b 1
)

echo [+] TCP server is running
echo [+] Starting agent...
echo.

mcp-demo.exe -mode tcp
exit /b 0

:run_discover
echo.
echo === Running MCP Demo - Auto-Discovery ===
echo.

if not exist "mcp-demo.exe" (
    echo [!] mcp-demo.exe not found. Running build first...
    call :build
)

call :check_ollama
if errorlevel 1 exit /b 1

echo [+] Starting agent with auto-discovery...
echo.

mcp-demo.exe -mode discover
exit /b 0

:start_tcp_server
echo.
echo === Starting TCP Math Server ===
echo.

if not exist "math-server.exe" (
    echo [!] math-server.exe not found. Running build first...
    call :build
)

set PORT=9999
if not "%2"=="" set PORT=%2

echo Starting math-server.exe on port %PORT%...
echo Press Ctrl+C to stop
echo.

math-server.exe -mode tcp -port %PORT%
exit /b 0

:test_stdio
echo.
echo === Testing Stdio Server ===
echo.

if not exist "math-server.exe" (
    echo [!] math-server.exe not found. Running build first...
    call :build
)

echo [+] Testing tools/list request...
echo.

echo {"jsonrpc":"2.0","method":"tools/list","params":{},"id":1} | math-server.exe

echo.
echo.
echo [+] Testing tools/call request ^(add 5 + 3^)...
echo.

echo {"jsonrpc":"2.0","method":"tools/call","params":{"name":"add","arguments":{"a":5,"b":3}},"id":2} | math-server.exe

echo.
exit /b 0

:test_tcp
echo.
echo === Testing TCP Server ===
echo.

REM Check if server is running
powershell -Command "$tcpClient = New-Object System.Net.Sockets.TcpClient; $result = $tcpClient.BeginConnect('localhost', 9999, $null, $null); $success = $result.AsyncWaitHandle.WaitOne(1000, $false); $tcpClient.Close(); exit 0"
if errorlevel 1 (
    echo [!] TCP server is not running on localhost:9999
    echo Please start it with: quick-start.bat start-tcp-server
    exit /b 1
)

echo [+] TCP server is running on localhost:9999
echo.
echo [+] Testing with PowerShell connection...
echo.

powershell -Command ^
    "$tcpClient = New-Object System.Net.Sockets.TcpClient('localhost', 9999); " ^ 
    "$stream = $tcpClient.GetStream(); " ^
    "$writer = New-Object System.IO.StreamWriter($stream); " ^
    "$reader = New-Object System.IO.StreamReader($stream); " ^
    "$request = '{\"jsonrpc\":\"2.0\",\"method\":\"tools/list\",\"params\":{},\"id\":1}'; " ^
    "$writer.WriteLine($request); " ^
    "$writer.Flush(); " ^
    "$response = $reader.ReadLine(); " ^
    "Write-Host \"Response: $response\" -ForegroundColor Green; " ^
    "$writer.Close(); " ^
    "$reader.Close(); " ^
    "$tcpClient.Close();"

echo.
exit /b 0

:clean
echo.
echo === Cleaning up ===
echo.

if exist "mcp-demo.exe" del mcp-demo.exe
if exist "math-server.exe" del math-server.exe

echo [+] Cleaned up binaries
echo.
exit /b 0

endlocal
