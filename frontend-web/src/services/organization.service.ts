import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

// OrgPolicy is the per-tenant attendance policy.
export interface OrgPolicy {
  methods: string[];
  salary_slip_enabled: boolean;
  shift_start_time: string;
  shift_end_time: string;
}

// Employee is a Keycloak-backed org user.
export interface Employee {
  id: string;
  keycloak_sub: string;
  role: string;
}

@Injectable({ providedIn: 'root' })
export class OrganizationService {
  constructor(private http: HttpClient) {}

  getPolicy(tenantId: string): Observable<OrgPolicy> {
    return this.http.get<OrgPolicy>(`/api/organization/${tenantId}/policies`);
  }

  setPolicy(tenantId: string, policy: OrgPolicy): Observable<OrgPolicy> {
    return this.http.post<OrgPolicy>(`/api/organization/${tenantId}/policies`, policy);
  }

  listEmployees(tenantId: string): Observable<Employee[]> {
    return this.http.get<Employee[]>(`/api/organization/${tenantId}/employees`);
  }

  createEmployee(tenantId: string, username: string, role: string): Observable<{ id: string }> {
    return this.http.post<{ id: string }>(
      `/api/organization/${tenantId}/employees`,
      { username, role }
    );
  }
}
