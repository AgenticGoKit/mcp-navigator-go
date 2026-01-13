# Test the math MCP server

Write-Host "Testing math-server.exe..." -ForegroundColor Cyan

# Test tools/list
$request = '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}'
Write-Host "`nSending tools/list request:" -ForegroundColor Yellow
Write-Host $request

Write-Host "`nResponse:" -ForegroundColor Yellow
$request | .\math-server.exe

Write-Host "`n"
