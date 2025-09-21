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
  is_message_card?: boolean;
  is_add_on?: boolean;
  has_stems?: boolean;
  stems?: Stem[];
}

export interface Stem {
  stem_count: string;
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
