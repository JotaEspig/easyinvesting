import { Asset } from './asset.model';

export interface AssetEntry {
  id: number;
  price: number;
  quantity: number;
  type: number;
  date: string;
  asset_id: number;
  asset: Asset;
}
