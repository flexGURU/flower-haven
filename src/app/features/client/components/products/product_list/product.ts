import { Component, inject } from '@angular/core';
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
  ],
  providers: [MessageService],
})
export class ProductComponent {
  allProducts: Product[] = []; // Store all products fetched
  products: Product[] = []; // Products displayed on the current page
  categories: Category[] = [];
  totalProducts = 0;
  currentPage = 1;
  pageSize = 12; // Initial page size
  pageTitle = 'All Products';

  // Filters
  priceRange = [0, 500];
  selectedCategories: string[] = [];
  inStockOnly = false;
  sortBy = 'name';

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
  constructor(private route: ActivatedRoute) {}

  ngOnInit() {
    this.loadCategories();
    this.productService.products$.subscribe((products) => {
      this.allProducts = products; // Get all products once
      this.applyFilters(); // Apply filters and pagination initially
    });

    // Listen for route changes
    this.route.queryParams.subscribe((params) => {
      if (params['category']) {
        this.selectedCategories = [params['category']];
      }
      this.applyFilters();
    });
  }

  loadCategories() {
    this.productService.getCategories().subscribe((categories) => {
      this.categories = categories;
    });
  }

  applyFilters() {
    // 1. Filter products based on current filter settings
    let filteredProducts = [...this.allProducts]; // Start with all products

    // Filter by Price Range
    filteredProducts = filteredProducts.filter(
      (p) => p.price >= this.priceRange[0] && p.price <= this.priceRange[1],
    );

    // Filter by Categories
    if (this.selectedCategories.length > 0) {
      filteredProducts = filteredProducts.filter(
        (p) => this.selectedCategories.includes(p.categoryId), // Assuming product has categoryId
      );
    }

    // Filter by In Stock Only
    if (this.inStockOnly) {
      filteredProducts = filteredProducts.filter((p) => p.stock > 0);
    }

    // 2. Sort the filtered products
    // filteredProducts.sort((a, b) => {
    //   switch (this.sortBy) {
    //     case 'name':
    //       return a.name.localeCompare(b.name);
    //     case 'name_desc':
    //       return b.name.localeCompare(a.name);
    //     case 'price':
    //       return a.price - b.price;
    //     case 'price_desc':
    //       return b.price - a.price;
    //     case 'created_desc':
    //       // Assuming product has a 'createdAt' property for newest first
    //       // You might need to adjust this based on your Product model
    //       return (
    //         new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
    //       );
    //     default:
    //       return 0;
    //   }
    // });

    this.totalProducts = filteredProducts.length; // Update totalProducts for paginator

    // 3. Apply pagination to the filtered and sorted products
    const start = (this.currentPage - 1) * this.pageSize;
    const end = start + this.pageSize;
    this.products = filteredProducts.slice(start, end);

    // Ensure currentPage is valid if filters reduce total products significantly
    if (
      this.currentPage > Math.ceil(this.totalProducts / this.pageSize) &&
      this.totalProducts > 0
    ) {
      this.currentPage = Math.ceil(this.totalProducts / this.pageSize);
      this.applyFilters(); // Re-apply to fetch correct page
    } else if (this.totalProducts === 0) {
      this.currentPage = 1;
    }
  }

  clearFilters() {
    this.priceRange = [0, 500];
    this.selectedCategories = [];
    this.inStockOnly = false;
    this.sortBy = 'name';
    this.currentPage = 1; // Reset current page when clearing filters
    this.applyFilters();
  }

  onPageChange(event: any) {
    this.currentPage = event.page + 1;
    this.pageSize = event.rows; // Update pageSize if rowsPerPageOptions is used
    this.applyFilters(); // Re-apply filters which will also paginate
  }

  addToCart(product: Product) {
    console.log("pp");
    
    this.cartService.addToCart(product);
    this.messageService.add({
      severity: 'success',
      summary: 'Info',
      detail: `${product.name} added to cart`,
      life: 3000,
    });
  }
}
