import { Injectable, signal } from '@angular/core';
import { CartItem } from './cart.model';
import { Product } from '../../../../shared/models/models';

@Injectable({
  providedIn: 'root',
})
export class CartService {
  // cartItem = signal<CartItem | []>([]);

  // private readonly cartKey = 'shopping_cart';
  // constructor() {}

  // addToCart(product: Product, quantity: number = 1): void {}

  // private loadCartFromStorage(): CartItem[] | [] {
  //   return JSON.parse(localStorage.getItem(this.cartKey) || '[]');
  // }

  // private saveCartToStorage(cart: CartItem[]): void {
  //   localStorage.setItem(this.cartKey, JSON.stringify(cart));
  // }
  // private clearCartStorage(): void {
  //   localStorage.removeItem(this.cartKey);
  // }
  // clearCart(): void {
  //   this.cartItem.set(null);
  //   this.clearCartStorage();
  // }
}
