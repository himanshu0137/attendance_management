# Backend Integration Example

This file demonstrates how the Angular frontend would communicate with the Go backend.

## API Endpoints Structure

### Authentication
- `POST /auth/login` - User login
- `POST /auth/refresh` - Token refresh  
- `POST /auth/logout` - User logout

### Attendance Management
- `POST /api/attendance/punch` - Record punch-in/out
- `GET /api/attendance/me` - Get current employee's records
- `GET /api/attendance/reports` - Get attendance reports

### Organization Management  
- `GET /api/organization/employees` - List employees (admin)
- `POST /api/organization/employees` - Create employee (admin)

## Example HTTP Requests from Angular

### Authentication Service Implementation

```typescript
@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:8080/api';
  
  constructor(private http: HttpClient) {}
  
  login(credentials: { username: string, password: string }) {
    return this.http.post<{access_token: string, refresh_token: string}>(
      `${this.apiUrl}/auth/login`, 
      credentials
    );
  }
  
  refreshToken() {
    // Implementation for token refresh
  }
  
  logout() {
    // Implementation for logout
  }
}
```

### Attendance Service Implementation

```typescript
@Injectable({
  providedIn: 'root'
})
export class AttendanceService {
  private apiUrl = 'http://localhost:8080/api';
  
  constructor(private http: HttpClient) {}
  
  punch(punchData: { type: 'in' | 'out' }) {
    return this.http.post(`${this.apiUrl}/attendance/punch`, punchData);
  }
  
  getMyRecords() {
    return this.http.get<AttendanceRecord[]>(`${this.apiUrl}/attendance/me`);
  }
  
  getReports() {
    return this.http.get<AttendanceReport[]>(`${this.apiUrl}/attendance/reports`);
  }
}
```

## Token Management

The frontend should implement proper token management:

1. Store access token in memory or secure cookies
2. Add Authorization header to all requests:
   ```
   Authorization: Bearer <access_token>
   ```

3. Implement automatic refresh when token expires
4. Handle token expiration gracefully with redirect to login

## Error Handling

All API calls should handle common errors:

- 401 Unauthorized - Redirect to login
- 403 Forbidden - Show access denied message  
- 500 Internal Server Error - Show generic error message