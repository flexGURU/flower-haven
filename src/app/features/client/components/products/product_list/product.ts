import { Component, computed, effect, inject, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { Category, Product } from '../../../../../shared/models/models';
import { ProductService } from '../../../../../shared/services/product.service';
import { CartService } from '../../cart/cart.service';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TabViewModule } from 'primeng/tabview';
import { SliderModule } from 'primeng/slider';
import { CheckboxModule } from 'primeng/checkbox';
import { SelectModule } from 'primeng/select';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { PaginatorModule } from 'primeng/paginator';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import {
  categoryQuery,
  productQuery,
} from '../../../../../shared/services/product.query';
import { ProgressSpinner } from 'primeng/progressspinner';
import { MessageModule } from 'primeng/message';
import { Meta, Title } from '@angular/platform-browser';

@Component({
  selector: 'app-product',
  templateUrl: './product.html',
  imports: [
    PaginatorModule,
    FormsModule,
    CommonModule,
    TabViewModule,
    SliderModule,
    CheckboxModule,
    SelectModule,
    CardModule,
    ButtonModule,
    RouterLink,
    ToastModule,
    ProgressSpinner,
    MessageModule,
  ],
  providers: [MessageService],
})
export class ProductComponent {
  allProducts: Product[] = [];
  products: Product[] = [];
  totalProducts = 0;
  currentPage = 1;
  pageSize = 12;
  pageTitle = 'Products';
  initialPriceTo = signal(50000);
  initialPriceFrom = signal(50);

  categoryQueryData = categoryQuery();

  first = signal(0);

  newProducts: Product[] = [];
  productQueryData = productQuery();

  priceRange = signal([this.initialPriceFrom(), this.initialPriceTo()]);
  selectedCategories = signal<string[]>([]);
  inStockOnly = false;
  sortBy = 'name';

  selectCategory = signal('');

  sortOptions = [
    { label: 'Name A-Z', value: 'name' },
    { label: 'Name Z-A', value: 'name_desc' },
    { label: 'Price Low to High', value: 'price' },
    { label: 'Price High to Low', value: 'price_desc' },
    { label: 'Newest First', value: 'created_desc' },
  ];

  private productService = inject(ProductService);
  private cartService = inject(CartService);
  private messageService = inject(MessageService);
  private title = inject(Title);
  private meta = inject(Meta);

  total = computed(() => this.productService.totalProducts());
  constructor() {}
  ngOnInit() {
    this.title.setTitle('Products - Floral Haven');
    this.meta.updateTag({
      name: 'description',
      content:
        'Explore our diverse range of products at Floral Haven. From fresh flowers to elegant arrangements, find the perfect choice for every occasion.',
    });
    this.productService.page.set(this.currentPage);
    this.productService.fetchProducts().subscribe((products) => {
      this.allProducts = products;
    });
  }

  categories = computed<Category[]>(() => this.categoryQueryData.data() ?? []);

  prices = computed(() => {
    return [this.priceRange()[0], this.priceRange()[1]];
  });

  applyFilters() {
    if (this.selectedCategories().length > 0) {
      this.productService.categoryId.update((previous) => [
        ...previous,
        ...this.selectedCategories(),
      ]);
    } else {
      this.productService.categoryId.set([]);
    }
    this.productService.priceFrom.set(this.prices()[0]);
    this.productService.priceTo.set(this.prices()[1]);
    this.productService.page.set(this.currentPage);
    this.productService.limit.set(this.pageSize);
  }
  clearFilters() {
    this.productService.page.set(1);
    this.productService.limit.set(10);
    this.productService.categoryId.set([]);
    this.productService.priceFrom.set(this.initialPriceFrom());
    this.productService.priceTo.set(this.initialPriceTo());
    this.productService.search.set('');
    this.initialPriceFrom.set(0);
    this.initialPriceTo.set(500000);
    this.priceRange.set([this.initialPriceFrom(), this.initialPriceTo()]);
  }

  refreshProducts() {
    this.productQueryData.refetch();
  }

  onPageChange(event: any) {
    this.productService.page.set(event.page + 1);
    this.productService.limit.set(event.rows);
    this.first.set(event.first);
  }

  addToCart(product: Product) {
    this.cartService.addToCart(product);
    this.messageService.add({
      severity: 'success',
      summary: 'Info',
      detail: `${product.name} added to cart`,
    });
  }
}
