import { Component, EventEmitter, Input, Output } from '@angular/core';
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

@Component({
  selector: 'app-category-form',
  templateUrl: './category-form.html',
  imports: [
    ButtonModule,
    ToastModule,
    FormsModule,
    ReactiveFormsModule,
    CardModule, InputTextModule, TextareaModule,
  ],
  providers: [MessageService],
})
export class CategoryForm {
  @Input() categoryData: Category | null = null;
  @Input() isEditMode: boolean = false;
  @Output() onSave = new EventEmitter<Category>();
  @Output() onCancel = new EventEmitter<void>();

  categoryForm!: FormGroup;
  saving = false;

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
    });
  }

  populateForm() {
    if (this.categoryData) {
      this.categoryForm.patchValue({
        name: this.categoryData.name,
        description: this.categoryData.description,
      });
    }
  }

  async onSubmit() {
    if (this.categoryForm.valid) {
      this.saving = true;

      try {
        // Simulate API call delay
        await new Promise((resolve) => setTimeout(resolve, 1000));

        const categoryData: Category = {
          id: this.isEditMode ? this.categoryData!.id : '',
          name: this.categoryForm.value.name.trim(),
          description: this.categoryForm.value.description.trim(),
        };

        this.onSave.emit(categoryData);

        // Reset form after successful save
        if (!this.isEditMode) {
          this.categoryForm.reset();
        }
      } catch (error) {
        console.error('Error saving category:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to save category. Please try again.',
        });
      } finally {
        this.saving = false;
      }
    } else {
      // Mark all fields as touched to show validation errors
      Object.keys(this.categoryForm.controls).forEach((key) => {
        this.categoryForm.get(key)?.markAsTouched();
      });

      this.messageService.add({
        severity: 'warn',
        summary: 'Validation Error',
        detail: 'Please fix the form errors before submitting.',
      });
    }
  }

  onCancelModal() {
    this.categoryForm.reset();
    this.onCancel.emit();
  }

  // Utility method to check if field has specific error
  hasError(fieldName: string, errorType: string): boolean {
    const field = this.categoryForm.get(fieldName);
    return !!(field?.errors?.[errorType] && field?.touched);
  }

  // Get field error message
  getFieldError(fieldName: string): string {
    const field = this.categoryForm.get(fieldName);
    if (field?.errors && field?.touched) {
      if (field.errors['required']) return `${fieldName} is required`;
      if (field.errors['minlength']) return `${fieldName} is too short`;
      if (field.errors['maxlength']) return `${fieldName} is too long`;
    }
    return '';
  }
}
