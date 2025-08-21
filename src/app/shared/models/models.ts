export interface Product {
  id?: string;
  name: string;
  description: string;
  price: number;
  image_url: string[];
  category_id: string;
  stock_quantity: number;
  createdAt?: Date;
}

export interface Category {
  id?: string;
  name: string;
  description: string;
  image_url?: string[];
  productCount?: number;
}
