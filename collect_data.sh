#!/bin/bash

# Data Collection Script for Air Quality Monitor
# This script fetches data from the PurpleAir sensor and stores it in the database

echo "Starting data collection for air quality monitoring..."
echo "This will collect data every 30 seconds for testing the graphing functionality."
echo "Press Ctrl+C to stop."
echo

DEVICE_URL="http://192.168.0.249/json"
SERVER_URL="http://localhost:8080"

# Function to fetch and store data
fetch_data() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - Fetching data..."
    
    # Make a request to the server to trigger data collection
    response=$(curl -s "$SERVER_URL/data/json")
    
    if [ $? -eq 0 ]; then
        echo "  ✓ Data collected successfully"
    else
        echo "  ✗ Failed to collect data"
    fi
}

# Main loop
while true; do
    fetch_data
    echo "  Waiting 30 seconds until next collection..."
    echo
    sleep 30
done
