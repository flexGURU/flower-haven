import { Component, Input } from '@angular/core';
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

  ngOnInit() {
    this.cartService.cart$.subscribe((cart) => {
      this.cart = cart;
      this.total = cart.total
    });
  }

  updateQuantity(productId: string, quantity: number) {
    this.cartService.updateQuantity(productId, quantity);
  }

  removeItem(productId: string) {
    console.log(productId);

    this.cartService.removeFromCart(productId);
  }

  clearCart() {
    this.cartService.clearCart();
  }

  getTotalItems(): number {
    return this.cart.items.reduce((total, item) => total + item.quantity, 0);
  }

  checkOut() {
    this.checkout = true

  }
}
