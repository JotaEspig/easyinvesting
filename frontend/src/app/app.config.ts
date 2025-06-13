import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';

import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    provideHttpClient(),
  ],
};

export var backendHost: string = "localhost:8000"; // Default backend host
// check ENV variable for backend host
if (typeof process !== 'undefined' && process.env && process.env["EASY_INVESTING_A_BACKEND_HOST"]) {
  backendHost = process.env["EASY_INVESTING_A_BACKEND_HOST"];
}
export const apiUrl: string = `http://${backendHost}/api/v1/`;
