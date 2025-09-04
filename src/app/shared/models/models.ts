export interface Product {
  id?: string;
  name: string;
  description: string;
  price: number;
  image_url: string[];
  category_id: string;
  stock_quantity: number;
  createdAt?: Date;
  category_data?: Category;
  stems?: stems[];
}

interface stems {
  label: string;
  value: string;
  price: number;
}
export interface Category {
  id?: string;
  name: string;
  description: string;
  image_url?: string[];
  productCount?: number;
}

export interface User {
  name: string;
  email: string;
  phone_number: string;
  password: string;
  is_admin: string;
}

export const initialProd: Product = {
  name: '',
  description: '',
  price: 10,
  image_url: [],
  category_id: '',
  stock_quantity: 0,
};
