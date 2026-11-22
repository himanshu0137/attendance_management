import { Routes } from '@angular/router';
import { AuthService } from '../core/auth.service';
import { TokenService } from '../core/token.service';
import { canActivate } from '../core/guards';
import { LoginComponent } from '../core/login.component';
import { PunchPage } from '../pages/punch.page';
import { ReportsPage } from '../pages/reports.page';
import { SlipsPage } from '../pages/slips.page';

export const routes: Routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  {
    path: 'punch',
    component: PunchPage,
    canActivate: [canActivate(tokenService, auth)],
    data: { title: 'Punch' },
  },
  {
    path: 'reports',
    component: ReportsPage,
    canActivate: [canActivate(tokenService, auth)],
    data: { title: 'Reports' },
  },
  {
    path: 'slips',
    component: SlipsPage,
    canActivate: [canActivate(tokenService, auth)],
    data: { title: 'Slips' },
  },
];
