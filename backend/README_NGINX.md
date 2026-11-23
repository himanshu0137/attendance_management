# Nginx Proxy Configuration

This document explains how the Nginx proxy is configured to expose both the backend API and web application.

## Overview

The Nginx configuration provides a unified entrypoint for accessing:
- Backend API at `/api/` 
- Web Application at `/`

## Configuration Details

### Ports
- HTTP: 80 (redirects to HTTPS)
- HTTPS: 443

### Proxy Rules
1. **API Requests**: All requests starting with `/api/` are proxied to the backend service on port 7080
2. **Web Application**: All other requests are proxied to the web application

## Accessing Services

Once the system is running:
- Backend API: `http://localhost/api/`
- Web Application: `http://localhost/`

## SSL Configuration

For production use, you should:
1. Replace the placeholder SSL certificates in `/etc/nginx/ssl/` with real certificates
2. Update the Nginx configuration to point to your actual certificate files

## Testing

You can test the setup by:
1. Accessing `http://localhost/` to see the web application
2. Making API calls to `http://localhost/api/` endpoints

## Customization

To modify the proxy behavior:
1. Edit `/nginx/nginx.conf`
2. Update the upstream server definitions if needed
3. Modify location blocks for different routing rules
4. Restart services with `docker-compose up -d`

## Health Checks

The Nginx service includes health checks that verify it can reach both backend and webapp services.

## Security Headers

The configuration includes several security headers:
- X-Frame-Options: SAMEORIGIN
- X-XSS-Protection: 1; mode=block  
- X-Content-Type-Options: nosniff
- Referrer-Policy: no-referrer-when-downgrade
- Content-Security-Policy: default-src * 'unsafe-inline' 'unsafe-eval'