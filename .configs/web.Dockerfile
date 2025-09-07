# Stage 1: Build the React application
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Stage 2: Create production image with Nginx
FROM nginxinc/nginx-unprivileged:1.29.1-alpine3.22

# Switch to non-root user (nginx-unprivileged already uses nginx user)
USER nginx

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built application from builder stage
COPY --from=builder --chown=nginx:nginx /app/build /usr/share/nginx/html

# Copy health check script
COPY healthcheck.sh /healthcheck.sh

# Install curl for health checks and monitoring
USER root
RUN apk add --no-cache curl && \
    chmod +x /healthcheck.sh && \
    chown nginx:nginx /healthcheck.sh
USER nginx

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/healthcheck.sh"]

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]