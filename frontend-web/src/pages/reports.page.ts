import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  MatButton,
  MatButtonModule,
  MatFormField,
  MatFormFieldModule,
  MatIcon,
  MatInput,
  MatInputModule,
  MatSelect,
  MatSelectModule,
} from '@angular/material';

import { AuthService } from '../core/auth.service';
import { DateRange, DateRangePicker } from '../shared/date-range.component';
import { AttendanceService, AttendanceReport, PunchLog } from '../services/attendance.service';
import { OrganizationService, OrgPolicy } from '../services/organization.service';
import { MatTableComponent, TableColumn } from '../shared/mat-table.component';
import { StatusBadge } from '../shared/badge.component';

@Component({
  selector: 'app-reports',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    DateRangePicker,
    MatTableComponent,
    StatusBadge,
  ],
  templateUrl: './reports.html',
  styleUrl: './reports.scss'
})
export class ReportsPage implements OnInit {
  employeeId = '';
  policy: OrgPolicy | null = null;
  range: DateRange = { date_from: '', date_to: '' };
  report: AttendanceReport | null = null;
  logs: PunchLog[] = [];
  loading = false;
  error = '';

  columns: TableColumn[] = [
    { key: 'start', label: 'Start' },
    { key: 'method', label: 'Method' },
    { key: 'end', label: 'End' },
  ];

  constructor(
    private auth: AuthService,
    private attendance: AttendanceService,
    private org: OrganizationService
  ) {}

  ngOnInit() {
    this.employeeId = this.auth.getToken() || '';
    this.loadPolicy();
  }

  loadPolicy() {
    if (!this.employeeId) return;
    this.org
      .getPolicy(this.employeeId)
      .subscribe((policy) => {
        this.policy = policy;
      });
  }

  load() {
    this.loading = true;
    this.error = '';
    this.attendance
      .report(this.employeeId, this.range.date_from, this.range.date_to)
      .subscribe({
        next: (r) => {
          this.report = r;
          this.loading = false;
        },
        error: (err) => {
          this.error = err.error?.detail || 'report failed';
          this.loading = false;
        },
      });
  }

  loadLogs() {
    if (!this.employeeId) return;
    this.attendance
      .getLogs(this.employeeId)
      .subscribe((logs) => {
        this.logs = logs;
      });
  }
}
