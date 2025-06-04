import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common'

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

  constructor(private fb: FormBuilder) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  onSubmit(): void {
    if (this.loginForm.valid) {
      const { email, password } = this.loginForm.value;
      // Handle login logic here
      console.log('Login successful', { email, password });
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
