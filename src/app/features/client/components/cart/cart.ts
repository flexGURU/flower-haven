import { Component, computed, inject, Input } from '@angular/core';
import { CartService } from './cart.service';
import { Cart } from './cart.model';
import { CardModule } from 'primeng/card';
import { RouterLink } from '@angular/router';
import { InputNumber, InputNumberModule } from 'primeng/inputnumber';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { DividerModule } from 'primeng/divider';
import { ButtonModule } from 'primeng/button';
import { Dialog, DialogModule } from 'primeng/dialog';
import { CheckoutComponent } from '../checkout/checkout.component';
import { CartSignalService } from './cart.signal.service';
import { Product } from '../../../../shared/models/models';

@Component({
  selector: 'app-cart',
  templateUrl: './cart.html',
  imports: [
    CardModule,
    DividerModule,
    RouterLink,
    InputNumberModule,
    FormsModule,
    CommonModule,
    ButtonModule,
    DialogModule,
    CheckoutComponent,
  ],
})
export class CartComponent {
  cart: Cart = {
    items: [],
    total: 0,
    itemCount: 0,
  };
  promoCode = '';
  checkout = false;
  total: number = 0;

  constructor(private cartService: CartService) {}

  #cartSignalService = inject(CartSignalService);
  cartItems = computed(() => this.#cartSignalService.cart());
  cartTotal = computed(() => this.#cartSignalService.cartTotal());
  cartCount = computed(() => this.#cartSignalService.cartCount());

  ngOnInit() {
    this.cartService.cart$.subscribe((cart) => {
      this.cart = cart;
      this.total = cart.total;
    });
  }

  updateQuantity(productId: Product, quantity: number) {
    this.#cartSignalService.updateCart(productId, quantity);
  }

  removeItem(productId: string) {
    this.#cartSignalService.removeFromCart(productId);
  }

  clearCart() {
    this.#cartSignalService.clearCart();
  }

  checkOut() {
    console.log('cart details', this.cart);

    this.checkout = true;
  }
}
