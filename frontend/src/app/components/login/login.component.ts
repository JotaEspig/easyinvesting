import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common'
import { ApiService } from '../../services/api.service';
import { Router } from '@angular/router';
import { MessageService } from 'primeng/api';

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

  constructor(
    private fb: FormBuilder,
    private apiService: ApiService,
    private router: Router,
    private messageService: MessageService
  ) {
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
          localStorage.setItem('token', data.token);
          this.router.navigate(['/portfolio']);
          this.messageService.add({
            severity: 'success',
            summary: 'Login Successful',
            detail: 'You have successfully logged in'
          });
        },
        error: (error) => {
          this.messageService.add({
            severity: 'error',
            summary: 'Login Failed',
            detail: 'Invalid email or password. Please try again.'
          });
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
