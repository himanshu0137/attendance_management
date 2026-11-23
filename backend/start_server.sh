#!/bin/bash

# Kill any existing server processes
echo "Stopping any existing server processes..."
pkill -f "./server"

# Wait a moment for cleanup
sleep 1

# Set up environment variables from .env file
set -a
source .env
set +a

# Run the server
echo "Starting Attendance Management Server on port $ATTENDANCE_PORT"
./server