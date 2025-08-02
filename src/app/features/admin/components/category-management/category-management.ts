import { Component } from '@angular/core';
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
    PaginatorModule, InputTextModule,ButtonModule
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

  // Mock product counts (replace with actual service calls)
  productCounts: { [categoryId: string]: number } = {};

  constructor(
    private confirmationService: ConfirmationService,
    private messageService: MessageService,
    // Inject your category service here
  ) {}

  ngOnInit() {
    this.loadCategories();
  }

  loadCategories() {
    this.loading = true;
    // Replace with actual service call
    setTimeout(() => {
      this.categories = [
        {
          id: '1',
          name: 'Electronics',
          description:
            'Electronic devices, gadgets, and accessories including smartphones, laptops, and home appliances',
        },
        {
          id: '2',
          name: 'Clothing',
          description:
            'Fashion items, apparel, and accessories for men, women, and children',
        },
        {
          id: '3',
          name: 'Books',
          description:
            'Physical and digital books, magazines, and other reading materials',
        },
      ];

      // Mock product counts
      this.productCounts = {
        '1': 25,
        '2': 18,
        '3': 12,
      };

      this.filteredCategories = [...this.categories];
      this.loading = false;
    }, 1000);
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
    const productCount = this.getProductCount(category.id);

    if (productCount > 0) {
      this.messageService.add({
        severity: 'error',
        summary: 'Cannot Delete',
        detail: `Cannot delete category "${category.name}" as it contains ${productCount} products. Please remove or reassign products first.`,
      });
      return;
    }

    this.confirmationService.confirm({
      message: `Are you sure you want to delete the category "${category.name}"? This action cannot be undone.`,
      header: 'Delete Category',
      icon: 'pi pi-exclamation-triangle',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.deleteCategory(category);
      },
    });
  }

  deleteCategory(category: Category) {
    // Replace with actual service call
    this.categories = this.categories.filter((c) => c.id !== category.id);
    delete this.productCounts[category.id];
    this.applyFilters();

    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: `Category "${category.name}" deleted successfully`,
    });
  }

  onCategorySave(categoryData: Category) {
    if (this.isEditMode) {
      // Update existing category
      const index = this.categories.findIndex((c) => c.id === categoryData.id);
      if (index !== -1) {
        this.categories[index] = categoryData;
        this.messageService.add({
          severity: 'success',
          summary: 'Success',
          detail: `Category "${categoryData.name}" updated successfully`,
        });
      }
    } else {
      // Add new category
      const newCategory: Category = {
        ...categoryData,
        id: this.generateId(), // Replace with actual ID generation
      };
      this.categories.push(newCategory);
      this.productCounts[newCategory.id] = 0;

      this.messageService.add({
        severity: 'success',
        summary: 'Success',
        detail: `Category "${categoryData.name}" created successfully`,
      });
    }

    this.categoryForm = false;
    this.applyFilters();
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
}
