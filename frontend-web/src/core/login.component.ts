import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormControl, FormGroup } from '@angular/forms';
import {
  MatButton,
  MatButtonModule,
  MatFormField,
  MatFormFieldModule,
  MatIcon,
  MatInput,
  MatInputModule,
} from '@angular/material';

import { AuthService, LoginCredentials } from './auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatButtonModule, MatFormFieldModule, MatInputModule],
  templateUrl: './login.html',
  styleUrl: './login.scss'
})
export class LoginComponent {
  form = new FormGroup({
    username: new FormControl('', { nonNullable: true }),
    password: new FormControl('', { nonNullable: true }),
  });

  loading = false;
  error = '';

  constructor(private auth: AuthService) {}

  login() {
    this.error = '';
    this.loading = true;
    this.auth
      .loginAndSave(this.form.value as LoginCredentials)
      .then(() => {
        this.loading = false;
      })
      .catch((err) => {
        this.error = err.error?.detail || 'login failed';
        this.loading = false;
      });
  }
}
