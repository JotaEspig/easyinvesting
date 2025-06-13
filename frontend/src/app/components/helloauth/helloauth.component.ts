import { Component } from '@angular/core';
import { ApiService } from '../../services/api.service';

@Component({
  selector: 'app-helloauth',
  imports: [],
  templateUrl: './helloauth.component.html',
  styleUrl: './helloauth.component.css'
})
export class HelloauthComponent {
  message: string = '';

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
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
        window.location.href = '/login';
      }
    });
  }
}
