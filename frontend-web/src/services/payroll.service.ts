import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ApiBaseService } from './api-base.service';
import { StatutoryCompliance, EmployeeStatutoryDeduction } from '../shared/models/statutory.model';

@Injectable({
  providedIn: 'root'
})
export class PayrollService extends ApiBaseService {
  
  constructor(protected http: HttpClient) {
    super(http);
  }

  // Statutory Compliance
  createStatutoryCompliance(tenantId: string, complianceData: StatutoryCompliance): Observable<StatutoryCompliance> {
    return this.http.post<StatutoryCompliance>(`${this.baseUrl}/payroll/statutory`, complianceData);
  }

  getStatutoryCompliance(tenantId: string, complianceId: string): Observable<StatutoryCompliance> {
    return this.http.get<StatutoryCompliance>(`${this.baseUrl}/payroll/statutory/${complianceId}`);
  }

  listStatutoryCompliance(tenantId: string, page: number = 1, limit: number = 20): Observable<{ compliances: StatutoryCompliance[], total: number }> {
    return this.http.get<{ compliances: StatutoryCompliance[], total: number }>(
      `${this.baseUrl}/payroll/statutory?page=${page}&limit=${limit}`
    );
  }

  updateStatutoryCompliance(tenantId: string, complianceId: string, complianceData: StatutoryCompliance): Observable<StatutoryCompliance> {
    return this.http.put<StatutoryCompliance>(`${this.baseUrl}/payroll/statutory/${complianceId}`, complianceData);
  }

  deleteStatutoryCompliance(tenantId: string, complianceId: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/payroll/statutory/${complianceId}`);
  }

  // Employee Statutory Deductions
  createEmployeeStatutoryDeduction(employeeId: string, deductionData: EmployeeStatutoryDeduction): Observable<EmployeeStatutoryDeduction> {
    return this.http.post<EmployeeStatutoryDeduction>(`${this.baseUrl}/payroll/employee/${employeeId}/statutory-deductions`, deductionData);
  }

  getEmployeeStatutoryDeduction(employeeId: string, deductionId: string): Observable<EmployeeStatutoryDeduction> {
    return this.http.get<EmployeeStatutoryDeduction>(`${this.baseUrl}/payroll/employee/${employeeId}/statutory-deductions/${deductionId}`);
  }

  listEmployeeStatutoryDeductions(employeeId: string, page: number = 1, limit: number = 20): Observable<{ deductions: EmployeeStatutoryDeduction[], total: number }> {
    return this.http.get<{ deductions: EmployeeStatutoryDeduction[], total: number }>(
      `${this.baseUrl}/payroll/employee/${employeeId}/statutory-deductions?page=${page}&limit=${limit}`
    );
  }

  updateEmployeeStatutoryDeduction(employeeId: string, deductionId: string, deductionData: EmployeeStatutoryDeduction): Observable<EmployeeStatutoryDeduction> {
    return this.http.put<EmployeeStatutoryDeduction>(`${this.baseUrl}/payroll/employee/${employeeId}/statutory-deductions/${deductionId}`, deductionData);
  }

  deleteEmployeeStatutoryDeduction(employeeId: string, deductionId: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/payroll/employee/${employeeId}/statutory-deductions/${deductionId}`);
  }
}