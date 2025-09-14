import { Component, inject } from '@angular/core';
import { Category } from '../../../../../shared/models/models';
import { ProductService } from '../../../../../shared/services/product.service';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-category',
  templateUrl: './category.html',
  imports: [CommonModule],
})
export class CategoryComponent {
  categories: Category[] = [];
  #productService = inject(ProductService);
  #router = inject(Router);

  ngOnInit() {
    this.loadCategories();
  }
  loadCategories() {
    this.#productService.getCategories().subscribe((categories) => {
      this.categories = categories;
    });
  }

  loadProductsByCategory(categoryId: string) {
    this.#productService.categoryId.update((previous) => [...previous, categoryId]);

    return this.#router.navigate(['/products']);
  }
}
