import { Injectable, signal, computed, inject } from '@angular/core';
import { Product } from '../../../../shared/models/models';
import { CartItem, Cart } from './cart.model';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class CartService {
  private cartSubject = new BehaviorSubject<Cart>({
    items: [],
    total: 0,
    itemCount: 0,
  });

  public cart$ = this.cartSubject.asObservable();

  constructor() {
    this.loadCartFromStorage();
  }

  addToCart(product: Product, quantity: number = 1): void {
    const currentCart = this.cartSubject.value;
    const existingItemIndex = currentCart.items.findIndex(
      (item) => item.product.id === product.id,
    );

    if (existingItemIndex > -1) {
      currentCart.items[existingItemIndex].quantity += quantity;
    } else {
      currentCart.items.push({ product, quantity, amount: product.price * quantity });
    }

    this.updateCart(currentCart);
  }


  
  removeFromCart(productId: string): void {
    const currentCart = this.cartSubject.value;

    const updatedCart: Cart = {
      ...currentCart,
      items: currentCart.items.filter((item) => item.product.id !== productId),
    };

    this.updateCart(updatedCart);
  }

  updateQuantity(productId: string, quantity: number): void {
    const currentCart = this.cartSubject.value;
    const itemIndex = currentCart.items.findIndex(
      (item) => item.product.id === productId,
    );

    if (itemIndex > -1) {
      if (quantity <= 0) {
        this.removeFromCart(productId);
      } else {
        currentCart.items[itemIndex].quantity = quantity;
        this.updateCart(currentCart);
      }
    }
  }

  clearCart(): void {
    const emptyCart: Cart = {
      items: [],
      total: 0,
      itemCount: 0,
    };
    this.updateCart(emptyCart);
  }

  private updateCart(cart: Cart): void {
    cart.total = cart.items.reduce(
      (sum, item) => sum + item.product.price * item.quantity,
      0,
    );

    cart.total = cart.total;

    this.cartSubject.next(cart);
    this.saveCartToStorage(cart);
  }

  private loadCartFromStorage(): void {
    const savedCart = localStorage.getItem('flower-cart');
    if (savedCart) {
      const cart = JSON.parse(savedCart);
      this.cartSubject.next(cart);
    }
  }

  private saveCartToStorage(cart: Cart): void {
    localStorage.setItem('flower-cart', JSON.stringify(cart));
  }

  getCartItemCount(): number {
    return this.cartSubject.value.items.reduce(
      (count, item) => count + item.quantity,
      0,
    );
  }
}
