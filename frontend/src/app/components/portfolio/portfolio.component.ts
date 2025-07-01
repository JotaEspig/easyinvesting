import { Component } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { Router } from '@angular/router';
import { Asset } from '../../models/asset.model';
import { AssetEntry } from '../../models/entry.model';
import { ApiService } from '../../services/api.service';
import { CommonModule } from '@angular/common';
import { AddAssetComponent } from '../add-asset/add-asset.component';
import { EditAssetComponent } from '../edit-asset/edit-asset.component';
import { Modal } from 'bootstrap';
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-portfolio',
  standalone: true,
  imports: [CommonModule, AddAssetComponent, EditAssetComponent],
  templateUrl: './portfolio.component.html',
  styleUrl: './portfolio.component.css'
})
export class PortfolioComponent {
  message: string = 'Welcome to your portfolio!';
  assets: Asset[] = [];
  modalAddAssetInstance?: Modal;
  modalEditAssetInstance?: Modal;
  currentlySelectedAsset?: Asset;

  constructor(
    private apiService: ApiService,
    private router: Router,
    public auth: AuthService,
    private messageService: MessageService
  ) {}

  ngOnInit(): void {
    if (!this.auth.isAuthenticated()) {
      this.messageService.add({
        severity: 'warn',
        summary: 'Not Logged In',
        detail: 'You are not logged in. Please log in to access this page.'
      });
      this.router.navigate(['/login']);
      return;
    }

    this.updateAssets();
  }

  ngAfterViewInit() {
    const modalElement = document.getElementById('addAssetModal');
    if (modalElement) {
      this.modalAddAssetInstance = new Modal(modalElement);
    }

    const editModalElement = document.getElementById('editAssetModal');
    if (editModalElement) {
      this.modalEditAssetInstance = new Modal(editModalElement);
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
        this.messageService.add({
          severity: 'success',
          summary: 'Update Successful',
          detail: data.message
        });
        this.updateAssets();
      },
      error: (err) => {
        this.messageService.add({
          severity: 'error',
          summary: 'Update Failed',
          detail: 'Failed to update prices. Please try again later.'
        });
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
              asset.profitability = ((asset.market_price - asset.hold_avg_price) / asset.hold_avg_price) * 100;
              asset.profitability = Math.round(asset.profitability * 100) / 100;
              console.log(`Updated market price for asset ${asset.id}: ${asset.market_price}`);
              console.log(`Profitability for asset ${asset.id}: ${asset.profitability}%`);
            },
            error: (err) => {
              console.error(`Failed to fetch market price for asset ${asset.id}:`, err);
            }
          });
        }
      },
      error: (err) => {
        this.messageService.add({
          severity: 'error',
          summary: 'Load Failed',
          detail: 'Failed to load portfolio. Please try again later.'
        });
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
        this.messageService.add({
          severity: 'success',
          summary: 'Asset Added',
          detail: 'Asset added successfully!'
        });
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
        this.modalAddAssetInstance?.hide();
      },
      error: (err) => {
        this.messageService.add({
          severity: 'error',
          summary: 'Add Failed',
          detail: 'Failed to add asset.' + (err.error?.message || 'Please try again later.')
        });
      }
    });
  }

  getProfitabilityStyle(asset: Asset): Object {
    if (asset.market_price > asset.hold_avg_price) {
      return {'color': 'green'};
    } else if (asset.market_price < asset.hold_avg_price) {
      return {'color': 'red'};
    }
    return {'color': 'gray'};
  }

  onAssetUpdated(entry: AssetEntry): void {
    // Should not happen, but just in case
    if (!this.currentlySelectedAsset) {
      this.messageService.add({
        severity: 'warn',
        summary: 'No Asset Selected',
        detail: 'Strange...'
      });
      return;
    }

    entry.asset_id = this.currentlySelectedAsset.id;

    const options: Object = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    };
    this.apiService.postRequest<{message: string, entry: AssetEntry}>(`asset/entry/add`, entry, options).subscribe({
      next: (data) => {
        this.messageService.add({
          severity: 'success',
          summary: 'Asset Updated',
          detail: 'Asset updated successfully!'
        });
        this.unselectAsset();
        this.modalEditAssetInstance?.hide();
        this.updateAssets();
      },
      error: (err) => {
        this.messageService.add({
          severity: 'error',
          summary: 'Update Failed',
          detail: 'Failed to update asset.' + (err.error?.message || 'Please try again later.')
        });
      }
    });
  }

  selectAsset(asset: Asset): void {
    this.currentlySelectedAsset = asset;
  }

  unselectAsset(): void {
    this.currentlySelectedAsset = undefined;
  }
}
