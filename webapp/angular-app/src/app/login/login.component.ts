import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule],
  template: `
    <div class="login-container">
      <div class="login-card">
        <h2>Login to Attendance Portal</h2>
        <form (ngSubmit)="onLogin()" #loginForm="ngForm" class="login-form">
          <div class="form-group">
            <label for="username">Username:</label>
            <input 
              type="text" 
              id="username" 
              name="username" 
              [(ngModel)]="credentials.username" 
              required
              class="form-control">
          </div>
          
          <div class="form-group">
            <label for="password">Password:</label>
            <input 
              type="password" 
              id="password" 
              name="password" 
              [(ngModel)]="credentials.password" 
              required
              class="form-control">
          </div>
          
          <button type="submit" [disabled]="isLoading" class="btn btn-primary">
            {{ isLoading ? 'Logging in...' : 'Login' }}
          </button>
        </form>
        
        <div *ngIf="error" class="error-message">{{ error }}</div>
      </div>
    </div>
  `,
  styles: [`
    .login-container {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 100vh;
      background-color: #f5f5f5;
    }
    
    .login-card {
      background: white;
      padding: 30px;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0,0,0,0.1);
      width: 100%;
      max-width: 400px;
    }
    
    .login-form {
      display: flex;
      flex-direction: column;
    }
    
    .form-group {
      margin-bottom: 20px;
    }
    
    label {
      display: block;
      margin-bottom: 5px;
      font-weight: bold;
    }
    
    .form-control {
      width: 100%;
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 4px;
      box-sizing: border-box;
    }
    
    .btn {
      padding: 12px;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 16px;
    }
    
    .btn-primary {
      background-color: #007bff;
      color: white;
    }
    
    .btn-primary:disabled {
      background-color: #6c757d;
    }
    
    .error-message {
      margin-top: 15px;
      padding: 10px;
      background-color: #f8d7da;
      color: #721c24;
      border-radius: 4px;
      border: 1px solid #f5c6cb;
    }
  `]
})
export class LoginComponent {
  credentials = { username: '', password: '' };
  isLoading = false;
  error = '';
  
  private router = inject(Router);
  private http = inject(HttpClient);
  
  onLogin() {
    this.isLoading = true;
    this.error = '';
    
    // In a real application, we would call our backend authentication endpoint
    console.log('Attempting login with:', this.credentials);
    
    // Simulate API call to backend
    setTimeout(() => {
      this.isLoading = false;
      // For demo purposes, redirect regardless of credentials
      this.router.navigate(['/employee']);
    }, 1000);
  }
}