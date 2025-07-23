import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { OverlayPanelModule } from 'primeng/overlaypanel';
import { BadgeModule } from 'primeng/badge';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { CartService } from '../../components/cart/cart.service';
import { ProductService } from '../../../../shared/services/product.service';

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
  ],
})
export class HeaderComponent {
  isMobileMenuOpen: boolean = false;
  cartItemCount: number = 9;
  cartTotal = 10;

  router = inject(Router);
  cartService = inject(CartService);
  productsService = inject(ProductService);

  searchQuery = '';

  onSearch() {
    if (this.searchQuery.trim()) {
      this.router.navigate(['/search'], {
        queryParams: { q: this.searchQuery },
      });
    }
  }

  cartItems = [];

  login() {
    // Logic for login action
  }

  categories: any[] = [];

  ngOnInit() {
    this.loadCategories();
  }

  loadCategories() {
    this.productsService.getCategories().subscribe((response) => {
      this.categories = response;
    });
  }
}
