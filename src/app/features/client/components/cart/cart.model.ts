import { Product } from "../../../../shared/models/models";

export interface CartItem {
  product: Product;
  quantity: number;
  amount: number;
  selectedDate?: Date;
  personalMessage?: string;
}

export interface Cart {
  items: CartItem[];
  total: number;
  itemCount: number;
}