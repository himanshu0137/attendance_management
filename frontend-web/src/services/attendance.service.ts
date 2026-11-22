import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

// PunchMethod is the biometric or login channel used for a punch.
export type PunchMethod =
  | 'FACESCAN_KIOSK'
  | 'FINGERSCAN_KIOSK'
  | 'WEB_LOGIN'
  | 'MOBILE_LOGIN';

// Raw punch event returned by the attendance service.
export interface PunchLog {
  id: number;
  tenant_id: string;
  employee_id: string;
  method: string;
  start: string;
  end: string | null;
}

// Aggregated attendance for a date range.
export interface AttendanceReport {
  present_days: number;
  worked_minutes: number;
}

@Injectable({ providedIn: 'root' })
export class AttendanceService {
  constructor(private http: HttpClient) {}

  // Records a raw IN/OUT event for the employee.
  punch(employeeId: string, method: PunchMethod, timestamp: string): Observable<PunchLog> {
    return this.http.post<PunchLog>('/api/attendance/punch', {
      employee_id: employeeId,
      method,
      timestamp,
    });
  }

  // Returns the tenant's raw punch events.
  getLogs(tenantId: string): Observable<PunchLog[]> {
    return this.http.get<PunchLog[]>(`/api/attendance/me`);
  }

  // Aggregates the tenant's raw punches for a date range.
  report(tenantId: string, dateFrom: string, dateTo: string): Observable<AttendanceReport> {
    return this.http.get<AttendanceReport>(`/api/attendance/reports`, {
      params: { date_from: dateFrom, date_to: dateTo },
    });
  }
}
