package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// AirQualityData represents the JSON response from the PurpleAir sensor
type AirQualityData struct {
	SensorId              string  `json:"SensorId"`
	DateTime              string  `json:"DateTime"`
	Geo                   string  `json:"Geo"`
	Mem                   int     `json:"Mem"`
	Memfrag               int     `json:"memfrag"`
	Memfb                 int     `json:"memfb"`
	Memcs                 int     `json:"memcs"`
	Id                    int     `json:"Id"`
	Lat                   float64 `json:"lat"`
	Lon                   float64 `json:"lon"`
	Adc                   float64 `json:"Adc"`
	Loggingrate           int     `json:"loggingrate"`
	Place                 string  `json:"place"`
	Version               string  `json:"version"`
	Uptime                int     `json:"uptime"`
	Rssi                  int     `json:"rssi"`
	Period                int     `json:"period"`
	Httpsuccess           int     `json:"httpsuccess"`
	Httpsends             int     `json:"httpsends"`
	Hardwareversion       string  `json:"hardwareversion"`
	Hardwarediscovered    string  `json:"hardwarediscovered"`
	CurrentTempF          float64 `json:"current_temp_f"`
	CurrentHumidity       int     `json:"current_humidity"`
	CurrentDewpointF      float64 `json:"current_dewpoint_f"`
	Pressure              float64 `json:"pressure"`
	CurrentTempF680       float64 `json:"current_temp_f_680"`
	CurrentHumidity680    int     `json:"current_humidity_680"`
	CurrentDewpointF680   float64 `json:"current_dewpoint_f_680"`
	Pressure680           float64 `json:"pressure_680"`
	Gas680                float64 `json:"gas_680"`
	P25aqicB              string  `json:"p25aqic_b"`
	Pm25AqiB              int     `json:"pm2.5_aqi_b"`
	Pm10Cf1B              float64 `json:"pm1_0_cf_1_b"`
	P03UmB                float64 `json:"p_0_3_um_b"`
	Pm25Cf1B              float64 `json:"pm2.5_cf_1_b"`
	P05UmB                float64 `json:"p_0_5_um_b"`
	Pm100Cf1B             float64 `json:"pm10_0_cf_1_b"`
	P10UmB                float64 `json:"p_1_0_um_b"`
	Pm10AtmB              float64 `json:"pm1_0_atm_b"`
	P25UmB                float64 `json:"p_2_5_um_b"`
	Pm25AtmB              float64 `json:"pm2_5_atm_b"`
	P50UmB                float64 `json:"p_5_0_um_b"`
	Pm100AtmB             float64 `json:"pm10_0_atm_b"`
	P100UmB               float64 `json:"p_10_0_um_b"`
	P25aqic               string  `json:"p25aqic"`
	Pm25Aqi               int     `json:"pm2.5_aqi"`
	Pm10Cf1               float64 `json:"pm1_0_cf_1"`
	P03Um                 float64 `json:"p_0_3_um"`
	Pm25Cf1               float64 `json:"pm2_5_cf_1"`
	P05Um                 float64 `json:"p_0_5_um"`
	Pm100Cf1              float64 `json:"pm10_0_cf_1"`
	P10Um                 float64 `json:"p_1_0_um"`
	Pm10Atm               float64 `json:"pm1_0_atm"`
	P25Um                 float64 `json:"p_2_5_um"`
	Pm25Atm               float64 `json:"pm2_5_atm"`
	P50Um                 float64 `json:"p_5_0_um"`
	Pm100Atm              float64 `json:"pm10_0_atm"`
	P100Um                float64 `json:"p_10_0_um"`
	PaLatency             int     `json:"pa_latency"`
	Wlstate               string  `json:"wlstate"`
	Status0               int     `json:"status_0"`
	Status1               int     `json:"status_1"`
	Status2               int     `json:"status_2"`
	Status3               int     `json:"status_3"`
	Status4               int     `json:"status_4"`
	Ssid                  string  `json:"ssid"`
}

// fetchAirQualityData makes an HTTP request to the IoT device and returns the parsed data
func fetchAirQualityData(deviceURL string) (*AirQualityData, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(deviceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data AirQualityData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &data, nil
}

// printAirQualityData displays the air quality data in a formatted way
func printAirQualityData(data *AirQualityData) {
	fmt.Println("=== Air Quality Sensor Data ===")
	fmt.Printf("Sensor ID: %s\n", data.SensorId)
	fmt.Printf("Location: %s (%.6f, %.6f)\n", data.Geo, data.Lat, data.Lon)
	fmt.Printf("DateTime: %s\n", data.DateTime)
	fmt.Printf("Place: %s\n", data.Place)
	fmt.Printf("Version: %s\n", data.Version)
	fmt.Printf("Hardware: %s\n", data.Hardwareversion)
	fmt.Printf("Uptime: %d seconds\n", data.Uptime)
	fmt.Printf("WiFi: %s (RSSI: %d)\n", data.Wlstate, data.Rssi)
	fmt.Printf("SSID: %s\n", data.Ssid)
	
	fmt.Println("\n=== Environmental Data ===")
	fmt.Printf("Temperature: %.1f°F\n", data.CurrentTempF)
	fmt.Printf("Humidity: %d%%\n", data.CurrentHumidity)
	fmt.Printf("Dew Point: %.1f°F\n", data.CurrentDewpointF)
	fmt.Printf("Pressure: %.2f hPa\n", data.Pressure)
	fmt.Printf("Gas (BME680): %.2f kΩ\n", data.Gas680)
	
	fmt.Println("\n=== Air Quality (Channel A) ===")
	fmt.Printf("PM2.5 AQI: %d (%s)\n", data.Pm25Aqi, data.P25aqic)
	fmt.Printf("PM1.0 (CF1): %.2f μg/m³\n", data.Pm10Cf1)
	fmt.Printf("PM2.5 (CF1): %.2f μg/m³\n", data.Pm25Cf1)
	fmt.Printf("PM10.0 (CF1): %.2f μg/m³\n", data.Pm100Cf1)
	fmt.Printf("PM1.0 (ATM): %.2f μg/m³\n", data.Pm10Atm)
	fmt.Printf("PM2.5 (ATM): %.2f μg/m³\n", data.Pm25Atm)
	fmt.Printf("PM10.0 (ATM): %.2f μg/m³\n", data.Pm100Atm)
	
	fmt.Println("\n=== Air Quality (Channel B) ===")
	fmt.Printf("PM2.5 AQI: %d (%s)\n", data.Pm25AqiB, data.P25aqicB)
	fmt.Printf("PM1.0 (CF1): %.2f μg/m³\n", data.Pm10Cf1B)
	fmt.Printf("PM2.5 (CF1): %.2f μg/m³\n", data.Pm25Cf1B)
	fmt.Printf("PM10.0 (CF1): %.2f μg/m³\n", data.Pm100Cf1B)
	fmt.Printf("PM1.0 (ATM): %.2f μg/m³\n", data.Pm10AtmB)
	fmt.Printf("PM2.5 (ATM): %.2f μg/m³\n", data.Pm25AtmB)
	fmt.Printf("PM10.0 (ATM): %.2f μg/m³\n", data.Pm100AtmB)
	
	fmt.Println("\n=== System Status ===")
	fmt.Printf("Memory: %d bytes (frag: %d%%, free: %d, cache: %d)\n", 
		data.Mem, data.Memfrag, data.Memfb, data.Memcs)
	fmt.Printf("ADC: %.2fV\n", data.Adc)
	fmt.Printf("HTTP Success/Sends: %d/%d\n", data.Httpsuccess, data.Httpsends)
	fmt.Printf("PurpleAir Latency: %dms\n", data.PaLatency)
	fmt.Printf("Status: %d/%d/%d/%d/%d\n", 
		data.Status0, data.Status1, data.Status2, data.Status3, data.Status4)
}

func main() {
	// Check if server mode is requested
	if len(os.Args) > 1 && os.Args[1] == "--server" {
		// Server mode
		deviceURL := "http://192.168.1.100/json"
		serverAddr := ":8080"
		
		if len(os.Args) > 2 {
			deviceURL = os.Args[2]
		}
		if len(os.Args) > 3 {
			serverAddr = os.Args[3]
		}
		
		fmt.Printf("Starting server mode...\n")
		fmt.Printf("Device URL: %s\n", deviceURL)
		fmt.Printf("Server address: %s\n", serverAddr)
		fmt.Printf("Web interface: http://localhost%s\n", serverAddr)
		fmt.Printf("Graphs: http://localhost%s/graphs\n", serverAddr)
		fmt.Printf("API endpoints:\n")
		fmt.Printf("  - GET /data/json - Raw JSON data\n")
		fmt.Printf("  - GET /data - Formatted text data\n")
		fmt.Printf("  - GET /health - Health check\n")
		fmt.Printf("  - GET /graphs - Historical graphs\n")
		fmt.Printf("  - GET /api/measurements - Measurement data for graphing\n")
		fmt.Printf("  - GET /api/stats - Statistics\n\n")
		
		// Initialize database
		dbPath := os.Getenv("DATABASE_PATH")
		if dbPath == "" {
			dbPath = "air_quality.db"
		}
		database, err := NewDatabase(dbPath)
		if err != nil {
			log.Printf("Warning: Failed to initialize database: %v", err)
			log.Printf("Data storage and graphing will be disabled\n")
			database = nil
		} else {
			log.Printf("Database initialized successfully\n")
		}
		
		server := NewServer(deviceURL, database)
		log.Fatal(server.Start(serverAddr))
	} else {
		// Command-line mode
		deviceURL := "http://192.168.1.100/json"
		if len(os.Args) > 1 {
			deviceURL = os.Args[1]
		}

		fmt.Printf("Fetching air quality data from: %s\n\n", deviceURL)

		// Fetch data from the IoT device
		data, err := fetchAirQualityData(deviceURL)
		if err != nil {
			log.Fatalf("Error fetching air quality data: %v", err)
		}

		// Display the data
		printAirQualityData(data)
	}
}
