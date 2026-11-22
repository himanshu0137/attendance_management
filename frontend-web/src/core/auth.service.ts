import { HttpClient } from '@angular/common/http';
import { Injectable, Signal, signal } from '@angular/core';
import { Observable } from 'rxjs';

// Credentials sent to the Angular auth proxy on the first login.
export interface LoginCredentials {
  username: string;
  password: string;
}

// The access token returned by POST /api/auth/login.
export interface AuthToken {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

// Keys used to persist the access token, refresh token and expiry.
const ACCESS_KEY = 'attendance_access_token';
const REFRESH_KEY = 'attendance_refresh_token';
const EXPIRY_KEY = 'attendance_expires_in';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private token = signal<string | null>(null);
  private expiry = signal<number | null>(null);

  constructor(private http: HttpClient) {}

  // --- token storage -------------------------------------------------

  getToken(): string | null {
    return this.token();
  }

  getExpiry(): number | null {
    return this.expiry();
  }

  setToken(token: AuthToken) {
    this.token.set(token.access_token);
    this.expiry.set(token.expires_in);
    localStorage.setItem(ACCESS_KEY, token.access_token);
    localStorage.setItem(REFRESH_KEY, token.refresh_token);
    localStorage.setItem(EXPIRY_KEY, String(token.expires_in));
  }

  logout() {
    this.token.set(null);
    this.expiry.set(null);
    localStorage.removeItem(ACCESS_KEY);
    localStorage.removeItem(REFRESH_KEY);
    localStorage.removeItem(EXPIRY_KEY);
  }

  // --- auth proxy ----------------------------------------------------

  login(credentials: LoginCredentials): Observable<AuthToken> {
    return this.http.post<AuthToken>('/api/auth/login', credentials);
  }

  loginAndSave(credentials: LoginCredentials): Promise<AuthToken> {
    return new Promise((resolve, reject) => {
      this.login(credentials).subscribe({
        next: (token) => {
          this.setToken(token);
          resolve(token);
        },
        error: (err) => reject(err),
      });
    });
  }
}
