import { Component, computed, effect, inject } from '@angular/core';
import { Category } from '../../../../../shared/models/models';
import { ProductService } from '../../../../../shared/services/product.service';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { categoryQuery } from '../../../../../shared/services/product.query';

@Component({
  selector: 'app-category',
  templateUrl: './category.html',
  imports: [CommonModule],
})
export class CategoryComponent {
  #productService = inject(ProductService);
  #router = inject(Router);
  #categoryQueryData = categoryQuery();

  categories = computed<Category[]>(() => this.#categoryQueryData.data() ?? []);

  constructor() {
    effect(() => {});
  }

  loadProductsByCategory(categoryId: string) {
    this.#productService.categoryId.set([categoryId]);

    return this.#router.navigate(['/products']);
  }
}
