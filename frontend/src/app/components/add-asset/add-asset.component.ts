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
    name: '',
    code: '',
    currency: 0,
  };

  isAssetValid(): boolean {
    return this.asset.name.trim() !== '' &&
           this.asset.code.trim() !== '' &&
           this.asset.currency >= 0 && this.asset.currency <= 1;
  }

  onSubmit() {
    this.assetCreated.emit({ ...this.asset });
    this.asset = { name: '', code: '', currency: 0 };
  }
}
