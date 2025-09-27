import {
  Component,
  computed,
  effect,
  EventEmitter,
  inject,
  OnInit,
  Output,
  signal,
} from '@angular/core';
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
import {
  addonQuery,
  productQuery,
} from '../../../../../shared/services/product.query';
import { PaginatorModule } from 'primeng/paginator';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { ScrollerModule } from 'primeng/scroller';

interface Addon extends Product {
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
    PaginatorModule,
    ProgressSpinnerModule,
    ScrollerModule,
  ],
  standalone: true,
})
export class ProductAddonsComponent {
  selected?: boolean;
  quantities: { [key: string]: number } = {};
  addOnsQuery = addonQuery();
  @Output() addonAddedToCart = new EventEmitter<Product>();
  @Output() addonRemovedFromCart = new EventEmitter<Product>();
  selectedAddon = signal<Addon>({} as Addon);
  first = signal(0);
  addons = computed<Addon[]>(
    () =>
      (this.addOnsQuery.data() ?? []).map((addon) => ({
        ...addon,
        selected: false,
      })) as Addon[],
  );

  #productService = inject(ProductService);
  total = computed(() => this.#productService.totalAddOns());

  constructor() {
    effect(() => {});
  }

  onAddonToggle(addon: Addon) {
    if (addon.selected) {
      this.addonAddedToCart.emit(addon);
    } else {
      if (addon) {
        this.addonRemovedFromCart.emit(addon);
      }
    }
  }

  onPageChange(event: any) {
    this.#productService.page.set(event.page + 1);
    this.#productService.limit.set(event.rows);
    this.first.set(event.first);
  }
}
