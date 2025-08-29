#!/bin/bash

# Air Quality Monitor Test Script

echo "=== Air Quality Monitor Test ==="
echo

# Test 1: Build the application
echo "1. Building the application..."
if go build -o air-quality-monitor .; then
    echo "   ✓ Build successful"
else
    echo "   ✗ Build failed"
    exit 1
fi

echo

# Test 2: Show help/usage
echo "2. Testing command-line usage..."
echo "   Command: ./air-quality-monitor"
echo "   Expected: Error (no device URL provided)"
./air-quality-monitor 2>&1 | head -5
echo

# Test 3: Test server mode startup
echo "3. Testing server mode startup..."
echo "   Starting server in background..."
./air-quality-monitor --server > server.log 2>&1 &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Test health endpoint
echo "   Testing health endpoint..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "   ✓ Server is running and responding"
else
    echo "   ✗ Server failed to start"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

# Test data endpoint (will fail without real device, but should return error properly)
echo "   Testing data endpoint..."
curl -s http://localhost:8080/data/json | head -3

# Stop server
echo "   Stopping server..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

echo

# Test 4: Show available endpoints
echo "4. Available endpoints when running in server mode:"
echo "   - GET /              - Web interface"
echo "   - GET /data/json     - Raw JSON data"
echo "   - GET /data          - Formatted text data"
echo "   - GET /health        - Health check"

echo

# Test 5: Show usage examples
echo "5. Usage examples:"
echo "   # Command-line mode (one-time fetch):"
echo "   ./air-quality-monitor http://192.168.1.100/json"
echo
echo "   # Server mode (web interface):"
echo "   ./air-quality-monitor --server"
echo
echo "   # Server mode with custom device URL:"
echo "   ./air-quality-monitor --server http://192.168.1.150/json"
echo
echo "   # Server mode with custom device URL and port:"
echo "   ./air-quality-monitor --server http://192.168.1.150/json :9090"

echo

echo "=== Test completed ==="
echo "The application is ready to use!"
echo "To find your PurpleAir device IP address, check your router's DHCP client list"
echo "or use network scanning tools like 'nmap -sn 192.168.1.0/24'"
