import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ApiBaseService } from './api-base.service';
import { Shift } from '../shared/models/shift.model';
import { AttendanceRule } from '../shared/models/attendance-rule.model';

@Injectable({
  providedIn: 'root'
})
export class AttendanceService extends ApiBaseService {
  
  constructor(protected http: HttpClient) {
    super(http);
  }

  // Shift Management
  createShift(tenantId: string, shiftData: Shift): Observable<Shift> {
    return this.http.post<Shift>(`${this.baseUrl}/attendance/shifts`, shiftData);
  }

  getShift(tenantId: string, shiftId: string): Observable<Shift> {
    return this.http.get<Shift>(`${this.baseUrl}/attendance/shifts/${shiftId}`);
  }

  listShifts(tenantId: string, page: number = 1, limit: number = 20): Observable<{ shifts: Shift[], total: number }> {
    return this.http.get<{ shifts: Shift[], total: number }>(
      `${this.baseUrl}/attendance/shifts?page=${page}&limit=${limit}`
    );
  }

  updateShift(tenantId: string, shiftId: string, shiftData: Shift): Observable<Shift> {
    return this.http.put<Shift>(`${this.baseUrl}/attendance/shifts/${shiftId}`, shiftData);
  }

  deleteShift(tenantId: string, shiftId: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/attendance/shifts/${shiftId}`);
  }

  // Attendance Rules
  createAttendanceRule(tenantId: string, ruleData: AttendanceRule): Observable<AttendanceRule> {
    return this.http.post<AttendanceRule>(`${this.baseUrl}/attendance/rules`, ruleData);
  }

  getAttendanceRule(tenantId: string, ruleId: string): Observable<AttendanceRule> {
    return this.http.get<AttendanceRule>(`${this.baseUrl}/attendance/rules/${ruleId}`);
  }

  listAttendanceRules(tenantId: string, page: number = 1, limit: number = 20): Observable<{ rules: AttendanceRule[], total: number }> {
    return this.http.get<{ rules: AttendanceRule[], total: number }>(
      `${this.baseUrl}/attendance/rules?page=${page}&limit=${limit}`
    );
  }

  updateAttendanceRule(tenantId: string, ruleId: string, ruleData: AttendanceRule): Observable<AttendanceRule> {
    return this.http.put<AttendanceRule>(`${this.baseUrl}/attendance/rules/${ruleId}`, ruleData);
  }

  deleteAttendanceRule(tenantId: string, ruleId: string): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/attendance/rules/${ruleId}`);
  }
}