import { Routes } from '@angular/router';
import { HomeComponent } from './components/home/home.component';
import { LoginComponent } from './components/login/login.component';
import { SignupComponent } from './components/signup/signup.component';
import { HelloComponent } from './components/hello/hello.component';
import { HelloauthComponent } from './components/helloauth/helloauth.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'signup', component: SignupComponent },
  { path: 'hello', component: HelloComponent },
  { path: 'helloauth', component: HelloauthComponent },
];
