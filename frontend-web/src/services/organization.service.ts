import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ApiBaseService } from './api-base.service';
import { Employee, CreateEmployeeRequest, EmployeeWithDocuments, OnboardEmployeeReq } from '../shared/models/employee.model';
import { SalaryStructure, CreateSalaryStructureRequest } from '../shared/models/salary.model';
import { HRLetterResponse } from '../shared/models/hr-letter.model';

@Injectable({
  providedIn: 'root'
})
export class OrganizationService extends ApiBaseService {
  
  constructor(protected http: HttpClient) {
    super(http);
  }

  // Employee Management
  createEmployee(tenantId: string, employeeData: CreateEmployeeRequest): Observable<Employee> {
    return this.http.post<Employee>(`${this.baseUrl}/organization/${tenantId}/employees`, employeeData);
  }

  getEmployee(tenantId: string, employeeId: string): Observable<Employee> {
    return this.http.get<Employee>(`${this.baseUrl}/organization/${tenantId}/employees/${employeeId}`);
  }

  // Employee Onboarding
  onboardEmployee(tenantId: string, employeeId: string, onboardingData: OnboardEmployeeReq): Observable<any> {
    return this.http.post(`${this.baseUrl}/organization/${tenantId}/employees/${employeeId}/onboard`, onboardingData);
  }

  getEmployeeDocuments(tenantId: string, employeeId: string): Observable<EmployeeWithDocuments> {
    return this.http.get<EmployeeWithDocuments>(`${this.baseUrl}/organization/${tenantId}/employees/${employeeId}/documents`);
  }

  // Salary Structure
  createSalaryStructure(tenantId: string, employeeId: string, salaryData: CreateSalaryStructureRequest): Observable<SalaryStructure> {
    return this.http.post<SalaryStructure>(`${this.baseUrl}/organization/${tenantId}/employees/${employeeId}/salary-structure`, salaryData);
  }

  getSalaryStructure(tenantId: string, employeeId: string): Observable<SalaryStructure> {
    return this.http.get<SalaryStructure>(`${this.baseUrl}/organization/${tenantId}/employees/${employeeId}/salary-structure`);
  }

  updateSalaryStructure(tenantId: string, employeeId: string, salaryData: CreateSalaryStructureRequest): Observable<SalaryStructure> {
    return this.http.put<SalaryStructure>(`${this.baseUrl}/organization/${tenantId}/employees/${employeeId}/salary-structure`, salaryData);
  }

  // HR Letters
  generateHRLetter(tenantId: string, employeeId: string, letterType: string): Observable<HRLetterResponse> {
    return this.http.post<HRLetterResponse>(`${this.baseUrl}/organization/${tenantId}/employees/${employeeId}/hr-letters/${letterType}`, {});
  }
}