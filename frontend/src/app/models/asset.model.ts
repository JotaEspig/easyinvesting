export interface Asset {
  id: number;
  name: string;
  code: string;
  asset_type: number; // 0: stock
  currency: number; // 0: BRL, 1: USD
  user_id: number;
  hold_avg_price: number;
  hold_quantity: number;
}
