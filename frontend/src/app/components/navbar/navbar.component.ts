import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-navbar',
  imports: [RouterLink, CommonModule],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css'
})
export class NavbarComponent {
  constructor(private router: Router, public auth: AuthService) {}

  askForLogout(): void {
    if (confirm('Are you sure you want to log out?')) {
      this.auth.logout();
      this.router.navigate(['/login']);
    }
  }
}
