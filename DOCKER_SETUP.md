# Docker Setup and Configuration

## Docker Compose Configuration

The system uses Docker Compose to orchestrate the following services:

### PostgreSQL Service
- **Image**: postgres:16
- **Port Mapping**: 5432:5432 (host:container)
- **Environment Variables**:
  - POSTGRES_USER: attendance
  - POSTGRES_PASSWORD: attendance  
  - POSTGRES_DB: attendance
- **Health Check**: Tests connection with pg_isready

### Keycloak Service
- **Image**: quay.io/keycloak/keycloak:latest
- **Port Mapping**: 9080:9080 (host:container)
- **Command**: `start-dev --http-port 9080`
- **Environment Variables**:
  - KC_BOOTSTRAP_ADMIN_USERS: admin
  - KC_HOSTNAME_STRICT: false
  - KC_HOSTNAME: localhost
  - KC_HTTP_HOST: 0.0.0.0
- **Health Check**: Tests HTTP endpoint

## Service Binding

The Keycloak service is configured to bind to `0.0.0.0` using the `KC_HTTP_HOST=0.0.0.0` environment variable, which allows external access to the Keycloak service.

## Running Services

### Start All Services
```bash
cd /home/carbon/prog/attendance_management/backend
docker compose up -d
```

### View Service Status
```bash
docker ps
```

### View Service Logs
```bash
# View PostgreSQL logs
docker logs <postgres_container_id>

# View Keycloak logs  
docker logs <keycloak_container_id>
```

### Stop Services
```bash
docker compose down
```

## Environment Variables for Backend

When the Docker services are running, the backend can be configured with:

```bash
export ATTENDANCE_PORT=7080
export ATTENDANCE_DB="postgresql://attendance:attendance@localhost:5432/attendance"
export KEYCLOAK_ISSUER="http://localhost:9080/realms/attendance"
export KEYCLOAK_CLIENT_ID="attendance-app"
export KEYCLOAK_CLIENT_SECRET="secret"
export LOG_LEVEL=debug
```

## Testing Connectivity

### Test PostgreSQL Connection
```bash
# Test from host machine
psql -h localhost -p 5432 -U attendance -d attendance
```

### Test Keycloak Accessibility  
```bash
# Test Keycloak endpoint
curl http://localhost:9080/
```

### Test Backend Server
```bash
# After starting the server with proper env vars
curl http://localhost:7080/health
```

## Security Considerations

1. **PostgreSQL**: Access is limited to localhost by default
2. **Keycloak**: Configured for development use, not production secure
3. **Backend**: Listens on 0.0.0.0 for external access (use in secure environments only)

## Troubleshooting

### Common Issues:

1. **Port Conflicts**:
   ```bash
   # Check if ports are in use
   netstat -tuln | grep -E "(5432|7080|9080)"
   ```

2. **Service Not Starting**:
   ```bash
   # View detailed logs
   docker compose logs <service_name>
   ```

3. **Network Issues**:
   ```bash
   # Verify Docker network
   docker network ls
   ```

## Production Considerations

For production deployment:
1. Use proper SSL certificates for Keycloak and backend
2. Configure firewall rules to restrict access
3. Set up proper database backups
4. Use environment-specific configurations