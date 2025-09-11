import { Component, effect, inject } from '@angular/core';
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
import { productQuery } from '../../../../shared/services/product.query';
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
  ],
  providers: [ConfirmationService, MessageService],
})
export class ProductManagement {
  products: Product[] = [];
  filteredProducts: Product[] = [];
  categories: Category[] = [];
  productForm = false;
  productQueryData = productQuery();

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
    private confirmationService: ConfirmationService,
    private messageService: MessageService,
  ) {

    effect(()=> {
      const products = this.productQueryData.data() ?? [];
      this.products = products;
      this.filteredProducts = [...products];
      this.applyFilters();
    })
  }

  private productService = inject(ProductService);

  ngOnInit() {
    this.loadCategories();
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

  applyFilters() {
    this.filteredProducts = (this.productQueryData.data() ?? []).filter(
      (product) => {
        const matchesSearch =
          product.name.toLowerCase().includes(this.searchTerm.toLowerCase()) ||
          product.description
            .toLowerCase()
            .includes(this.searchTerm.toLowerCase());

        const matchesCategory =
          !this.selectedCategory ||
          product.category_id === this.selectedCategory;

        let matchesStock = true;
        if (this.stockFilter === 'inStock') {
          matchesStock = product.stock_quantity > 10;
        } else if (this.stockFilter === 'lowStock') {
          matchesStock =
            product.stock_quantity > 0 && product.stock_quantity <= 10;
        } else if (this.stockFilter === 'outOfStock') {
          matchesStock = product.stock_quantity === 0;
        }

        return matchesSearch && matchesCategory && matchesStock;
      },
    );
  }
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
      message: `Are you sure you want to delete the product "${product.name}"?`,
      header: 'Confirm Deletion',
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        this.productService.deleteProduct(product.id!).subscribe({
          next: () => {
            this.messageService.add({
              severity: 'success',
              summary: 'Success',
              detail: `Category "${product.name}" has been deleted.`,
            });
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
    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: `Product "${productData.name}" has been ${
        this.isEditMode ? 'updated' : 'added'
      }.`,
    });
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

  getStockSeverity(
    stock: number,
  ): 'info' | 'success' | 'warn' | 'danger' | 'secondary' | 'contrast' {
    if (stock === 0) return 'danger';
    if (stock <= 10) return 'warn';
    return 'success';
  }

  deleteProduct(productId: string) {}
}
