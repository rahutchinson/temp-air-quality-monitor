package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database connection and operations
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Database{db: db}, nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	// Main measurements table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS measurements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		sensor_id TEXT,
		datetime TEXT,
		geo TEXT,
		lat REAL,
		lon REAL,
		place TEXT,
		version TEXT,
		uptime INTEGER,
		rssi INTEGER,
		wlstate TEXT,
		ssid TEXT,
		
		-- Environmental data
		current_temp_f REAL,
		current_humidity INTEGER,
		current_dewpoint_f REAL,
		pressure REAL,
		gas_680 REAL,
		
		-- Air quality data (Channel A)
		pm25_aqi INTEGER,
		pm10_cf1 REAL,
		pm25_cf1 REAL,
		pm100_cf1 REAL,
		pm10_atm REAL,
		pm25_atm REAL,
		pm100_atm REAL,
		
		-- Air quality data (Channel B)
		pm25_aqi_b INTEGER,
		pm10_cf1_b REAL,
		pm25_cf1_b REAL,
		pm100_cf1_b REAL,
		pm10_atm_b REAL,
		pm25_atm_b REAL,
		pm100_atm_b REAL,
		
		-- System data
		mem INTEGER,
		memfrag INTEGER,
		memfb INTEGER,
		memcs INTEGER,
		adc REAL,
		httpsuccess INTEGER,
		httpsends INTEGER,
		pa_latency INTEGER,
		status_0 INTEGER,
		status_1 INTEGER,
		status_2 INTEGER,
		status_3 INTEGER,
		status_4 INTEGER
	);
	
	CREATE INDEX IF NOT EXISTS idx_measurements_timestamp ON measurements(timestamp);
	CREATE INDEX IF NOT EXISTS idx_measurements_sensor_id ON measurements(sensor_id);
	`

	_, err := db.Exec(createTableSQL)
	return err
}

// StoreMeasurement stores a single air quality measurement
func (d *Database) StoreMeasurement(data *AirQualityData) error {
	insertSQL := `
	INSERT INTO measurements (
		sensor_id, datetime, geo, lat, lon, place, version, uptime, rssi, wlstate, ssid,
		current_temp_f, current_humidity, current_dewpoint_f, pressure, gas_680,
		pm25_aqi, pm10_cf1, pm25_cf1, pm100_cf1, pm10_atm, pm25_atm, pm100_atm,
		pm25_aqi_b, pm10_cf1_b, pm25_cf1_b, pm100_cf1_b, pm10_atm_b, pm25_atm_b, pm100_atm_b,
		mem, memfrag, memfb, memcs, adc, httpsuccess, httpsends, pa_latency,
		status_0, status_1, status_2, status_3, status_4
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(insertSQL,
		data.SensorId, data.DateTime, data.Geo, data.Lat, data.Lon, data.Place, data.Version, data.Uptime, data.Rssi, data.Wlstate, data.Ssid,
		data.CurrentTempF, data.CurrentHumidity, data.CurrentDewpointF, data.Pressure, data.Gas680,
		data.Pm25Aqi, data.Pm10Cf1, data.Pm25Cf1, data.Pm100Cf1, data.Pm10Atm, data.Pm25Atm, data.Pm100Atm,
		data.Pm25AqiB, data.Pm10Cf1B, data.Pm25Cf1B, data.Pm100Cf1B, data.Pm10AtmB, data.Pm25AtmB, data.Pm100AtmB,
		data.Mem, data.Memfrag, data.Memfb, data.Memcs, data.Adc, data.Httpsuccess, data.Httpsends, data.PaLatency,
		data.Status0, data.Status1, data.Status2, data.Status3, data.Status4,
	)

	if err != nil {
		return fmt.Errorf("failed to insert measurement: %w", err)
	}

	return nil
}

// GetRecentMeasurements retrieves recent measurements for graphing
func (d *Database) GetRecentMeasurements(hours int) ([]Measurement, error) {
	query := `
	SELECT 
		timestamp, sensor_id, current_temp_f, current_humidity, pressure, gas_680,
		pm25_aqi, pm25_cf1, pm100_cf1, pm25_aqi_b, pm25_cf1_b, pm100_cf1_b,
		mem, rssi, pa_latency
	FROM measurements 
	WHERE timestamp >= datetime('now', '-` + fmt.Sprintf("%d", hours) + ` hours')
	ORDER BY timestamp ASC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query measurements: %w", err)
	}
	defer rows.Close()

	var measurements []Measurement
	for rows.Next() {
		var m Measurement
		err := rows.Scan(
			&m.Timestamp, &m.SensorID, &m.Temperature, &m.Humidity, &m.Pressure, &m.Gas680,
			&m.PM25AQI, &m.PM25CF1, &m.PM100CF1, &m.PM25AQIB, &m.PM25CF1B, &m.PM100CF1B,
			&m.Memory, &m.RSSI, &m.PaLatency,
		)
		if err != nil {
			log.Printf("Error scanning measurement: %v", err)
			continue
		}
		measurements = append(measurements, m)
	}

	return measurements, nil
}

// GetMeasurementStats returns statistics for the specified time period
func (d *Database) GetMeasurementStats(hours int) (*MeasurementStats, error) {
	query := `
	SELECT 
		COUNT(*) as count,
		COALESCE(AVG(current_temp_f), 0) as avg_temp,
		COALESCE(AVG(current_humidity), 0) as avg_humidity,
		COALESCE(AVG(pressure), 0) as avg_pressure,
		COALESCE(AVG(pm25_aqi), 0) as avg_pm25_aqi,
		COALESCE(AVG(pm25_cf1), 0) as avg_pm25_cf1,
		COALESCE(AVG(pm100_cf1), 0) as avg_pm100_cf1,
		COALESCE(MAX(pm25_aqi), 0) as max_pm25_aqi,
		COALESCE(MIN(pm25_aqi), 0) as min_pm25_aqi,
		COALESCE(MAX(current_temp_f), 0) as max_temp,
		COALESCE(MIN(current_temp_f), 0) as min_temp
	FROM measurements 
	WHERE timestamp >= datetime('now', '-` + fmt.Sprintf("%d", hours) + ` hours')
	`

	var stats MeasurementStats
	err := d.db.QueryRow(query).Scan(
		&stats.Count, &stats.AvgTemp, &stats.AvgHumidity, &stats.AvgPressure,
		&stats.AvgPM25AQI, &stats.AvgPM25CF1, &stats.AvgPM100CF1,
		&stats.MaxPM25AQI, &stats.MinPM25AQI, &stats.MaxTemp, &stats.MinTemp,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &stats, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// Measurement represents a simplified measurement for graphing
type Measurement struct {
	Timestamp   time.Time `json:"timestamp"`
	SensorID    string    `json:"sensor_id"`
	Temperature float64   `json:"temperature"`
	Humidity    int       `json:"humidity"`
	Pressure    float64   `json:"pressure"`
	Gas680      float64   `json:"gas_680"`
	PM25AQI     int       `json:"pm25_aqi"`
	PM25CF1     float64   `json:"pm25_cf1"`
	PM100CF1    float64   `json:"pm100_cf1"`
	PM25AQIB    int       `json:"pm25_aqi_b"`
	PM25CF1B    float64   `json:"pm25_cf1_b"`
	PM100CF1B   float64   `json:"pm100_cf1_b"`
	Memory      int       `json:"memory"`
	RSSI        int       `json:"rssi"`
	PaLatency   int       `json:"pa_latency"`
}

// MeasurementStats represents statistics for a time period
type MeasurementStats struct {
	Count       int     `json:"count"`
	AvgTemp     float64 `json:"avg_temp"`
	AvgHumidity float64 `json:"avg_humidity"`
	AvgPressure float64 `json:"avg_pressure"`
	AvgPM25AQI  float64 `json:"avg_pm25_aqi"`
	AvgPM25CF1  float64 `json:"avg_pm25_cf1"`
	AvgPM100CF1 float64 `json:"avg_pm100_cf1"`
	MaxPM25AQI  int     `json:"max_pm25_aqi"`
	MinPM25AQI  int     `json:"min_pm25_aqi"`
	MaxTemp     float64 `json:"max_temp"`
	MinTemp     float64 `json:"min_temp"`
}
