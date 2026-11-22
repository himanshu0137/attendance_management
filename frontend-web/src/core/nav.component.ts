import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  MatButton,
  MatButtonModule,
  MatButtonToggle,
  MatButtonToggleModule,
  MatIcon,
  MatIconButton,
  MatIconModule,
  MatToolbar,
  MatToolbarModule,
  MatTooltip,
  MatTooltipModule,
} from '@angular/material';

import { AuthService } from './auth.service';

@Component({
  selector: 'app-nav',
  standalone: true,
  imports: [
    CommonModule,
    MatToolbarModule,
    MatButtonModule,
    MatButtonToggleModule,
    MatIconModule,
    MatTooltipModule,
  ],
  templateUrl: './nav.html',
  styleUrl: './nav.scss'
})
export class NavComponent {
  logout() {
    this.auth.logout();
    window.location.href = '/login';
  }
}
