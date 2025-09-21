import { Component, computed, effect, inject, signal } from '@angular/core';
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
import {
  categoryQuery,
  productQuery,
} from '../../../../shared/services/product.query';
import { MessageModule } from 'primeng/message';
import { ProgressSpinnerModule } from 'primeng/progressspinner';

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
    MessageModule,
    ProgressSpinnerModule,
    PaginatorModule,
    MessageModule,
  ],
  providers: [ConfirmationService, MessageService],
})
export class ProductManagement {
  products: Product[] = [];
  filteredProducts: Product[] = [];
  productForm = signal(false);
  productQueryData = productQuery();
  categoryQueryData = categoryQuery();
  first = signal(0);
  loading = signal(false);

  // Filters
  searchTerm = signal('');
  selectedCategory = signal('');
  stockFilter = signal('');

  categoryOptions: any[] = [];
  stockOptions = [
    { label: 'In Stock', value: 'inStock' },
    { label: 'Out of Stock', value: 'outOfStock' },
    { label: 'Low Stock', value: 'lowStock' },
  ];

  #productService = inject(ProductService);

  total = computed(() => this.#productService.totalProducts());
  constructor(
    private confirmationService: ConfirmationService,
    private messageService: MessageService,
  ) {
    effect(() => {
      const products = this.productQueryData.data() ?? [];
      this.products = products;
      this.filteredProducts = [...products];

      this.categoryOptions = [
        { name: 'All Categories', id: '' },
        ...(this.categoryQueryData.data() ?? []),
      ];
    });
  }

  private productService = inject(ProductService);

  onPageChange(event: any) {
    this.#productService.page.set(event.page + 1);
    this.#productService.limit.set(event.rows);
    this.first.set(event.first);
  }

  searchProducts() {
    this.#productService.search.set(this.searchTerm());
  }

  searchByCategory() {
    if (this.selectedCategory()) {
      this.#productService.categoryId.set([this.selectedCategory()]);
    } else {
      this.#productService.categoryId.set([]);
    }
  }

  searchByStock() {
    const stock = this.stockFilter();

    this.filteredProducts =
      this.productQueryData.data()?.filter((product) => {
        if (stock === 'inStock') return product.stock_quantity > 10;
        if (stock === 'lowStock')
          return product.stock_quantity > 0 && product.stock_quantity <= 10;
        if (stock === 'outOfStock') return product.stock_quantity === 0;
        return true;
      }) ?? [];
  }

  applyFilters() {
    this.searchProducts();
    this.searchByCategory();
    this.searchByStock();
  }
  clearFilters() {
    this.searchTerm.set('');
    this.selectedCategory.set('');
    this.stockFilter.set('');
    this.filteredProducts = [...this.products];
    this.#productService.search.set('');
    this.#productService.categoryId.set([]);
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
      message: `Are you sure you want to delete the product "${product.name}"?`,
      header: 'Confirm Deletion',
      icon: 'pi pi-exclamation-triangle',
      rejectButtonProps: {
        label: 'Cancel',
        severity: 'secondary',
        outlined: true,
      },
      acceptButtonProps: {
        label: 'Confirm',
        loading: this.loading(),
      },
      accept: () => {
        this.loading.set(true);
        this.productService.deleteProduct(product.id!).subscribe({
          next: () => {
            this.loading.set(false);
            this.messageService.add({
              severity: 'success',
              summary: 'Success',
              detail: `Category "${product.name}" has been deleted.`,
            });
            this.productQueryData.refetch();
          },
          error: (err) => {
            console.error('Deletion failed:', err);
            this.messageService.add({
              severity: 'error',
              summary: 'Error',
              detail: 'Failed to delete product.',
            });
          },
        });
      },
    });
  }

  showProductDetails = false;
  isEditMode = false;
  selectedProduct: Product | null = null;

  showProductForm() {
    this.isEditMode = false;
    this.selectedProduct = null;
    this.productForm.set(true);
  }

  editProduct(product: Product) {
    this.isEditMode = true;
    this.selectedProduct = { ...product };
    this.productForm.set(true);
  }

  viewProduct(product: any) {
    this.selectedProduct = product;
    this.showProductDetails = true;
  }

  editProductFromDetails() {
    this.showProductDetails = false;
    this.editProduct(this.selectedProduct!);
  }

  onProductSave(productData: Product) {
    // Handle save logic
    this.productForm.set(false);
    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: `Product "${productData.name}" has been ${
        this.isEditMode ? 'updated' : 'added'
      }.`,
    });
    this.productQueryData.refetch();
  }

  onFormCancel() {
    this.productForm.set(false);
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

  getStockSeverity(
    stock: number,
  ): 'info' | 'success' | 'warn' | 'danger' | 'secondary' | 'contrast' {
    if (stock === 0) return 'danger';
    if (stock <= 10) return 'warn';
    return 'success';
  }

  deleteProduct(productId: string) {}
}
