import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
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
import { AttendanceService, PunchMethod } from '../services/attendance.service';
import { OrganizationService, OrgPolicy } from '../services/organization.service';

@Component({
  selector: 'app-punch',
  standalone: true,
  imports: [CommonModule, FormsModule, MatButtonModule, MatFormFieldModule, MatInputModule, MatSelectModule],
  templateUrl: './punch.html',
  styleUrl: './punch.scss'
})
export class PunchPage implements OnInit {
  employeeId = '';
  policy: OrgPolicy | null = null;
  methods: PunchMethod[] = [];
  method = 'FACESCAN_KIOSK' as PunchMethod;
  startControl = new FormControl('');
  endControl = new FormControl('');
  loading = false;
  error = '';
  logs: PunchLog[] = [];

  constructor(
    private auth: AuthService,
    private attendance: AttendanceService,
    private org: OrganizationService
  ) {}

  ngOnInit() {
    this.employeeId = this.auth.getToken() || '';
    this.loadPolicy();
    this.loadLogs();
  }

  loadPolicy() {
    if (!this.employeeId) return;
    this.org
      .getPolicy(this.employeeId)
      .subscribe((policy) => {
        this.policy = policy;
        this.methods = (policy.methods || []).map((m) => m as PunchMethod);
        if (!this.methods.includes(this.method)) {
          this.method = this.methods[0];
        }
        const s = new Date(policy.shift_start_time);
        const e = new Date(policy.shift_end_time);
        if (!isNaN(s.getTime())) this.startControl.setValue(s.toISOString().slice(0, 16));
        if (!isNaN(e.getTime())) this.endControl.setValue(e.toISOString().slice(0, 16));
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

  punch() {
    this.error = '';
    this.loading = true;
    const ts = this.startControl.value;
    this.attendance
      .punch(this.employeeId, this.method, ts)
      .subscribe({
        next: (log) => {
          this.logs = [log, ...this.logs];
          this.loading = false;
        },
        error: (err) => {
          this.error = err.error?.detail || 'punch failed';
          this.loading = false;
        },
      });
  }
}
