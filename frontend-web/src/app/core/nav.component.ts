import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-nav',
  template: `
    <nav class="navbar">
      <div class="nav-brand">
        <a routerLink="/dashboard">Attendance Management</a>
      </div>
      <ul class="nav-menu">
        <li><a routerLink="/dashboard" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">Dashboard</a></li>
        <li><a routerLink="/organization" routerLinkActive="active">Organization</a></li>
        <li><a routerLink="/attendance" routerLinkActive="active">Attendance</a></li>
        <li><a routerLink="/payroll" routerLinkActive="active">Payroll</a></li>
        <li><a routerLink="/leave" routerLinkActive="active">Leave</a></li>
      </ul>
      <div class="nav-user">
        <span>Admin User</span>
      </div>
    </nav>
  `,
  styles: [`
    .navbar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px 24px;
      background-color: #1a365d;
      color: white;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    .nav-brand a {
      color: white;
      text-decoration: none;
      font-size: 1.2rem;
      font-weight: bold;
    }

    .nav-menu {
      display: flex;
      list-style: none;
      margin: 0;
      padding: 0;
    }

    .nav-menu li {
      margin: 0 12px;
    }

    .nav-menu a {
      color: #e2e8f0;
      text-decoration: none;
      padding: 8px 12px;
      border-radius: 4px;
      transition: all 0.2s ease;
    }

    .nav-menu a:hover {
      background-color: rgba(255,255,255,0.1);
      color: white;
    }

    .nav-menu a.active {
      background-color: #3182ce;
      color: white;
    }

    .nav-user {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    @media (max-width: 768px) {
      .navbar {
        flex-direction: column;
        gap: 12px;
      }
      
      .nav-menu {
        flex-wrap: wrap;
        justify-content: center;
      }
      
      .nav-menu li {
        margin: 4px 8px;
      }
    }
  `]
})
export class NavComponent {}