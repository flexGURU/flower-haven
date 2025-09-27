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
import { TestimonialsComponent } from '../testimonials/testimonials';
import { FeaturedComponent } from '../featured/featured';

@Component({
  selector: 'app-home',
  templateUrl: './home.html',
  imports: [
    CarouselModule,
    CardModule,
    ButtonModule,
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
}
