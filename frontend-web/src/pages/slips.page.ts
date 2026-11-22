import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormControl,
  FormGroup,
  FormsModule,
} from '@angular/forms';
import {
  MatButton,
  MatButtonModule,
  MatFormField,
  MatFormFieldModule,
  MatInput,
  MatInputModule,
  MatSelect,
  MatSelectModule,
} from '@angular/material';

import { AuthService } from '../core/auth.service';
import { PayrollService, SalarySlip } from '../services/payroll.service';
import { MatTableComponent, TableColumn } from '../shared/mat-table.component';

@Component({
  selector: 'app-slips',
  standalone: true,
  imports: [CommonModule, FormsModule, MatButtonModule, MatFormFieldModule, MatInputModule, MatSelectModule, MatTableComponent],
  templateUrl: './slips.html',
  styleUrl: './slips.scss'
})
export class SlipsPage implements OnInit {
  employeeId = '';
  months = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12];
  month = 1;
  year = new Date().getFullYear();
  loading = false;
  error = '';
  slips: SalarySlip[] = [];

  columns: TableColumn[] = [
    { key: 'id', label: 'ID' },
    { key: 'period', label: 'Period' },
    { key: 'base_salary', label: 'Base salary' },
    { key: 'days', label: 'Days' },
  ];

  form = {
    month: new FormControl(1),
    year: new FormControl(new Date().getFullYear()),
  } as FormGroup;

  constructor(private auth: AuthService, private payroll: PayrollService) {}

  ngOnInit() {
    this.employeeId = this.auth.getToken() || '';
    this.year = new Date().getFullYear();
    this.load();
  }

  generate() {
    this.error = '';
    this.loading = true;
    this.payroll
      .generateSlip(this.employeeId, this.month, this.year)
      .subscribe({
        next: (slip) => {
          this.slips = [slip, ...this.slips];
          this.loading = false;
        },
        error: (err) => {
          this.error = err.error?.detail || 'slip generation failed';
          this.loading = false;
        },
      });
  }

  load() {
    if (!this.employeeId) return;
    this.payroll
      .getMySlips(this.employeeId)
      .subscribe((slips) => {
        this.slips = slips;
      });
  }
}
