#!/bin/bash

# Keycloak initialization script
# This script sets up the required realm and admin user for the attendance management system

set -e

KEYCLOAK_URL="http://localhost:9080"
REALM_NAME="attendance-management"
ADMIN_USERNAME="admin"
ADMIN_PASSWORD="admin123"
ADMIN_EMAIL="admin@attendance.com"

echo "Initializing Keycloak..."

# Wait for Keycloak to be ready
echo "Waiting for Keycloak to start..."
until curl -f -s $KEYCLOAK_URL/realms > /dev/null; do
    sleep 2
done

echo "Keycloak is ready."

# Get admin token
echo "Getting admin token..."
ADMIN_TOKEN=$(curl -s -X POST \
  "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'username=admin' \
  -d 'password=admin' \
  -d 'grant_type=password' \
  -d 'client_id=admin-cli' | jq -r '.access_token')

if [ "$ADMIN_TOKEN" = "null" ]; then
  echo "Failed to get admin token"
  exit 1
fi

echo "Got admin token successfully."

# Create realm if it doesn't exist
echo "Checking if realm exists..."
REALM_EXISTS=$(curl -s -X GET \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  "$KEYCLOAK_URL/admin/realms/$REALM_NAME" | jq -r '.id')

if [ "$REALM_EXISTS" = "null" ]; then
  echo "Creating realm: $REALM_NAME"
  
  # Create realm
  curl -s -X POST \
    "$KEYCLOAK_URL/admin/realms" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '{
      "id": "'"$REALM_NAME"'",
      "realm": "'"$REALM_NAME"'",
      "enabled": true,
      "loginTheme": "attendance-management",
      "displayName": "Attendance Management System",
      "displayNameHtml": "<div class=\"kc-logo-text\"><span>Attendance Management</span></div>",
      "users": [],
      "roles": {
        "realm": []
      }
    }'
    
  echo "Realm created successfully."
else
  echo "Realm $REALM_NAME already exists."
fi

# Create client if it doesn't exist
echo "Checking if client exists..."
CLIENT_EXISTS=$(curl -s -X GET \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients?clientId=attendance-app" | jq -r '.[0].id')

if [ "$CLIENT_EXISTS" = "null" ]; then
  echo "Creating client: attendance-app"
  
  # Create client
  CLIENT_ID=$(curl -s -X POST \
    "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '{
      "clientId": "attendance-app",
      "enabled": true,
      "redirectUris": ["http://localhost:4200/*"],
      "webOrigins": ["http://localhost:4200"],
      "publicClient": true,
      "protocol": "openid-connect"
    }' | jq -r '.id')
    
  echo "Client created successfully with ID: $CLIENT_ID."
else
  echo "Client attendance-app already exists."
fi

# Create admin user if it doesn't exist
echo "Checking if admin user exists..."
USER_EXISTS=$(curl -s -X GET \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users?username=$ADMIN_USERNAME" | jq -r '.[0].id')

if [ "$USER_EXISTS" = "null" ]; then
  echo "Creating admin user: $ADMIN_USERNAME"
  
  # Create user
  USER_ID=$(curl -s -X POST \
    "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '{
      "username": "'"$ADMIN_USERNAME"'",
      "email": "'"$ADMIN_EMAIL"'",
      "enabled": true,
      "emailVerified": true,
      "firstName": "Admin",
      "lastName": "User"
    }' | jq -r '.id')
    
  echo "User created successfully with ID: $USER_ID."
  
  # Set user password
  curl -s -X PUT \
    "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users/$USER_ID/reset-password" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '{
      "type": "password",
      "value": "'"$ADMIN_PASSWORD"'",
      "temporary": false
    }'
    
  echo "User password set successfully."
  
  # Assign admin role to user
  curl -s -X POST \
    "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users/$USER_ID/role-mappings/realm" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '[{
      "id": "admin",
      "name": "admin"
    }]'
    
  echo "Admin role assigned to user successfully."
else
  echo "Admin user already exists."
fi

echo "Keycloak initialization completed successfully!"