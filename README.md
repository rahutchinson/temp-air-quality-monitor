# Air Quality Monitor

A Go application that fetches air quality data from PurpleAir IoT sensors and provides both command-line and web interface access to the data.

## Features

- **HTTP Client**: Makes requests to PurpleAir sensors on your local network
- **JSON Parsing**: Parses the complex JSON response from PurpleAir devices
- **Command-line Interface**: Simple CLI for one-time data fetching
- **Web Server**: Beautiful web interface with real-time data updates
- **REST API**: JSON endpoints for programmatic access
- **Health Monitoring**: System status and health check endpoints
- **Data Storage**: SQLite database for historical data collection
- **Interactive Graphs**: Real-time charts using Chart.js
- **Statistics**: Time-based analytics and trends
- **Auto-refresh**: Automatic data collection and visualization updates

## Installation

1. **Prerequisites**: Go 1.21 or later
2. **Clone and build**:
   ```bash
   git clone <repository-url>
   cd temp-air-quality-monitor
   go mod tidy
   go build -o air-quality-monitor .
   ```

## Usage

### Command-line Mode

Fetch air quality data once and display it in the terminal:

```bash
# Use default device URL (http://192.168.1.100/json)
./air-quality-monitor

# Specify custom device URL
./air-quality-monitor http://192.168.1.150/json

# Example output:
# === Air Quality Sensor Data ===
# Sensor ID: c8:c9:a3:2d:fd:4f
# Location: PurpleAir-fd4f (29.798401, -95.407898)
# Temperature: 101.0°F
# Humidity: 40%
# PM2.5 AQI: 33
```

### Web Server Mode

Start a web server to provide a beautiful web interface:

```bash
# Start server with default settings
./air-quality-monitor --server

# Specify custom device URL and port
./air-quality-monitor --server http://192.168.1.150/json :9090

# Access the web interface at http://localhost:8080
```

### API Endpoints

When running in server mode, the following endpoints are available:

- `GET /` - Web interface with real-time data
- `GET /graphs` - Interactive historical graphs and charts
- `GET /data/json` - Raw JSON data from the sensor
- `GET /data` - Formatted text data
- `GET /health` - Health check endpoint
- `GET /api/measurements` - Historical measurement data for graphing
- `GET /api/stats` - Statistical data for the specified time period

## Configuration

Edit `config.json` to customize the application behavior:

```json
{
  "device": {
    "url": "http://192.168.1.100/json",
    "timeout": 10,
    "retry_attempts": 3,
    "retry_delay": 5
  },
  "server": {
    "port": 8080,
    "host": "0.0.0.0",
    "refresh_interval": 30
  }
}
```

## Data Storage and Graphing

The application automatically stores all air quality measurements in a SQLite database (`air_quality.db`) for historical analysis and graphing.

### Graphing Features

- **Interactive Charts**: Real-time line charts using Chart.js
- **Multiple Time Ranges**: View data for 1 hour, 6 hours, 24 hours, or 1 week
- **Multiple Metrics**: 
  - PM2.5 Air Quality Index (both channels)
  - Temperature and Humidity
  - PM2.5 Concentration (μg/m³)
  - System metrics (Memory, WiFi signal strength)
- **Statistics Dashboard**: Average, min, max values for all metrics
- **Auto-refresh**: Charts update automatically every 5 minutes

### Data Collection

The application automatically stores data when:
- The web interface is accessed
- API endpoints are called
- The server is running and receiving requests

For continuous data collection, you can use the provided script:
```bash
./collect_data.sh
```

## Data Structure

The application parses the following PurpleAir sensor data:

### Environmental Data
- Temperature (°F)
- Humidity (%)
- Dew Point (°F)
- Pressure (hPa)
- Gas resistance (kΩ)

### Air Quality Data
- PM1.0, PM2.5, PM10.0 concentrations (μg/m³)
- Air Quality Index (AQI)
- Particle counts by size
- Both CF1 (correction factor 1) and ATM (atmospheric) measurements

### System Information
- Sensor ID and location
- Hardware version and uptime
- WiFi status and signal strength
- Memory usage and system status

## Finding Your PurpleAir Device

1. **Network Discovery**: PurpleAir devices typically serve data on port 80
2. **Common URLs**:
   - `http://[device-ip]/json` - JSON data endpoint
   - `http://[device-ip]/` - Web interface
3. **Finding IP Address**:
   - Check your router's DHCP client list
   - Use network scanning tools like `nmap`
   - Look for devices with hostnames containing "PurpleAir"

## Troubleshooting

### Connection Issues
- Verify the device IP address is correct
- Ensure the device is on the same network
- Check if the device is powered on and connected to WiFi

### Data Parsing Errors
- Verify the device is a PurpleAir sensor
- Check if the JSON response format matches expectations
- Ensure the device firmware is up to date

### Web Interface Issues
- Check if the port is available and not blocked by firewall
- Verify the server is running on the expected address
- Check browser console for JavaScript errors

## Development

### Project Structure
```
temp-air-quality-monitor/
├── main.go              # Main application logic
├── server.go            # Web server implementation
├── database.go          # Database operations and data storage
├── go.mod               # Go module definition
├── config.json          # Configuration file
├── air_quality.db       # SQLite database (created automatically)
├── collect_data.sh      # Data collection script
├── test.sh              # Test script
└── README.md            # This file
```

### Building
```bash
# Build for current platform
go build -o air-quality-monitor .

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o air-quality-monitor-linux .
GOOS=windows GOARCH=amd64 go build -o air-quality-monitor.exe .
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Acknowledgments

- PurpleAir for providing the sensor hardware and data format
- The Go community for excellent HTTP and JSON libraries
