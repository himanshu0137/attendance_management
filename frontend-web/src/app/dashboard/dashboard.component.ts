import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  template: `
    <div class="dashboard">
      <header class="dashboard-header">
        <h1>Attendance Management Dashboard</h1>
        <p>Welcome to the comprehensive workforce management system</p>
      </header>

      <div class="features-grid">
        <!-- Organization Management Card -->
        <div class="feature-card">
          <div class="card-icon">
            <i class="fas fa-users"></i>
          </div>
          <h2>Organization Management</h2>
          <p>Employee management, policy settings, and salary structure definition</p>
          <ul class="feature-list">
            <li>Employee onboarding workflow</li>
            <li>HR letter generation</li>
            <li>Organization policy configuration</li>
          </ul>
          <a routerLink="/organization" class="btn btn-primary">Manage Organization</a>
        </div>

        <!-- Attendance Tracking Card -->
        <div class="feature-card">
          <div class="card-icon">
            <i class="fas fa-clock"></i>
          </div>
          <h2>Attendance Tracking</h2>
          <p>Real-time punch tracking and attendance reporting</p>
          <ul class="feature-list">
            <li>Punch-in/punch-out functionality</li>
            <li>Location-based tracking</li>
            <li>Overtime management</li>
          </ul>
          <a routerLink="/attendance" class="btn btn-primary">Track Attendance</a>
        </div>

        <!-- Payroll Processing Card -->
        <div class="feature-card">
          <div class="card-icon">
            <i class="fas fa-money-bill-wave"></i>
          </div>
          <h2>Payroll Processing</h2>
          <p>Salary calculations and payroll management</p>
          <ul class="feature-list">
            <li>Salary record creation</li>
            <li>Net salary calculation</li>
            <li>Payment tracking</li>
          </ul>
          <a routerLink="/payroll" class="btn btn-primary">Process Payroll</a>
        </div>

        <!-- Leave Management Card -->
        <div class="feature-card">
          <div class="card-icon">
            <i class="fas fa-calendar-alt"></i>
          </div>
          <h2>Leave Management</h2>
          <p>Employee leave requests and approvals</p>
          <ul class="feature-list">
            <li>Leave request submission</li>
            <li>Approval workflow</li>
            <li>Leave balance tracking</li>
          </ul>
          <a routerLink="/leave" class="btn btn-primary">Manage Leaves</a>
        </div>
      </div>

      <!-- System Overview -->
      <section class="system-overview">
        <h2>System Overview</h2>
        <div class="stats-grid">
          <div class="stat-card">
            <h3>Active Employees</h3>
            <p class="stat-value">1,248</p>
          </div>
          <div class="stat-card">
            <h3>This Month's Attendance</h3>
            <p class="stat-value">98.7%</p>
          </div>
          <div class="stat-card">
            <h3>Pending Leaves</h3>
            <p class="stat-value">24</p>
          </div>
          <div class="stat-card">
            <h3>Payroll Processed</h3>
            <p class="stat-value">100%</p>
          </div>
        </div>
      </section>

      <!-- Recent Activity -->
      <section class="recent-activity">
        <h2>Recent Activity</h2>
        <div class="activity-list">
          <div class="activity-item">
            <div class="activity-icon bg-blue">
              <i class="fas fa-user-plus"></i>
            </div>
            <div class="activity-content">
              <h4>New Employee Onboarded</h4>
              <p>John Smith joined the Marketing department</p>
              <span class="activity-time">2 hours ago</span>
            </div>
          </div>
          <div class="activity-item">
            <div class="activity-icon bg-green">
              <i class="fas fa-clock"></i>
            </div>
            <div class="activity-content">
              <h4>Attendance Recorded</h4>
              <p>Sarah Johnson punched in at 9:05 AM</p>
              <span class="activity-time">4 hours ago</span>
            </div>
          </div>
          <div class="activity-item">
            <div class="activity-icon bg-purple">
              <i class="fas fa-file-invoice-dollar"></i>
            </div>
            <div class="activity-content">
              <h4>Payroll Processed</h4>
              <p>Monthly payroll completed for 1,248 employees</p>
              <span class="activity-time">1 day ago</span>
            </div>
          </div>
        </div>
      </section>
    </div>
  `,
  styles: [`
    .dashboard {
      padding: 24px;
      max-width: 1200px;
      margin: 0 auto;
    }

    .dashboard-header {
      text-align: center;
      margin-bottom: 40px;
    }

    .dashboard-header h1 {
      font-size: 2.5rem;
      color: #1a365d;
      margin-bottom: 12px;
    }

    .dashboard-header p {
      font-size: 1.1rem;
      color: #4a5568;
      max-width: 600px;
      margin: 0 auto;
    }

    .features-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 24px;
      margin-bottom: 40px;
    }

    .feature-card {
      background: white;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      transition: transform 0.3s ease, box-shadow 0.3s ease;
      border: 1px solid #e2e8f0;
    }

    .feature-card:hover {
      transform: translateY(-5px);
      box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
    }

    .card-icon {
      font-size: 2.5rem;
      color: #3182ce;
      margin-bottom: 16px;
    }

    .feature-card h2 {
      font-size: 1.5rem;
      color: #2d3748;
      margin-bottom: 12px;
    }

    .feature-card p {
      color: #4a5568;
      margin-bottom: 16px;
      line-height: 1.5;
    }

    .feature-list {
      list-style-type: none;
      padding: 0;
      margin: 16px 0;
    }

    .feature-list li {
      padding: 8px 0;
      border-bottom: 1px solid #edf2f7;
      display: flex;
      align-items: center;
    }

    .feature-list li:before {
      content: "✓";
      color: #38a169;
      font-weight: bold;
      margin-right: 8px;
    }

    .btn {
      display: inline-block;
      padding: 10px 20px;
      background-color: #3182ce;
      color: white;
      border-radius: 6px;
      text-decoration: none;
      font-weight: 500;
      transition: background-color 0.3s ease;
      border: none;
      cursor: pointer;
    }

    .btn:hover {
      background-color: #2b6cb0;
    }

    .btn-primary {
      background-color: #3182ce;
    }

    .btn-primary:hover {
      background-color: #2b6cb0;
    }

    .system-overview {
      background: white;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      margin-bottom: 40px;
    }

    .system-overview h2 {
      color: #2d3748;
      margin-bottom: 24px;
      text-align: center;
    }

    .stats-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 16px;
    }

    .stat-card {
      background: #f7fafc;
      padding: 20px;
      border-radius: 8px;
      text-align: center;
      border: 1px solid #e2e8f0;
    }

    .stat-card h3 {
      font-size: 1rem;
      color: #4a5568;
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 1.8rem;
      font-weight: bold;
      color: #2d3748;
    }

    .recent-activity {
      background: white;
      border-radius: 12px;
      padding: 24px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
    }

    .recent-activity h2 {
      color: #2d3748;
      margin-bottom: 24px;
    }

    .activity-list {
      display: flex;
      flex-direction: column;
      gap: 16px;
    }

    .activity-item {
      display: flex;
      align-items: flex-start;
      gap: 16px;
    }

    .activity-icon {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 18px;
    }

    .bg-blue {
      background-color: #3182ce;
    }

    .bg-green {
      background-color: #38a169;
    }

    .bg-purple {
      background-color: #805ad5;
    }

    .activity-content {
      flex: 1;
    }

    .activity-content h4 {
      margin: 0 0 8px 0;
      color: #2d3748;
    }

    .activity-content p {
      margin: 0 0 8px 0;
      color: #4a5568;
      font-size: 0.9rem;
    }

    .activity-time {
      color: #718096;
      font-size: 0.8rem;
    }

    @media (max-width: 768px) {
      .dashboard {
        padding: 16px;
      }
      
      .dashboard-header h1 {
        font-size: 2rem;
      }
      
      .features-grid {
        grid-template-columns: 1fr;
      }
      
      .stats-grid {
        grid-template-columns: 1fr 1fr;
      }
    }

    @media (max-width: 480px) {
      .stats-grid {
        grid-template-columns: 1fr;
      }
      
      .dashboard-header h1 {
        font-size: 1.5rem;
      }
      
      .activity-item {
        flex-direction: column;
        gap: 8px;
      }
    }
  `]
})
export class DashboardComponent {}