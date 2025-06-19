import { Routes } from '@angular/router';
import { HomeComponent } from './components/home/home.component';
import { LoginComponent } from './components/login/login.component';
import { SignupComponent } from './components/signup/signup.component';
import { HelloauthComponent } from './components/helloauth/helloauth.component';
import { PortfolioComponent } from './components/portfolio/portfolio.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'signup', component: SignupComponent },
  { path: 'helloauth', component: HelloauthComponent },
  { path: 'portfolio', component: PortfolioComponent },
  { path: '**', redirectTo: '/' }
];
