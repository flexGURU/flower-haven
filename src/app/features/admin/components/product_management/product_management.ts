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
import { ProductFormComponent } from '../product-form/product-form';
import { DialogModule } from 'primeng/dialog';

@Component({
  selector: 'app-product-management',
  templateUrl: './product_management.html',
  imports: [
    SelectModule,
    ConfirmDialog,
    FormsModule,
    TableModule,
    BadgeModule,
DialogModule,
    ToastModule,
    ConfirmDialog,
    CommonModule,
    ButtonModule,
    PaginatorModule,
    InputTextModule,
    ProductFormComponent,
  ],
  providers: [ConfirmationService, MessageService],
})
export class ProductManagement {
  products: Product[] = [];
  filteredProducts: Product[] = [];
  categories: Category[] = [];
  loading = false;
  productForm = false;

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

  

  showProductDetails = false;
  isEditMode = false;
  selectedProduct: any = null;

  showProductForm() {
    this.isEditMode = false;
    this.selectedProduct = null;
    this.productForm = true;
  }

  editProduct(product: any) {
    this.isEditMode = true;
    this.selectedProduct = { ...product };
    this.productForm = true;
  }

  viewProduct(product: any) {
    this.selectedProduct = product;
    this.showProductDetails = true;
  }

  editProductFromDetails() {
    this.showProductDetails = false;
    this.editProduct(this.selectedProduct);
  }

  onProductSave(productData: any) {
    // Handle save logic
    this.productForm = false;
    // Refresh your products list
  }

  onFormCancel() {
    this.productForm = false;
  }

  onFormClose() {
    this.selectedProduct = null;
    this.isEditMode = false;
  }

  getStockLabel(stock: number): string {
    if (stock === 0) return 'Out of Stock';
    if (stock <= 10) return 'Low Stock';
    return 'In Stock';
  }

  getStockSeverity(stock: number): "info" | "success" | "warn" | "danger" | "secondary" | "contrast" {
    if (stock === 0) return 'danger';
    if (stock <= 10) return 'warn';
    return 'success';
  }

  deleteProduct(productId: string) {}
}
