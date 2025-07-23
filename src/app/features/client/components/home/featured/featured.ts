import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { CardModule } from 'primeng/card';
import { Product } from '../../../../../shared/models/models';
import { ButtonModule } from 'primeng/button';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-featured',
  templateUrl: './featured.html',
  imports: [CardModule, CommonModule, RouterLink, ButtonModule],
})
export class FeaturedComponent {
  featuredProducts: Product[] = [];

  addToCart(product: Product) {
    // Add to cart logic
    console.log('Adding to cart:', product);
  }

  test(){
    console.log("llooe");
    
  }
}
