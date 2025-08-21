import { Component, inject } from '@angular/core';
import { Category } from '../../../../shared/models/models';
import { ConfirmationService, MessageService } from 'primeng/api';
import { DialogModule } from 'primeng/dialog';
import { CategoryForm } from '../category-form/category-form';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { ToastModule } from 'primeng/toast';
import { BadgeModule } from 'primeng/badge';
import { TableModule } from 'primeng/table';
import { FormsModule } from '@angular/forms';
import { PaginatorModule } from 'primeng/paginator';
import { InputTextModule } from 'primeng/inputtext';
import { ButtonModule } from 'primeng/button';
import { ProductService } from '../../../../shared/services/product.service';

@Component({
  selector: 'app-category-management',
  templateUrl: './category-management.html',
  imports: [
    DialogModule,
    CategoryForm,
    ConfirmDialog,
    ToastModule,
    BadgeModule,
    TableModule,
    FormsModule,
    PaginatorModule,
    InputTextModule,
    ButtonModule,
  ],
  providers: [MessageService, ConfirmationService],
})
export class CategoryManagement {
  categories: Category[] = [];
  filteredCategories: Category[] = [];
  loading = false;

  // Form dialog properties
  categoryForm = false;
  showCategoryDetails = false;
  isEditMode = false;
  selectedCategory!: Category;

  // Filter properties
  searchTerm = '';

  private productService = inject(ProductService);
  // Mock product counts (replace with actual service calls)
  productCounts: { [categoryId: string]: number } = {};

  constructor(
    private confirmationService: ConfirmationService,
    private messageService: MessageService,
  ) {}

  ngOnInit() {
    this.loadCategories();
  }

  loadCategories() {
    console.log();

    this.loading = true;
    this.productService.getCategories().subscribe({
      next: (categories) => {
        this.categories = categories;
        this.filteredCategories = [...this.categories];
        this.loading = false;
      },
      error: (error) => {
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to load categories. Please try again later.',
        });
        this.loading = false;
      },
    });
  }

  applyFilters() {
    this.filteredCategories = this.categories.filter(
      (category) =>
        category.name.toLowerCase().includes(this.searchTerm.toLowerCase()) ||
        category.description
          .toLowerCase()
          .includes(this.searchTerm.toLowerCase()),
    );
  }

  clearFilters() {
    this.searchTerm = '';
    this.filteredCategories = [...this.categories];
  }

  showCategoryForm() {
    this.isEditMode = false;
    // this.selectedCategory = {};
    this.categoryForm = true;
  }

  editCategory(category: Category) {
    this.isEditMode = true;
    this.selectedCategory = { ...category };
    this.categoryForm = true;
  }

  viewCategory(category: Category) {
    this.selectedCategory = category;
    this.showCategoryDetails = true;
  }

  editCategoryFromDetails() {
    this.showCategoryDetails = false;
    this.editCategory(this.selectedCategory!);
  }

  confirmDelete(category: Category) {
    this.confirmationService.confirm({
      message: `Are you sure you want to delete the category "${category.name}"?`,
      header: 'Confirm Deletion',
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        this.productService.deleteCategory(category.id!).subscribe({
          next: () => {
            this.messageService.add({
              severity: 'success',
              summary: 'Success',
              detail: `Category "${category.name}" has been deleted.`,
            });
          },
          error: (err) => {
            console.error('Deletion failed:', err);
            this.messageService.add({
              severity: 'error',
              summary: 'Error',
              detail: 'Failed to delete category.',
            });
          },
        });
      },
    });
  }

  deleteCategory(category: Category) {
    console.log('Deleting category:', category);
  }

  onCategorySave(categoryData: Category) {
    this.categoryForm = false; 
    this.loadCategories();
    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: `Category "${categoryData.name}" has been ${
        this.isEditMode ? 'updated' : 'added'
      }.`,
    });
  }

  onFormCancel() {
    this.categoryForm = false;
  }

  onFormClose() {
    // this.selectedCategory = null;
    this.isEditMode = false;
  }

  getProductCount(categoryId: string): number {
    return this.productCounts[categoryId] || 0;
  }

  getActiveProductCount(categoryId: string): number {
    // Replace with actual service call to get active product count
    return Math.floor(this.getProductCount(categoryId) * 0.8);
  }

  private generateId(): string {
    return Date.now().toString();
  }
  getStockSeverity(
    stock: number,
  ): 'info' | 'success' | 'warn' | 'danger' | 'secondary' | 'contrast' {
    if (stock === 0) return 'danger';
    if (stock <= 10) return 'warn';
    return 'success';
  }
}
