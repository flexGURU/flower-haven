export interface Product {
  id?: string;
  name: string;
  description: string;
  price: number;
  image: string;
  categoryId: string;
  stock: number;
  createdAt?: Date;
}

export interface Category {
  id?: string;
  name: string;
  description: string;
  image_url?: string[];
  productCount?: number;
}
