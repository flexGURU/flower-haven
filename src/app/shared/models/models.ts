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
  id?: number;
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

export interface OrderItem {
  product_id: string;
  quantity: number;
  amount: number;
  payment_method?: 'subscription' | 'one_time';
  frequency?: 'weekly' | 'monthly';
}

export interface OrderPayload {
  user_name: string;
  user_phone_number: string;
  payment_status: boolean;
  status: string;
  delivery_date: string;
  time_slot: string;
  total: number;
  reference: string;
  items: OrderItem[];
}
