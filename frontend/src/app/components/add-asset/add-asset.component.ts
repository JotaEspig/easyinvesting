import { Component, Output, EventEmitter } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-add-asset',
  imports: [FormsModule],
  templateUrl: './add-asset.component.html',
  styleUrls: ['./add-asset.component.css']

})
export class AddAssetComponent {
  @Output() assetCreated = new EventEmitter<any>();

  asset = {
    code: '',
    currency: 0,
  };

  isAssetValid(): boolean {
    return this.asset.code.trim() !== '' &&
           this.asset.currency >= 0 && this.asset.currency <= 1;
  }

  onSubmit() {
    if (typeof this.asset.currency === 'string') {
      this.asset.currency = parseInt(this.asset.currency, 10);
    }
    this.assetCreated.emit({ ...this.asset });
    this.asset = { code: '', currency: 0 };
  }
}
