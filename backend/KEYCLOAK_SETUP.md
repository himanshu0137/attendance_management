# Keycloak Setup Instructions

This document explains how to set up Keycloak for the Attendance Management System.

## Overview

The system uses Keycloak as the authentication and authorization server. The setup includes:
- A custom realm called "attendance-management"
- An application client named "attendance-app"
- An admin user with appropriate permissions

## Prerequisites

- Docker and Docker Compose installed
- Port 9080 available for Keycloak
- Port 5432 available for PostgreSQL

## Setup Process

### 1. Running the System

To start the entire system including Keycloak with auto-initialization:

```bash
cd backend
docker-compose up -d
```

This will:
- Start PostgreSQL database
- Start Keycloak server
- Automatically initialize Keycloak with the required realm and admin user

### 2. Manual Initialization (if needed)

If you need to run the initialization script manually, use:

```bash
cd backend
chmod +x keycloak-init/init-keycloak.sh
./keycloak-init/init-keycloak.sh
```

### 3. Accessing Keycloak Admin Console

Once running, access the Keycloak admin console at:
- URL: http://localhost:9080
- Username: admin
- Password: admin123

## Keycloak Configuration Details

### Realm
- Name: attendance-management
- Display Name: Attendance Management System

### Client
- Client ID: attendance-app
- Redirect URIs: http://localhost:4200/*
- Web Origins: http://localhost:4200
- Public Client: true

### Admin User
- Username: admin
- Email: admin@attendance.com
- Password: admin123
- Role: admin

## Troubleshooting

### Keycloak Not Starting
If Keycloak fails to start, check the logs:
```bash
docker-compose logs keycloak
```

### Initialization Issues
If the initialization script fails, ensure:
1. PostgreSQL is running and accessible
2. Keycloak has started before the script runs
3. Required ports are available

### Manual Reset
To reset Keycloak configuration:
```bash
docker-compose down
docker-compose up -d
```
This will recreate the containers with fresh initialization.

## Security Notes

- The default admin password should be changed in production environments
- In production, ensure Keycloak is configured with HTTPS
- The realm and client configurations are defined in the initialization script