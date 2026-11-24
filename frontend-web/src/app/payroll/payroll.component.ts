import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-payroll',
  template: `
    <div class="payroll-page">
      <header class="page-header">
        <h1>Payroll Processing</h1>
        <p>Manage salary records and generate payroll reports</p>
      </header>

      <div class="content-layout">
        <!-- Sidebar Navigation -->
        <nav class="sidebar-nav">
          <ul>
            <li><a routerLink="/payroll/salary-records" class="nav-link">Salary Records</a></li>
            <li><a routerLink="/payroll/payroll-summary" class="nav-link">Payroll Summary</a></li>
            <li><a routerLink="/payroll/tax-brackets" class="nav-link">Tax Configuration</a></li>
            <li><a routerLink="/payroll/reports" class="nav-link">Payroll Reports</a></li>
          </ul>
        </nav>

        <!-- Main Content -->
        <main class="main-content">
          <router-outlet></router-outlet>
        </main>
      </div>
    </div>
  `,
  styles: [`
    .payroll-page {
      padding: 24px;
      max-width: 1200px;
      margin: 0 auto;
    }

    .page-header {
      text-align: center;
      margin-bottom: 32px;
    }

    .page-header h1 {
      font-size: 2rem;
      color: #1a365d;
      margin-bottom: 8px;
    }

    .page-header p {
      color: #4a5568;
      font-size: 1.1rem;
    }

    .content-layout {
      display: grid;
      grid-template-columns: 250px 1fr;
      gap: 24px;
    }

    .sidebar-nav ul {
      list-style: none;
      padding: 0;
    }

    .sidebar-nav li {
      margin-bottom: 8px;
    }

    .nav-link {
      display: block;
      padding: 12px 16px;
      color: #4a5568;
      text-decoration: none;
      border-radius: 6px;
      transition: all 0.2s ease;
    }

    .nav-link:hover {
      background-color: #edf2f7;
      color: #3182ce;
    }

    .nav-link.active {
      background-color: #3182ce;
      color: white;
    }

    .main-content {
      background: white;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
    }

    @media (max-width: 768px) {
      .content-layout {
        grid-template-columns: 1fr;
      }
      
      .sidebar-nav {
        order: 2;
      }
      
      .main-content {
        order: 1;
      }
    }
  `]
})
export class PayrollComponent {}