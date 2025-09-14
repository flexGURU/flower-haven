import { CommonModule } from '@angular/common';
import { Component, effect, inject, signal } from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { OverlayPanelModule } from 'primeng/overlaypanel';
import { BadgeModule } from 'primeng/badge';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { CartService } from '../../components/cart/cart.service';
import { ProductService } from '../../../../shared/services/product.service';
import { productQuery } from '../../../../shared/services/product.query';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { MessageModule } from 'primeng/message';

@Component({
  selector: 'app-header',
  templateUrl: './header.html',
  imports: [
    CommonModule,
    RouterLink,
    OverlayPanelModule,
    BadgeModule,
    FormsModule,
    ButtonModule,
    ProgressSpinnerModule,
    MessageModule,
    RouterLink,
  ],
})
export class HeaderComponent {
  isMobileMenuOpen: boolean = false;
  cartItemCount = 0;
  cartTotal = 10;

  products = productQuery();

  #router = inject(Router);
  #cartService = inject(CartService);
  #productsService = inject(ProductService);

  searchQuery = signal('');

  onSearch() {
    if (this.searchQuery().trim()) {
      this.#router.navigate(['/products']);
    }
  }

  test() {
    this.#productsService.search.set(this.searchQuery());
  }

  cartItems = [];

  categories: any[] = [];

  ngOnInit() {
    this.loadCategories();
    this.#cartService.cart$.subscribe((cart) => {
      this.cartItemCount = cart.items.reduce(
        (count, item) => count + item.quantity,
        0,
      );
    });
  }

  loadCategories() {
    this.#productsService.getCategories().subscribe((response) => {
      this.categories = response;
    });
  }

  cartCount() {
    this.cartItemCount = this.#cartService.getCartItemCount();
  }
}
