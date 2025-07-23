import { Component, inject } from '@angular/core';
import { Category } from '../../../../../shared/models/models';
import { ProductService } from '../../../../../shared/services/product.service';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-category',
  templateUrl: './category.html',
  imports: [RouterLink, CommonModule],
})
export class CategoryComponent {
  categories: Category[] = [];
  productService = inject(ProductService);

  ngOnInit() {
    this.loadCategories();
  }
  loadCategories() {
    this.productService.getCategories().subscribe((categories) => {
      console.log('categories', categories);

      this.categories = categories; // Show first 4 categories
    });
  }
}
