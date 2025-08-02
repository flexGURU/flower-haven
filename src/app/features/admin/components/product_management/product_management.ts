import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { BadgeModule } from 'primeng/badge';
import { DropdownModule } from 'primeng/dropdown';
import { TableModule } from 'primeng/table';
import { ToastModule } from 'primeng/toast';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { ConfirmationService, MessageService } from 'primeng/api';
import { Product, Category } from '../../../../shared/models/models';
import { ProductService } from '../../../../shared/services/product.service';
import { CommonModule } from '@angular/common';
import { SelectModule } from 'primeng/select';
import { ButtonModule } from 'primeng/button';
import { PaginatorModule } from 'primeng/paginator';
import { InputTextModule } from 'primeng/inputtext';

@Component({
  selector: 'app-product-management',
  templateUrl: './product_management.html',
  imports: [
    SelectModule,
    ConfirmDialog,
    FormsModule,
    TableModule,
    BadgeModule,
    RouterLink,
    ToastModule,
    ConfirmDialog,
    CommonModule,
    ButtonModule,
    PaginatorModule,
    InputTextModule,
  ],
  providers: [ConfirmationService, MessageService],
})
export class ProductManagement {
  products: Product[] = [];
  filteredProducts: Product[] = [];
  categories: Category[] = [];
  loading = false;

  productStatus = 'add';

  // Filters
  searchTerm = '';
  selectedCategory = '';
  stockFilter = '';

  categoryOptions: any[] = [];
  stockOptions = [
    { label: 'In Stock', value: 'inStock' },
    { label: 'Out of Stock', value: 'outOfStock' },
    { label: 'Low Stock', value: 'lowStock' },
  ];

  constructor(
    private productService: ProductService,
    private confirmationService: ConfirmationService,
    private messageService: MessageService,
  ) {}

  ngOnInit() {
    this.loadProducts();
    this.loadCategories();
  }

  loadProducts() {
    this.loading = true;
    this.productService.products$.subscribe({
      next: (products) => {
        this.products = products;
        this.filteredProducts = [...products];
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading products:', error);
        this.loading = false;
      },
    });
  }

  loadCategories() {
    this.productService.getCategories().subscribe((categories) => {
      this.categories = categories;
      this.categoryOptions = [
        { name: 'All Categories', id: '' },
        ...categories,
      ];
    });
  }

  applyFilters() {}

  clearFilters() {
    this.searchTerm = '';
    this.selectedCategory = '';
    this.stockFilter = '';
    this.filteredProducts = [...this.products];
  }

  getStockClass(
    quantity: number,
  ): 'info' | 'success' | 'warn' | 'danger' | 'secondary' | 'contrast' {
    if (quantity === 0) return 'danger';
    if (quantity < 10) return 'warn';
    return 'success';
  }

  confirmDelete(product: Product) {
    this.confirmationService.confirm({
      message: `Are you sure you want to delete "${product.name}"?`,
      header: 'Confirm Delete',
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        this.deleteProduct(product.id);
      },
    });
  }

  addProduct(status: string) {
    console.log('status', status);
  }

  editProduct(product: Product, status: string) {
    console.log('status', status);
  }

  deleteProduct(productId: string) {}
}
