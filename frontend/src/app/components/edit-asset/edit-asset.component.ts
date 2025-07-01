import { Component, Output, EventEmitter } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-edit-asset',
  imports: [FormsModule],
  templateUrl: './edit-asset.component.html',
  styleUrl: './edit-asset.component.css'
})
export class EditAssetComponent {
  @Output() assetEntryCreated = new EventEmitter<any>();

  assetEntry = {
    type: -1,
    price: 0,
    quantity: 0,
    date: new Date().toISOString().split('T')[0]
  };

  isAssetEntryValid(): boolean {
    return this.assetEntry.price > 0 && this.assetEntry.quantity > 0;
  }

  onSubmit() {
    if (typeof this.assetEntry.price === 'string') {
      this.assetEntry.price = parseFloat(this.assetEntry.price);
    }
    if (typeof this.assetEntry.quantity === 'string') {
      this.assetEntry.quantity = parseInt(this.assetEntry.quantity, 10);
    }
    if (typeof this.assetEntry.type === 'string') {
      this.assetEntry.type = parseInt(this.assetEntry.type, 10);
    }

    this.assetEntry.date = new Date(this.assetEntry.date).toISOString();
    this.assetEntryCreated.emit({ ...this.assetEntry });
    this.assetEntry = {
      type: -1,
      price: 0,
      quantity: 0,
      date: new Date().toISOString().split('T')[0]
    };
  }
}
