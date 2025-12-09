#!/bin/bash

# Exit on error
set -e

echo "Starting deployment process..."

# Define paths
APP_DIR="/home/plx/app"
LOG_DIR="/home/plx/logs"
BUILD_DIR="$APP_DIR/build"
SUPERVISOR_CONF="/etc/supervisor/conf.d/pulse.conf"

# Ensure directories exist
echo "Creating necessary directories..."
mkdir -p $LOG_DIR
mkdir -p $BUILD_DIR

# Navigate to project directory
cd $APP_DIR

echo "Installing/Updating Go dependencies..."
go mod download

echo "Building server binary..."
go build -o $BUILD_DIR/server cmd/server/main.go

echo "Building worker binary..."
go build -o $BUILD_DIR/worker cmd/worker/main.go

echo "Running database migrations..."
go run cmd/db/migrate.go up

# Restart services
echo "Restarting services..."
if [ -f "$SUPERVISOR_CONF" ]; then
    sudo supervisorctl reread
    sudo supervisorctl update
    sudo supervisorctl restart pulse:*  # Restart all programs in pulse group
else
    echo "Warning: Supervisor configuration file not found at $SUPERVISOR_CONF"
    echo "Please ensure supervisor configuration is properly installed"
    exit 1
fi

echo "Deployment completed successfully!"
