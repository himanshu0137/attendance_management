import { bootstrapApplication } from '@angular/platform-browser';
import { provideBrowserGlobalErrorListeners } from '@angular/core';

import { routes } from './app/app.routes';
import { AppComponent } from './app/app';
import { provideRouter } from '@angular/router';

bootstrapApplication(
  AppComponent,
  provideRouter(routes),
  provideBrowserGlobalErrorListeners()
);
