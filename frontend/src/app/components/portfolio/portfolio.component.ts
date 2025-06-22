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

  updatePrices(): void {
    const options: Object = {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
      }
    };
    this.apiService.postRequest<{message: string}>("realtimeupdate", {}, options).subscribe({
      next: (data) => {
        alert(data.message);
        this.updateAssets();
      },
      error: (err) => {
        alert('Failed to update prices. Please try again later.');
        console.error(err);
      }
    });
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

        for (let asset of this.assets) {
          this.apiService.getRequest<{market_price: number}>(`asset/${asset.id}/realtime`, options).subscribe({
            next: (data) => {
              asset.market_price = data.market_price;
              console.log(`Updated market price for asset ${asset.id}: ${asset.market_price}`);
            },
            error: (err) => {
              console.error(`Failed to fetch market price for asset ${asset.id}:`, err);
            }
          });
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
        this.apiService.getRequest<{market_price: number}>(`asset/${data.asset.id}/realtime`, options).subscribe({
          next: (data2) => {
            data.asset.market_price = data2.market_price;
            console.log(`Updated market price for asset ${data.asset.id}: ${data.asset.market_price}`);
          },
          error: (err) => {
            console.error(`Failed to fetch market price for asset ${data.asset.id}:`, err);
          }
        });
        this.assets.push(data.asset);
        this.modalInstance?.hide();
      },
      error: (err) => {
        alert('Failed to add asset.' + (err.error?.message || 'Please try again later.'));
      }
    });
  }
}
