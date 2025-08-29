package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server represents the web server for serving air quality data
type Server struct {
	deviceURL string
	router    *mux.Router
	database  *Database
}

// NewServer creates a new server instance
func NewServer(deviceURL string, database *Database) *Server {
	s := &Server{
		deviceURL: deviceURL,
		router:    mux.NewRouter(),
		database:  database,
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	s.router.HandleFunc("/", s.handleHome).Methods("GET")
	s.router.HandleFunc("/data", s.handleGetData).Methods("GET")
	s.router.HandleFunc("/data/json", s.handleGetDataJSON).Methods("GET")
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
	s.router.HandleFunc("/graphs", s.handleGraphs).Methods("GET")
	s.router.HandleFunc("/api/measurements", s.handleGetMeasurements).Methods("GET")
	s.router.HandleFunc("/api/stats", s.handleGetStats).Methods("GET")
}

// handleHome serves the home page
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Air Quality Monitor</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px; margin-bottom: 20px; }
        .data-section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .data-section h3 { color: #007bff; margin-top: 0; }
        .metric { display: flex; justify-content: space-between; margin: 5px 0; padding: 5px 0; border-bottom: 1px solid #eee; }
        .metric:last-child { border-bottom: none; }
        .value { font-weight: bold; color: #28a745; }
        .refresh-btn { background: #007bff; color: white; border: none; padding: 10px 20px; border-radius: 5px; cursor: pointer; margin: 10px 0; }
        .refresh-btn:hover { background: #0056b3; }
        .last-updated { text-align: center; color: #666; font-size: 0.9em; margin-top: 20px; }
        .loading { text-align: center; color: #666; font-style: italic; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Air Quality Monitor</h1>
            <p>Real-time air quality data from PurpleAir sensor</p>
        </div>
        
        <button class="refresh-btn" onclick="location.reload()">Refresh Data</button>
        <a href="/graphs" class="refresh-btn" style="text-decoration: none; display: inline-block; margin-left: 10px;">View Graphs</a>
        
        <div id="data-container">
            <div class="loading">Loading data...</div>
        </div>
        
        <div class="last-updated" id="last-updated"></div>
    </div>

    <script>
        function updateData() {
            fetch('/data/json')
                .then(response => response.json())
                .then(data => {
                    const container = document.getElementById('data-container');
                    let html = '';
                    
                    html += '<div class="data-section">';
                    html += '<h3>Current Air Quality</h3>';
                    html += '<div class="metric"><span>PM2.5 AQI (Channel A):</span><span class="value">' + data.pm2_5_aqi + '</span></div>';
                    html += '<div class="metric"><span>PM2.5 AQI (Channel B):</span><span class="value">' + data.pm2_5_aqi_b + '</span></div>';
                    html += '<div class="metric"><span>PM1.0 (CF1):</span><span class="value">' + data.pm1_0_cf_1.toFixed(2) + ' ug/m3</span></div>';
                    html += '<div class="metric"><span>PM2.5 (CF1):</span><span class="value">' + data.pm2_5_cf_1.toFixed(2) + ' ug/m3</span></div>';
                    html += '<div class="metric"><span>PM10.0 (CF1):</span><span class="value">' + data.pm10_0_cf_1.toFixed(2) + ' ug/m3</span></div>';
                    html += '</div>';
                    
                    html += '<div class="data-section">';
                    html += '<h3>Environmental Conditions</h3>';
                    html += '<div class="metric"><span>Temperature:</span><span class="value">' + data.current_temp_f.toFixed(1) + ' F</span></div>';
                    html += '<div class="metric"><span>Humidity:</span><span class="value">' + data.current_humidity + '%</span></div>';
                    html += '<div class="metric"><span>Dew Point:</span><span class="value">' + data.current_dewpoint_f.toFixed(1) + ' F</span></div>';
                    html += '<div class="metric"><span>Pressure:</span><span class="value">' + data.pressure.toFixed(2) + ' hPa</span></div>';
                    html += '<div class="metric"><span>Gas (BME680):</span><span class="value">' + data.gas_680.toFixed(2) + ' kOhm</span></div>';
                    html += '</div>';
                    
                    html += '<div class="data-section">';
                    html += '<h3>System Information</h3>';
                    html += '<div class="metric"><span>Sensor ID:</span><span class="value">' + data.SensorId + '</span></div>';
                    html += '<div class="metric"><span>Location:</span><span class="value">' + data.Geo + '</span></div>';
                    html += '<div class="metric"><span>Uptime:</span><span class="value">' + Math.floor(data.uptime / 3600) + 'h ' + Math.floor((data.uptime % 3600) / 60) + 'm</span></div>';
                    html += '<div class="metric"><span>WiFi Status:</span><span class="value">' + data.wlstate + ' (RSSI: ' + data.rssi + ')</span></div>';
                    html += '<div class="metric"><span>Memory:</span><span class="value">' + data.Mem + ' bytes</span></div>';
                    html += '</div>';
                    
                    container.innerHTML = html;
                    
                    document.getElementById('last-updated').textContent = 
                        'Last updated: ' + new Date().toLocaleString();
                })
                .catch(error => {
                    console.error('Error fetching data:', error);
                    document.getElementById('data-container').innerHTML = 
                        '<div style="color: red; text-align: center; padding: 20px;">Error loading data. Please try again.</div>';
                });
        }
        
        // Load data immediately and refresh every 30 seconds
        updateData();
        setInterval(updateData, 30000);
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// handleGetData serves the data as formatted text
func (s *Server) handleGetData(w http.ResponseWriter, r *http.Request) {
	data, err := fetchAirQualityData(s.deviceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	printAirQualityData(data)
}

// handleGetDataJSON serves the data as JSON
func (s *Server) handleGetDataJSON(w http.ResponseWriter, r *http.Request) {
	data, err := fetchAirQualityData(s.deviceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the measurement in the database
	if s.database != nil {
		if err := s.database.StoreMeasurement(data); err != nil {
			log.Printf("Warning: Failed to store measurement: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(data)
}

// handleHealth serves a health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "air-quality-monitor",
	})
}

// handleGraphs serves the graphs page
func (s *Server) handleGraphs(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Air Quality Graphs</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px; margin-bottom: 20px; }
        .chart-container { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .chart-container h3 { color: #007bff; margin-top: 0; }
        .controls { margin: 20px 0; text-align: center; }
        .controls select, .controls button { margin: 0 10px; padding: 8px 16px; border: 1px solid #ddd; border-radius: 4px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; margin: 20px 0; }
        .stat-card { background: #f8f9fa; padding: 15px; border-radius: 5px; text-align: center; }
        .stat-value { font-size: 24px; font-weight: bold; color: #007bff; }
        .stat-label { color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Air Quality Graphs</h1>
            <p>Historical air quality data from PurpleAir sensor</p>
        </div>
        
        <div class="controls">
            <label for="timeRange">Time Range:</label>
            <select id="timeRange" onchange="loadData()">
                <option value="1">Last Hour</option>
                <option value="6">Last 6 Hours</option>
                <option value="24" selected>Last 24 Hours</option>
                <option value="168">Last Week</option>
            </select>
            <button onclick="loadData()">Refresh</button>
        </div>
        
        <div id="stats" class="stats">
            <div class="stat-card">
                <div class="stat-value" id="avgPM25">-</div>
                <div class="stat-label">Avg PM2.5 AQI</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="avgTemp">-</div>
                <div class="stat-label">Avg Temperature (°F)</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="avgHumidity">-</div>
                <div class="stat-label">Avg Humidity (%)</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="dataPoints">-</div>
                <div class="stat-label">Data Points</div>
            </div>
        </div>
        
        <div class="chart-container">
            <h3>PM2.5 Air Quality Index</h3>
            <canvas id="pm25Chart" width="400" height="200"></canvas>
        </div>
        
        <div class="chart-container">
            <h3>Temperature and Humidity</h3>
            <canvas id="tempHumidityChart" width="400" height="200"></canvas>
        </div>
        
        <div class="chart-container">
            <h3>PM2.5 Concentration (μg/m³)</h3>
            <canvas id="pm25ConcentrationChart" width="400" height="200"></canvas>
        </div>
        
        <div class="chart-container">
            <h3>System Metrics</h3>
            <canvas id="systemChart" width="400" height="200"></canvas>
        </div>
    </div>

    <script>
        let pm25Chart, tempHumidityChart, pm25ConcentrationChart, systemChart;
        
        function initCharts() {
            const ctx1 = document.getElementById('pm25Chart').getContext('2d');
            const ctx2 = document.getElementById('tempHumidityChart').getContext('2d');
            const ctx3 = document.getElementById('pm25ConcentrationChart').getContext('2d');
            const ctx4 = document.getElementById('systemChart').getContext('2d');
            
            pm25Chart = new Chart(ctx1, {
                type: 'line',
                data: { labels: [], datasets: [] },
                options: {
                    responsive: true,
                    scales: { y: { beginAtZero: true } }
                }
            });
            
            tempHumidityChart = new Chart(ctx2, {
                type: 'line',
                data: { labels: [], datasets: [] },
                options: {
                    responsive: true,
                    scales: { y: { beginAtZero: true } }
                }
            });
            
            pm25ConcentrationChart = new Chart(ctx3, {
                type: 'line',
                data: { labels: [], datasets: [] },
                options: {
                    responsive: true,
                    scales: { y: { beginAtZero: true } }
                }
            });
            
            systemChart = new Chart(ctx4, {
                type: 'line',
                data: { labels: [], datasets: [] },
                options: {
                    responsive: true,
                    scales: { y: { beginAtZero: true } }
                }
            });
        }
        
        function loadData() {
            const hours = document.getElementById('timeRange').value;
            
            // Load measurements
            fetch('/api/measurements?hours=' + hours)
                .then(response => response.json())
                .then(data => {
                    updateCharts(data);
                })
                .catch(error => {
                    console.error('Error loading measurements:', error);
                });
            
            // Load stats
            fetch('/api/stats?hours=' + hours)
                .then(response => response.json())
                .then(data => {
                    updateStats(data);
                })
                .catch(error => {
                    console.error('Error loading stats:', error);
                });
        }
        
        function updateCharts(measurements) {
            const labels = measurements.map(m => new Date(m.timestamp).toLocaleTimeString());
            const pm25Data = measurements.map(m => m.pm25_aqi);
            const pm25BData = measurements.map(m => m.pm25_aqi_b);
            const tempData = measurements.map(m => m.temperature);
            const humidityData = measurements.map(m => m.humidity);
            const pm25ConcentrationData = measurements.map(m => m.pm25_cf1);
            const pm25BConcentrationData = measurements.map(m => m.pm25_cf1_b);
            const memoryData = measurements.map(m => m.memory);
            const rssiData = measurements.map(m => m.rssi);
            
            // Update PM2.5 AQI chart
            pm25Chart.data.labels = labels;
            pm25Chart.data.datasets = [
                { label: 'PM2.5 AQI (Channel A)', data: pm25Data, borderColor: 'rgb(75, 192, 192)', backgroundColor: 'rgba(75, 192, 192, 0.2)' },
                { label: 'PM2.5 AQI (Channel B)', data: pm25BData, borderColor: 'rgb(255, 99, 132)', backgroundColor: 'rgba(255, 99, 132, 0.2)' }
            ];
            pm25Chart.update();
            
            // Update Temperature and Humidity chart
            tempHumidityChart.data.labels = labels;
            tempHumidityChart.data.datasets = [
                { label: 'Temperature (°F)', data: tempData, borderColor: 'rgb(255, 159, 64)', backgroundColor: 'rgba(255, 159, 64, 0.2)', yAxisID: 'y' },
                { label: 'Humidity (%)', data: humidityData, borderColor: 'rgb(153, 102, 255)', backgroundColor: 'rgba(153, 102, 255, 0.2)', yAxisID: 'y1' }
            ];
            tempHumidityChart.options.scales.y = { type: 'linear', display: true, position: 'left' };
            tempHumidityChart.options.scales.y1 = { type: 'linear', display: true, position: 'right', grid: { drawOnChartArea: false } };
            tempHumidityChart.update();
            
            // Update PM2.5 Concentration chart
            pm25ConcentrationChart.data.labels = labels;
            pm25ConcentrationChart.data.datasets = [
                { label: 'PM2.5 CF1 (Channel A)', data: pm25ConcentrationData, borderColor: 'rgb(54, 162, 235)', backgroundColor: 'rgba(54, 162, 235, 0.2)' },
                { label: 'PM2.5 CF1 (Channel B)', data: pm25BConcentrationData, borderColor: 'rgb(255, 205, 86)', backgroundColor: 'rgba(255, 205, 86, 0.2)' }
            ];
            pm25ConcentrationChart.update();
            
            // Update System chart
            systemChart.data.labels = labels;
            systemChart.data.datasets = [
                { label: 'Memory (bytes)', data: memoryData, borderColor: 'rgb(201, 203, 207)', backgroundColor: 'rgba(201, 203, 207, 0.2)' },
                { label: 'RSSI (dBm)', data: rssiData, borderColor: 'rgb(255, 99, 132)', backgroundColor: 'rgba(255, 99, 132, 0.2)' }
            ];
            systemChart.update();
        }
        
        function updateStats(stats) {
            document.getElementById('avgPM25').textContent = stats.avg_pm25_aqi ? stats.avg_pm25_aqi.toFixed(1) : '-';
            document.getElementById('avgTemp').textContent = stats.avg_temp ? stats.avg_temp.toFixed(1) : '-';
            document.getElementById('avgHumidity').textContent = stats.avg_humidity ? stats.avg_humidity.toFixed(1) : '-';
            document.getElementById('dataPoints').textContent = stats.count || '-';
        }
        
        // Initialize charts and load data
        initCharts();
        loadData();
        
        // Auto-refresh every 5 minutes
        setInterval(loadData, 300000);
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// handleGetMeasurements serves measurement data for graphing
func (s *Server) handleGetMeasurements(w http.ResponseWriter, r *http.Request) {
	if s.database == nil {
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}

	hours := 24 // default to 24 hours
	if h := r.URL.Query().Get("hours"); h != "" {
		if parsed, err := fmt.Sscanf(h, "%d", &hours); err != nil || parsed != 1 {
			hours = 24
		}
	}

	measurements, err := s.database.GetRecentMeasurements(hours)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching measurements: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(measurements)
}

// handleGetStats serves measurement statistics
func (s *Server) handleGetStats(w http.ResponseWriter, r *http.Request) {
	if s.database == nil {
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}

	hours := 24 // default to 24 hours
	if h := r.URL.Query().Get("hours"); h != "" {
		if parsed, err := fmt.Sscanf(h, "%d", &hours); parsed != 1 || err != nil {
			hours = 24
		}
	}

	stats, err := s.database.GetMeasurementStats(hours)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(stats)
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
