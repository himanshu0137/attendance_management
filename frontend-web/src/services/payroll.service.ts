import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

// Money is the integer-based currency amount returned by a slip.
export interface Money {
  amount: string;
}

// One payable day in a salary slip.
export interface PayDay {
  date: string;
  minutes: number;
}

// Generated monthly salary slip.
export interface SalarySlip {
  id: string;
  tenant_id: string;
  employee_id: string;
  period: string;
  base_salary: Money;
  days: PayDay[];
}

@Injectable({ providedIn: 'root' })
export class PayrollService {
  constructor(private http: HttpClient) {}

  // Generates and persists a prorated monthly salary slip.
  generateSlip(employeeId: string, month: number, year: number): Observable<SalarySlip> {
    return this.http.post<SalarySlip>('/api/payroll/slips/generate', {
      employee_id: employeeId,
      month,
      year,
    });
  }

  // Returns the tenant's generated salary slips.
  getSlips(tenantId: string): Observable<SalarySlip[]> {
    return this.http.get<SalarySlip[]>(`/api/payroll/slips`);
  }

  getMySlips(tenantId: string): Observable<SalarySlip[]> {
    return this.http.get<SalarySlip[]>(`/api/payroll/slips/me`);
  }
}
