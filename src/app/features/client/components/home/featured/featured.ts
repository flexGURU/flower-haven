import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { CardModule } from 'primeng/card';
import { Product } from '../../../../../shared/models/models';
import { ButtonModule } from 'primeng/button';
import { Router, RouterLink } from '@angular/router';
import { ProductService } from '../../../../../shared/services/product.service';

@Component({
  selector: 'app-featured',
  templateUrl: './featured.html',
  imports: [CardModule, CommonModule, ButtonModule],
})
export class FeaturedComponent {
  featuredProducts: Product[] = [];
  #productService = inject(ProductService);
  #router = inject(Router);

  addToCart(product: Product) {
    // Add to cart logic
    console.log('Adding to cart:', product);
  }

  viewAllProducts() {
    this.#productService.page.set(1);
    this.#productService.limit.set(10);
    this.#productService.categoryId.set([]);
    this.#productService.priceFrom.set(0);
    this.#productService.priceTo.set(0);
    this.#productService.search.set('');
    this.#router.navigate(['/products']);
  }
}
