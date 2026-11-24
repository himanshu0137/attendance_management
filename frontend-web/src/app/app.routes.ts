import { Routes } from '@angular/router';
import { AuthService } from '../core/auth.service';
import { TokenService } from '../core/token.service';
import { canActivate } from '../core/guards';
import { LoginComponent } from '../core/login.component';
import { PunchPage } from '../pages/punch.page';
import { ReportsPage } from '../pages/reports.page';
import { SlipsPage } from '../pages/slips.page';
import { DashboardComponent } from './dashboard/dashboard.component';
import { OrganizationComponent } from './organization/organization.component';
import { AttendanceComponent } from './attendance/attendance.component';
import { PayrollComponent } from './payroll/payroll.component';
import { LeaveComponent } from './leave/leave.component';

export const routes: Routes = [
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'dashboard', component: DashboardComponent, canActivate: [canActivate(TokenService, AuthService)] },
  {
    path: 'organization',
    component: OrganizationComponent,
    canActivate: [canActivate(TokenService, AuthService)],
    children: [
      { path: '', redirectTo: 'employees', pathMatch: 'full' },
      { path: 'employees', component: DashboardComponent },
      { path: 'policies', component: DashboardComponent },
      { path: 'salary-structures', component: DashboardComponent },
      { path: 'onboarding', component: DashboardComponent },
      { path: 'hr-letters', component: DashboardComponent },
    ]
  },
  {
    path: 'attendance',
    component: AttendanceComponent,
    canActivate: [canActivate(TokenService, AuthService)],
    children: [
      { path: '', redirectTo: 'punch', pathMatch: 'full' },
      { path: 'punch', component: DashboardComponent },
      { path: 'shifts', component: DashboardComponent },
      { path: 'reports', component: DashboardComponent },
      { path: 'overtime', component: DashboardComponent },
    ]
  },
  {
    path: 'payroll',
    component: PayrollComponent,
    canActivate: [canActivate(TokenService, AuthService)],
    children: [
      { path: '', redirectTo: 'salary-records', pathMatch: 'full' },
      { path: 'salary-records', component: DashboardComponent },
      { path: 'payroll-summary', component: DashboardComponent },
      { path: 'tax-brackets', component: DashboardComponent },
      { path: 'reports', component: DashboardComponent },
    ]
  },
  {
    path: 'leave',
    component: LeaveComponent,
    canActivate: [canActivate(TokenService, AuthService)],
    children: [
      { path: '', redirectTo: 'requests', pathMatch: 'full' },
      { path: 'requests', component: DashboardComponent },
      { path: 'types', component: DashboardComponent },
      { path: 'balances', component: DashboardComponent },
      { path: 'reports', component: DashboardComponent },
    ]
  },
  {
    path: 'punch',
    component: PunchPage,
    canActivate: [canActivate(TokenService, AuthService)],
    data: { title: 'Punch' },
  },
  {
    path: 'reports',
    component: ReportsPage,
    canActivate: [canActivate(TokenService, AuthService)],
    data: { title: 'Reports' },
  },
  {
    path: 'slips',
    component: SlipsPage,
    canActivate: [canActivate(TokenService, AuthService)],
    data: { title: 'Slips' },
  },
];
