import { Injectable } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class AuthService {
  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  }

  logout(): void {
    localStorage.removeItem('token');
  }
}
