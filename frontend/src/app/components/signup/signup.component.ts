import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common'
import { ApiService } from '../../services/api.service';
import { Router } from "@angular/router"
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule
  ],
  templateUrl: './signup.component.html',
  styleUrl: './signup.component.css'
})
export class SignupComponent {
  signupForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private apiService: ApiService,
    private router: Router,
    private messageService: MessageService
  ) {
    this.signupForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      confirmPassword: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  isPasswordMismatch(): boolean {
    const password = this.signupForm.get('password')?.value;
    const confirmPassword = this.signupForm.get('confirmPassword')?.value;
    return password !== confirmPassword;
  }

  onSubmit(): void {
    if (this.signupForm.valid && !this.isPasswordMismatch()) {
      const { email, password, _ } = this.signupForm.value;
      const userData = {
        "email": email,
        "password": password
      };
      console.log('Form is valid, submitting data:', userData);

      this.apiService.postRequest("signup", userData).subscribe({
        next: (response) => {
          this.router.navigate(['/login']);
          this.messageService.add({
            severity: 'success',
            summary: 'Signup Successful',
            detail: 'You have successfully signed up. Please log in.'
          });
        },
        error: (error) => {
          console.error('Signup failed:', error);
          let details = 'An error occurred during signup. Please try again.';
          if (error.status === 409) {
            details = 'Email already exists. Please use a different email.';
          }
          this.messageService.add({
            severity: 'error',
            summary: 'Signup Failed',
            detail: details
          });
        }
      });
    } else {
      console.log('Form is invalid');
      this.signupForm.markAllAsTouched();
      const errors = this.signupForm.errors;
      if (errors) {
        console.log('Form errors:', errors);
      } else {
        console.log('Email or password is invalid');
      }
    }
  }
}
