import { Injectable, signal, computed } from '@angular/core';
import { Product } from '../../../../shared/models/models';
import { CartItem, Cart } from './cart.model';


@Injectable({
  providedIn: 'root',
})
export class CartService {
  private cartItems = signal<CartItem[]>([]);

  cart = computed<Cart>(() => {
    const items = this.cartItems();
    const total = items.reduce(
      (sum, item) => sum + item.product.price * item.quantity,
      0
    );
    const itemCount = items.reduce((sum, item) => sum + item.quantity, 0);

    return {
      items,
      total,
      itemCount,
    };
  });

  addToCart(
    product: Product,
    quantity: number = 1,
    selectedDate?: Date,
    personalMessage?: string
  ) {
    const existingItem = this.cartItems().find(
      (item) => item.product.id === product.id
    );

    if (existingItem) {
      this.updateQuantity(product.id, existingItem.quantity + quantity);
    } else {
      this.cartItems.update((items) => [
        ...items,
        {
          product,
          quantity,
          selectedDate,
          personalMessage,
        },
      ]);
    }
  }

  removeFromCart(productId: string) {
    this.cartItems.update((items) =>
      items.filter((item) => item.product.id !== productId)
    );
  }

  updateQuantity(productId: string, quantity: number) {
    if (quantity <= 0) {
      this.removeFromCart(productId);
      return;
    }

    this.cartItems.update((items) =>
      items.map((item) =>
        item.product.id === productId ? { ...item, quantity } : item
      )
    );
  }

  clearCart() {
    this.cartItems.set([]);
  }
}
