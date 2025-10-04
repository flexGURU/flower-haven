import { CommonModule } from '@angular/common';
import {
  Component,
  computed,
  effect,
  ElementRef,
  inject,
  signal,
  ViewChild,
} from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { OverlayPanelModule } from 'primeng/overlaypanel';
import { BadgeModule } from 'primeng/badge';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { CartService } from '../../components/cart/cart.service';
import { ProductService } from '../../../../shared/services/product.service';
import {
  categoryQuery,
  productQuery,
} from '../../../../shared/services/product.query';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { MessageModule } from 'primeng/message';
import { Category, Product } from '../../../../shared/models/models';
import { Popover, PopoverModule } from 'primeng/popover';
import { CartSignalService } from '../../components/cart/cart.signal.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.html',
  imports: [
    CommonModule,
    RouterLink,
    OverlayPanelModule,
    BadgeModule,
    FormsModule,
    ButtonModule,
    ProgressSpinnerModule,
    MessageModule,
    RouterLink,
    PopoverModule,
  ],
})
export class HeaderComponent {
  @ViewChild('op') op!: Popover;
  isMobileMenuOpen: boolean = false;
  cartTotal = 10;

  productsQueryData = productQuery();
  categoryQueryData = categoryQuery();

  #router = inject(Router);
  #cartService = inject(CartService);
  #productsService = inject(ProductService);
  #cartSignalService = inject(CartSignalService);

  searchQuery = signal('');

  onSearch() {
    if (this.searchQuery().trim()) {
      this.#router.navigate(['/products']);
    }
  }

  test(event: any) {
    if (this.op) {
      this.op.toggle(event);
    }

    this.#productsService.search.set(this.searchQuery());
  }

  cartItems = [];

  categories = computed<Category[]>(() => this.categoryQueryData.data() ?? []);
  products = computed<Product[]>(() => this.productsQueryData.data() ?? []);

  ngOnInit() {}
  cartItemCount = computed(() => this.#cartSignalService.cartCount());
  cartTotalAmount = computed(() => this.#cartSignalService.cartTotal());

  navigateToProduct(productId: string) {
    this.#router.navigate(['/product', productId]);
    this.searchQuery.set('');
  }


  viewCategory(categoryId: string) {
    this.#productsService.categoryId.set([categoryId]);
    this.#router.navigate(['/products']);
  }
}
