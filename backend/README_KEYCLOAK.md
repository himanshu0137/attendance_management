# Keycloak Integration

The Attendance Management System uses Keycloak for authentication and authorization. This document explains how Keycloak is configured and used.

## Keycloak Configuration

Keycloak is automatically initialized with:
1. A realm named `attendance-management`
2. A client application `attendance-app` 
3. An admin user with credentials `admin:admin123`

## How It Works

1. **Automatic Setup**: When you run `docker-compose up`, the Keycloak container automatically runs our initialization script
2. **Initialization Script**: The script creates the realm, client, and admin user if they don't already exist
3. **Security**: All authentication is handled through Keycloak's OAuth 2.0 and OpenID Connect protocols

## Accessing Keycloak

Once running:
- Admin Console: http://localhost:9080
- Admin Credentials: `admin` / `admin123`
- Realm: `attendance-management`

## Customization

To modify the Keycloak setup:
1. Edit `/keycloak-init/init-keycloak.sh`
2. Update realm, client, or user configurations as needed
3. Restart services with `docker-compose up -d`

## Testing Authentication

You can test authentication by:
1. Accessing the Keycloak admin console
2. Creating additional users and roles
3. Using the attendance-app client to authenticate users

## Production Considerations

For production environments:
1. Change default admin password
2. Configure HTTPS for Keycloak
3. Set up proper backup procedures
4. Implement proper monitoring