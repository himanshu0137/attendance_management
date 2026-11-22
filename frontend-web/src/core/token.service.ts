import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

// Token pair returned by POST /api/auth/refresh.
export interface RefreshToken {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

@Injectable({ providedIn: 'root' })
export class TokenService {
  constructor(private http: HttpClient) {}

  refresh(): Observable<RefreshToken> {
    return this.http.post<RefreshToken>('/api/auth/refresh', {});
  }
}
