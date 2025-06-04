import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common'

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

  constructor(private fb: FormBuilder) {
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
      const { email, password, confirmPassword } = this.signupForm.value;
      // Handle signup logic here
      console.log('Signup successful', { email, password, confirmPassword });
    } else {
      console.log('Form is invalid');
      this.signupForm.markAllAsTouched();
      // tell why the form is invalid
      const errors = this.signupForm.errors;
      if (errors) {
        console.log('Form errors:', errors);
      } else {
        console.log('Email or password is invalid');
      }
    }
  }
}
