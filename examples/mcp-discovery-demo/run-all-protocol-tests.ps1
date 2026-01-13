# Protocol Testing Script
# Tests the mcp-navigator-go library with all supported transports

param(
    [switch]$Debug = $false,
    [switch]$SkipBuild = $false,
    [int]$Port = 9999
)

Write-Host "=== MCP Navigator - Protocol Testing Suite ===" -ForegroundColor Cyan
Write-Host "Testing all transport protocols" -ForegroundColor Gray
Write-Host ""

# Colors for output
$successColor = "Green"
$errorColor = "Red"
$infoColor = "Cyan"
$warningColor = "Yellow"

# Track results
$results = @()

# Step 1: Build executables if needed
if (-not $SkipBuild) {
    Write-Host "Building executables..." -ForegroundColor $infoColor
    
    if (-not (Test-Path "math-server.exe")) {
        Write-Host "  📦 Building math-server..." -ForegroundColor Gray
        go build -o math-server.exe math-server.go
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✅ math-server built successfully" -ForegroundColor $successColor
        } else {
            Write-Host "  ❌ Failed to build math-server" -ForegroundColor $errorColor
            exit 1
        }
    }
    
    if (-not (Test-Path "test-client-tcp.exe")) {
        Write-Host "  📦 Building test-client-tcp..." -ForegroundColor Gray
        go build -o test-client-tcp.exe test-client-tcp.go
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✅ test-client-tcp built successfully" -ForegroundColor $successColor
        } else {
            Write-Host "  ❌ Failed to build test-client-tcp" -ForegroundColor $errorColor
            exit 1
        }
    }
    
    if (-not (Test-Path "test-client-stdio.exe")) {
        Write-Host "  📦 Building test-client-stdio..." -ForegroundColor Gray
        go build -o test-client-stdio.exe test-client-stdio.go
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✅ test-client-stdio built successfully" -ForegroundColor $successColor
        } else {
            Write-Host "  ❌ Failed to build test-client-stdio" -ForegroundColor $errorColor
            exit 1
        }
    }
    
    if (-not (Test-Path "test-client-websocket.exe")) {
        Write-Host "  📦 Building test-client-websocket..." -ForegroundColor Gray
        go build -o test-client-websocket.exe test-client-websocket.go
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✅ test-client-websocket built successfully" -ForegroundColor $successColor
        } else {
            Write-Host "  ❌ Failed to build test-client-websocket" -ForegroundColor $errorColor
            exit 1
        }
    }
    
    Write-Host ""
}

# Step 2: Kill any existing servers on the port
Write-Host "Cleaning up existing processes..." -ForegroundColor $infoColor
$existingProcess = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
if ($existingProcess) {
    Write-Host "  🔄 Stopping process on port $Port..." -ForegroundColor Gray
    Stop-Process -Id (Get-Process | Where-Object {$_.MainWindowHandle -eq (Get-NetTCPConnection -LocalPort $Port).OwningProcess}) -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 1
}
Write-Host ""

# Step 3: Test TCP Protocol
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor $infoColor
Write-Host "║ Test 1: TCP Protocol                   ║" -ForegroundColor $infoColor
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor $infoColor

Write-Host "Starting math-server in TCP mode on port $Port..." -ForegroundColor Gray
$serverProcess = Start-Process -FilePath ".\math-server.exe" -ArgumentList @("-mode", "tcp", "-port", $Port) -PassThru -NoNewWindow
Start-Sleep -Seconds 2

Write-Host "Running TCP test client..." -ForegroundColor Gray
if ($Debug) {
    & ".\test-client-tcp.exe" -port $Port -debug
} else {
    & ".\test-client-tcp.exe" -port $Port
}

$tcpResult = $LASTEXITCODE
if ($tcpResult -eq 0) {
    Write-Host "✅ TCP Test PASSED" -ForegroundColor $successColor
    $results += @{"Protocol" = "TCP"; "Status" = "PASSED"}
} else {
    Write-Host "❌ TCP Test FAILED" -ForegroundColor $errorColor
    $results += @{"Protocol" = "TCP"; "Status" = "FAILED"}
}

# Cleanup
Stop-Process -Id $serverProcess.Id -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 1
Write-Host ""

# Step 4: Test STDIO Protocol
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor $infoColor
Write-Host "║ Test 2: STDIO Protocol                 ║" -ForegroundColor $infoColor
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor $infoColor

Write-Host "Running STDIO test client (launches server automatically)..." -ForegroundColor Gray
if ($Debug) {
    & ".\test-client-stdio.exe" -debug
} else {
    & ".\test-client-stdio.exe"
}

$stdioResult = $LASTEXITCODE
if ($stdioResult -eq 0) {
    Write-Host "✅ STDIO Test PASSED" -ForegroundColor $successColor
    $results += @{"Protocol" = "STDIO"; "Status" = "PASSED"}
} else {
    Write-Host "❌ STDIO Test FAILED" -ForegroundColor $errorColor
    $results += @{"Protocol" = "STDIO"; "Status" = "FAILED"}
}

Write-Host ""

# Step 5: Test WebSocket Protocol
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor $infoColor
Write-Host "║ Test 3: WebSocket Protocol             ║" -ForegroundColor $infoColor
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor $infoColor

Write-Host "Starting math-server in WebSocket mode on port $Port..." -ForegroundColor Gray
$serverProcess = Start-Process -FilePath ".\math-server.exe" -ArgumentList @("-mode", "websocket", "-port", $Port) -PassThru -NoNewWindow
Start-Sleep -Seconds 2

Write-Host "Running WebSocket test client..." -ForegroundColor Gray
if ($Debug) {
    & ".\test-client-websocket.exe" -port $Port -debug
} else {
    & ".\test-client-websocket.exe" -port $Port
}

$wsResult = $LASTEXITCODE
if ($wsResult -eq 0) {
    Write-Host "✅ WebSocket Test PASSED" -ForegroundColor $successColor
    $results += @{"Protocol" = "WebSocket"; "Status" = "PASSED"}
} else {
    Write-Host "❌ WebSocket Test FAILED" -ForegroundColor $errorColor
    $results += @{"Protocol" = "WebSocket"; "Status" = "FAILED"}
}

# Cleanup
Stop-Process -Id $serverProcess.Id -Force -ErrorAction SilentlyContinue
Write-Host ""

# Step 6: Summary
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor $infoColor
Write-Host "║ Test Summary                           ║" -ForegroundColor $infoColor
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor $infoColor
Write-Host ""

$passCount = ($results | Where-Object { $_.Status -eq "PASSED" }).Count
$failCount = ($results | Where-Object { $_.Status -eq "FAILED" }).Count

foreach ($result in $results) {
    if ($result.Status -eq "PASSED") {
        Write-Host "  ✅ $($result.Protocol): $($result.Status)" -ForegroundColor $successColor
    } else {
        Write-Host "  ❌ $($result.Protocol): $($result.Status)" -ForegroundColor $errorColor
    }
}

Write-Host ""
Write-Host "Passed: $passCount / Failed: $failCount" -ForegroundColor (if ($failCount -eq 0) { $successColor } else { $errorColor })

if ($failCount -eq 0) {
    Write-Host "✅ All tests PASSED!" -ForegroundColor $successColor
    exit 0
} else {
    Write-Host "❌ Some tests FAILED!" -ForegroundColor $errorColor
    exit 1
}
