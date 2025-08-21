import { Component, EventEmitter, inject, Input, Output } from '@angular/core';
import {
  FormGroup,
  FormBuilder,
  Validators,
  ReactiveFormsModule,
  FormsModule,
} from '@angular/forms';
import { MessageService } from 'primeng/api';
import { Category } from '../../../../shared/models/models';
import { ButtonModule } from 'primeng/button';
import { ToastModule } from 'primeng/toast';
import { CardModule } from 'primeng/card';
import { CommonModule } from '@angular/common';
import { InputTextModule } from 'primeng/inputtext';
import { TextareaModule } from 'primeng/textarea';
import { ProductService } from '../../../../shared/services/product.service';

@Component({
  selector: 'app-category-form',
  templateUrl: './category-form.html',
  imports: [
    ButtonModule,
    ToastModule,
    FormsModule,
    ReactiveFormsModule,
    CardModule,
    InputTextModule,
    TextareaModule,
  ],
  providers: [MessageService],
})
export class CategoryForm {
  @Input() categoryData: Category | null = null;
  @Input() isEditMode: boolean = false;
  @Output() onSave = new EventEmitter<Category>();
  @Output() onCancel = new EventEmitter<void>();
  statusMessage: string | null = null;

  categoryForm!: FormGroup;
  saving = false;

  private productService = inject(ProductService);

  constructor(
    private fb: FormBuilder,
    private messageService: MessageService,
  ) {}

  ngOnInit() {
    this.initializeForm();

    if (this.isEditMode && this.categoryData) {
      this.populateForm();
    }
  }

  initializeForm() {
    this.categoryForm = this.fb.group({
      name: [
        '',
        [
          Validators.required,
          Validators.minLength(2),
          Validators.maxLength(50),
        ],
      ],
      description: [
        '',
        [
          Validators.required,
          Validators.minLength(10),
          Validators.maxLength(200),
        ],
      ],
      imageUrl: ['', Validators.required],
    });
  }

  populateForm() {
    if (this.categoryData) {
      this.categoryForm.patchValue({
        name: this.categoryData.name,
        description: this.categoryData.description,
        imageUrl: this.categoryData.image_url?.[0] || '',
      });
    }
  }

  onSubmit() {
    if (!this.categoryForm.valid) {
      return;
    }

    const categoryData: Category = {
      name: this.categoryForm.value.name,
      description: this.categoryForm.value.description,
      image_url: [this.categoryForm.value.imageUrl],
      id: this.categoryData?.id, // Include id for update
    };

    this.saving = true; // Show a loading indicator
    let apiCall$;

    if (this.isEditMode) {
      apiCall$ = this.productService.updateCategory(categoryData);
    } else {
      apiCall$ = this.productService.addCategory(categoryData);
    }

    apiCall$.subscribe({
      next: (response) => {
        const message = this.isEditMode
          ? 'Category updated successfully!'
          : 'Category added successfully!';
        this.onSave.emit(response.data);
        this.categoryForm.reset();
      },
      error: (err) => {
        this.messageService.add({
          severity: 'error',
          summary: 'Operation failed',
          detail: err.message,
        });
        console.error('API Error:', err);
      },
      complete: () => {
        this.saving = false; // Hide loading indicator
      },
    });
  }

  onCancelModal() {
    this.categoryForm.reset();
    this.onCancel.emit();
  }
}
