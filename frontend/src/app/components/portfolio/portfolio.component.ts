import { Component } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { Router } from '@angular/router';
import { Asset } from '../../models/asset.model';
import { ApiService } from '../../services/api.service';
import { CommonModule } from '@angular/common';
import { AddAssetComponent } from '../add-asset/add-asset.component';
import { Modal } from 'bootstrap';

@Component({
  selector: 'app-portfolio',
  imports: [CommonModule, AddAssetComponent],
  templateUrl: './portfolio.component.html',
  styleUrl: './portfolio.component.css'
})
export class PortfolioComponent {
  message: string = 'Welcome to your portfolio!';
  assets: Asset[] = [];
  modalInstance?: Modal;

  constructor(private apiService: ApiService, private router: Router, public auth: AuthService) {}

  ngOnInit(): void {
    if (!this.auth.isAuthenticated()) {
      alert('You are not logged in. Please log in to access this page.');
      this.router.navigate(['/login']);
      return;
    }

    this.updateAssets();
  }

  ngAfterViewInit() {
    const modalElement = document.getElementById('addAssetModal');
    if (modalElement) {
      this.modalInstance = new Modal(modalElement);
    }
  }

  updateAssets(): void {
    const options: Object = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    };
    this.apiService.getRequest<{assets: Asset[]}>("asset/list", options).subscribe({
      next: (data) => {
        this.assets = data.assets;
        if (this.assets.length === 0) {
          this.message = 'Your portfolio is empty.';
        }
      },
      error: (err) => {
        alert('Failed to load portfolio. Please try again later.');
        console.error(err);
      }
    });
  }

  onAssetCreated(asset: Asset) {
    const options: Object = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    };
    this.apiService.postRequest<{message: string, asset: Asset}>("asset/add", asset, options).subscribe({
      next: (data) => {
        alert('Asset added successfully!');
        this.assets.push(data.asset);
        this.modalInstance?.hide();
      },
      error: (err) => {
        alert('Failed to add asset.' + (err.error?.message || 'Please try again later.'));
      }
    });
  }
}
