import { Component, inject } from '@angular/core';
import { CarouselModule } from 'primeng/carousel';
import { Category, Product } from '../../../../../shared/models/models';
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ProductService } from '../../../../../shared/services/product.service';
import { CategoryComponent } from '../category/category';
import { HeroComponent } from '../hero/hero';
import { FeaturedComponent } from '../featured/featured';
import { TestimonialsComponent } from '../testimonials/testimonials';

@Component({
  selector: 'app-home',
  templateUrl: './home.html',
  imports: [
    CarouselModule,
    CardModule,
    ButtonModule,
    RouterLink,
    CommonModule,
    CategoryComponent,
    HeroComponent,
    FeaturedComponent,
    TestimonialsComponent,
  ],
})
export class HomeComponent {
  categories: Category[] = [];
  featuredProducts: Product[] = [];

  productService = inject(ProductService);
  constructor() {}

  ngOnInit() {}

  loadFeaturedProducts() {
    this.productService.products$.subscribe((products) => {
      this.featuredProducts = products
        .filter((product) => product.stock > 0)
        .slice(0, 6); // Show first 6 available products
    });
  }
}
