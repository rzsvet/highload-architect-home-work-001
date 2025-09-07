#!/bin/sh

# Health check script for React app with Nginx

# Check if Nginx is running
if ! pgrep nginx > /dev/null; then
    echo "Nginx is not running"
    exit 1
fi

# Check if the application is responding
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "Application is healthy"
    exit 0
else
    echo "Application health check failed"
    exit 1
fi