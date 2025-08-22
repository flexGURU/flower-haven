import { Component, inject, OnInit } from '@angular/core';
import { Product } from '../../../../../shared/models/models';
import { ProductService } from '../../../../../shared/services/product.service';
import { DynamicDialogRef, DynamicDialogConfig } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { CommonModule } from '@angular/common';
import { CartService } from '../../cart/cart.service';
import { InputNumber } from 'primeng/inputnumber';

@Component({
  selector: 'app-product-addons',
  templateUrl: './product-addons.html',
  imports: [ButtonModule, CommonModule, InputNumber],
  standalone: true,
})
export class ProductAddonsComponent {
  addons: Product[] = [];
  quantities: { [key: string]: number } = {};

  constructor(private productService: ProductService) {}
  private cartService = inject(CartService);

  ngOnInit() {
    this.productService.getProducts().subscribe((relatedProducts) => {
      this.addons = relatedProducts;
    });
  }

  addAddonToCart(addon: Product) {
    this.cartService.addToCart(addon);
  }
}
