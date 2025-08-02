import { Injectable, signal } from '@angular/core';
import { Category, Product } from '../models/models';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ProductService {
  private productsSubject = new BehaviorSubject<Product[]>([
    {
      id: '1',
      name: 'Red Roses Bouquet',
      description: 'Beautiful red roses perfect for any occasion',
      price: 45.99,
      image:
        'https://images.unsplash.com/photo-1587300003388-59208cc962cb?ixlib=rb-4.0.3&auto=format&fit=crop&w=600&q=60',
      categoryId: '1',
      stock: 0,
      createdAt: new Date(),
    },
    {
      id: '2',
      name: 'Mixed Wildflowers',
      description: 'Colorful wildflower arrangement',
      price: 32.99,
      image:
        'https://images.unsplash.com/photo-1490750967868-88aa4486c946?w=400',
      categoryId: '1',
      stock: 15,
      createdAt: new Date(),
    },
    {
      id: '3',
      name: 'Chocolate Gift Box',
      description: 'Premium chocolate assortment',
      price: 28.99,
      image: 'https://images.unsplash.com/photo-1549007908-17b1cebb6630?w=400',
      categoryId: '2',
      stock: 25,
      createdAt: new Date(),
    },
    {
      id: '4',
      name: 'Tulip Spring Bouquet',
      description: 'Fresh tulips in vibrant spring colors',
      price: 38.5,
      image:
        'https://images.unsplash.com/photo-1526397751294-331021109fbd?w=400',
      categoryId: '1',
      stock: 12,
      createdAt: new Date(),
    },
    {
      id: '5',
      name: 'Luxury Chocolate Strawberries',
      description: 'Chocolate-dipped strawberries with gold decoration',
      price: 49.99,
      image:
        'https://images.unsplash.com/photo-1481391319762-47dff72954d9?w=400',
      categoryId: '2',
      stock: 18,
      createdAt: new Date(),
    },
    {
      id: '6',
      name: 'Monstera Deliciosa Plant',
      description: 'Trendy tropical houseplant with split leaves',
      price: 34.99,
      image:
        'https://images.unsplash.com/photo-1606787366850-de6330128bfc?w=400',
      categoryId: '6',
      stock: 10,
      createdAt: new Date(),
    },
    {
      id: '7',
      name: 'Wedding Centerpiece',
      description: 'Elegant white floral arrangement for weddings',
      price: 89.99,
      image:
        'https://images.unsplash.com/photo-1519225421980-715cb0215aed?w=400',
      categoryId: '2',
      stock: 8,
      createdAt: new Date(),
    },
    {
      id: '8',
      name: 'Succulent Collection',
      description: 'Set of 4 low-maintenance succulents in ceramic pots',
      price: 42.5,
      image:
        'https://images.unsplash.com/photo-1519336056116-3e1f1a0a1e1b?w=400',
      categoryId: '6',
      stock: 14,
      createdAt: new Date(),
    },
  ]);

  products$ = this.productsSubject.asObservable();

  createProduct(product: Product) {
    const products = this.productsSubject.value;
    this.productsSubject.next([...products, product]);
  }

  updateProduct(id: string, product: Product) {
    const products = this.productsSubject.value;
    const index = products.findIndex((p) => p.id === product.id);
    if (index !== -1) {
      products[index] = product;
      this.productsSubject.next([...products]);
    }
  }

  getProductById(id: string) {
    return this.productsSubject.value[0];
  }

  deleteProduct(id: string) {
    const products = this.productsSubject.value.filter((p) => p.id !== id);
    this.productsSubject.next(products);
  }

  private categorySubject = new BehaviorSubject<Category[]>([
    {
      id: '1',
      name: 'Roses',
      description: 'Classic red and pink roses for special moments',
      image:
        'https://images.unsplash.com/photo-1518895949257-7621c3c786d7?w=600&h=400&fit=crop',
      productCount: 15,
    },
    {
      id: '2',
      name: 'Wedding Arrangements',
      description: 'Elegant floral designs for weddings',
      image:
        'https://images.unsplash.com/photo-1465495976277-4387d4b0e4a6?w=600&h=400&fit=crop',
      productCount: 23,
    },
    {
      id: '3',
      name: 'Graduation Gifts',
      description: 'Celebrate grads with floral and chocolate gifts',
      image:
        'https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=600&h=400&fit=crop',
      productCount: 18,
    },
    {
      id: '4',
      name: 'Birthday Bouquets',
      description: 'Cheerful birthday flowers and arrangements',
      image:
        'https://images.unsplash.com/photo-1563241527-3004b7be0ffd?w=600&h=400&fit=crop',
      productCount: 12,
    },
    {
      id: '5',
      name: 'Chocolate Hampers',
      description: 'Premium chocolate gift sets for every occasion',
      image:
        'https://images.unsplash.com/photo-1511381939415-e44015466834?w=600&h=400&fit=crop',
      productCount: 9,
    },
    {
      id: '6',
      name: 'Indoor Plants',
      description: 'Lively greens to brighten indoor spaces',
      image:
        'https://images.unsplash.com/photo-1586953208448-b95a79798f07?w=600&h=400&fit=crop',
      productCount: 9,
    },
    {
      id: '7',
      name: 'Outdoor Plants',
      description: 'Beautiful and resilient outdoor plants',
      image:
        'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=600&h=400&fit=crop',
      productCount: 9,
    },
  ]);

  getCategories() {
    return this.categorySubject.asObservable();
  }
  addCategory(category: Category) {
    const categories = this.categorySubject.value;
    this.categorySubject.next([...categories, category]);
  }

  updateCategory(category: Category) {
    const categories = this.categorySubject.value;
    const index = categories.findIndex((c) => c.id === category.id);
    if (index !== -1) {
      categories[index] = category;
      this.categorySubject.next([...categories]);
    }
  }

  getProductsByCategory(categoryId: string) {
    return this.productsSubject.value.filter(
      (product) => product.categoryId === categoryId,
    );
  }
}
