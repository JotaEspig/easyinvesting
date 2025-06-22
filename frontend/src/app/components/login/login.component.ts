import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common'
import { ApiService } from '../../services/api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  loginForm: FormGroup;

  constructor(private fb: FormBuilder, private apiService: ApiService, private router: Router) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  onSubmit(): void {
    if (this.loginForm.valid) {
      const { email, password } = this.loginForm.value;
      const userData = {
        "email": email,
        "password": password
      };

      console.log('Form is valid, submitting data:', userData);

      this.apiService.postRequest<{token: string}>("login", userData).subscribe({
        next: (data) => {
          console.log('Login successful:', data);
          localStorage.setItem('token', data.token);
          this.router.navigate(['/helloauth']); // TODO just temporary
        },
        error: (error) => {
          alert('Login failed');
          console.error('Login error:', error);
          this.loginForm.reset();
        }
      });
    } else {
      console.log('Form is invalid');
      this.loginForm.markAllAsTouched();
      // tell why the form is invalid
      const errors = this.loginForm.errors;
      if (errors) {
        console.log('Form errors:', errors);
      } else {
        console.log('Email or password is invalid');
      }
    }
  }
}
