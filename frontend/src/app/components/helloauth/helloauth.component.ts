import { Component } from '@angular/core';
import { ApiService } from '../../services/api.service';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-helloauth',
  imports: [],
  templateUrl: './helloauth.component.html',
  styleUrl: './helloauth.component.css'
})
export class HelloauthComponent {
  message: string = '';

  constructor(private apiService: ApiService, private router: Router, public auth: AuthService) {}

  ngOnInit(): void {
    if (!this.auth.isAuthenticated()) {
      alert('You are not logged in. Please log in to access this page.');
      this.router.navigate(['/login']);
      return;
    }
    const options: Object = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    };
    this.apiService.getRequest<{message: string, user: {id: string}}>("helloauth", options).subscribe({
      next: (data) => this.message = data.message + ' ' + data.user.id,
      error: (err) => {
        alert('You are not logged in. Please log in to access this page.');
        this.router.navigate(['/login']);
      }
    });
  }
}
