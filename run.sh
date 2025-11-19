#!/bin/bash

# Build the Go application
echo "Building the application..."
go build -o hubfly-storage ./cmd/hubfly-storage

# Check if the build was successful
if [ $? -ne 0 ]; then
    echo "Build failed. Please check the error messages."
    exit 1
fi

# Run the application with sudo
echo "Running the application with sudo..."
sudo ./hubfly-storage