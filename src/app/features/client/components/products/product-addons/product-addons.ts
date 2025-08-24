import { Component, EventEmitter, inject, OnInit, Output } from '@angular/core';
import { Product } from '../../../../../shared/models/models';
import { ProductService } from '../../../../../shared/services/product.service';
import { DynamicDialogRef, DynamicDialogConfig } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { CommonModule } from '@angular/common';
import { CartService } from '../../cart/cart.service';
import { InputNumber } from 'primeng/inputnumber';
import { CardModule } from 'primeng/card';
import { CheckboxModule } from 'primeng/checkbox';
import { FormsModule } from '@angular/forms';

export interface AddsOnInterface extends Product {
  selected: boolean;
}

@Component({
  selector: 'app-product-addons',
  templateUrl: './product-addons.html',
  imports: [
    ButtonModule,
    CommonModule,
    CardModule,
    CheckboxModule,
    FormsModule,
  ],
  standalone: true,
})
export class ProductAddonsComponent {
  selected?: boolean;
  addons: AddsOnInterface[] = [];
  quantities: { [key: string]: number } = {};
  @Output() addonAddedToCart = new EventEmitter<Product>();
  @Output() addonRemovedFromCart = new EventEmitter<Product>();

  constructor(private productService: ProductService) {}
  private cartService = inject(CartService);

  ngOnInit() {
    this.productService.getProducts().subscribe((products) => {
      this.addons = products.map((product) => ({
        ...product,
        selected: false,
      }));
    });
  }

  onAddonToggle(addon: AddsOnInterface) {
    if (addon.selected) {
      this.addonAddedToCart.emit(addon);
    } else {
      // Remove from cart
      if (addon) {
        this.addonRemovedFromCart.emit(addon);
      }
    }
  }

  addAddonToCart(addon: Product) {
    this.cartService.addToCart(addon);
  }

  getSelectedAddons() {
    return this.addons.filter((addon) => addon.selected);
  }

  getTotalAddonsPrice(): number {
    return this.getSelectedAddons().reduce(
      (total, addon) => total + addon.price,
      0,
    );
  }
}
