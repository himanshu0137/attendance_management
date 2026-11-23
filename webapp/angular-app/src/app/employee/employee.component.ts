import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-employee',
  standalone: true,
  template: `
    <div class="employee-container">
      <h2>Employee Dashboard</h2>
      
      <div class="dashboard-actions">
        <button (click)="punchIn()" class="btn btn-primary">Punch In</button>
        <button (click)="punchOut()" class="btn btn-secondary">Punch Out</button>
      </div>
      
      <div class="attendance-history">
        <h3>Recent Attendance Records</h3>
        <div *ngIf="loading">Loading...</div>
        <div *ngIf="!loading && records.length === 0">No attendance records found.</div>
        <table *ngIf="!loading && records.length > 0" class="records-table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Check In</th>
              <th>Check Out</th>
              <th>Duration</th>
            </tr>
          </thead>
          <tbody>
            <tr *ngFor="let record of records">
              <td>{{ record.date }}</td>
              <td>{{ record.checkIn }}</td>
              <td>{{ record.checkOut }}</td>
              <td>{{ record.duration }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  `,
  styles: [`
    .employee-container {
      padding: 20px;
    }
    
    .dashboard-actions {
      margin-bottom: 30px;
    }
    
    .btn {
      padding: 12px 24px;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 16px;
      margin-right: 10px;
    }
    
    .btn-primary {
      background-color: #28a745;
      color: white;
    }
    
    .btn-secondary {
      background-color: #6c757d;
      color: white;
    }
    
    .records-table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 20px;
    }
    
    .records-table th,
    .records-table td {
      padding: 12px;
      text-align: left;
      border-bottom: 1px solid #ddd;
    }
    
    .records-table th {
      background-color: #f8f9fa;
      font-weight: bold;
    }
    
    .records-table tr:hover {
      background-color: #f5f5f5;
    }
  `]
})
export class EmployeeComponent {
  records: any[] = [];
  loading = false;
  
  private router = inject(Router);
  private http = inject(HttpClient);
  
  constructor() {
    this.loadAttendanceRecords();
  }
  
  loadAttendanceRecords() {
    this.loading = true;
    // In a real application, we would call our backend API
    // this.http.get('/api/attendance/reports').subscribe(...)
    
    // Simulate API response with sample data
    setTimeout(() => {
      this.records = [
        { date: '2026-11-23', checkIn: '09:00 AM', checkOut: '05:00 PM', duration: '8h' },
        { date: '2026-11-22', checkIn: '09:15 AM', checkOut: '05:30 PM', duration: '8h 15m' },
        { date: '2026-11-21', checkIn: '08:45 AM', checkOut: '05:15 PM', duration: '8h 30m' }
      ];
      this.loading = false;
    }, 500);
  }
  
  punchIn() {
    console.log('Punch In action');
    // In a real app, this would call the backend API
    alert('Punch In recorded!');
  }
  
  punchOut() {
    console.log('Punch Out action');
    // In a real app, this would call the backend API
    alert('Punch Out recorded!');
  }
}