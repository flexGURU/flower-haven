import { computed, effect, Injectable, signal } from '@angular/core';
import { Cart, CartItem } from './cart.model';
import { Product } from '../../../../shared/models/models';

@Injectable({
  providedIn: 'root',
})
export class CartSignalService {
  cart = signal<CartItem[]>([]);

  constructor() {
    this.loadFromCartStorage();
    effect(() => {});
  }

  addToCart(product: Product, quantity: number = 1): void {
    const existingItem = this.cart().find(
      (item) => item.product.id === product.id,
    );
    if (existingItem) {
      this.updateCart(product, existingItem.quantity + quantity);
    } else {
      this.cart.update((currentCart) => [
        ...currentCart,
        { product, quantity, amount: product.price * quantity },
      ]);
      this.saveToCartStorage();
    }
  }

  removeFromCart(productId: string): void {
    this.cart.update((currentCart) =>
      currentCart.filter((item) => item.product.id !== productId),
    );
    this.saveToCartStorage();
  }
  clearCart(): void {
    localStorage.removeItem('floral-cart');
    this.cart.set([]);
  }

  updateCart(product: Product, quantity: number): void {
    this.cart.update((currentCart) =>
      currentCart.map((item) =>
        item.product.id === product.id
          ? { ...item, quantity, amount: item.product.price * quantity }
          : item,
      ),
    );
    this.saveToCartStorage();
  }

  private loadFromCartStorage(): void {
    const cartData = localStorage.getItem('floral-cart');
    if (cartData) {
      this.cart.set(JSON.parse(cartData));
    }
  }

  private saveToCartStorage(): void {
    localStorage.setItem('floral-cart', JSON.stringify(this.cart()));
  }

  cartCount = computed(() =>
    this.cart().reduce((acc, item) => acc + item.quantity, 0),
  );
  cartTotal = computed(() =>
    this.cart().reduce(
      (acc, item) => acc + item.product.price * item.quantity,
      0,
    ),
  );

  perCartItemTotal = computed(() =>
    this.cart().map((item) => item.product.price * item.quantity),
  );
}
