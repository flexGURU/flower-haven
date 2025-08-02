import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormGroup, FormBuilder, Validators, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { MessageService } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { Category } from '../../../../shared/models/models';
import { ProductService } from '../../../../shared/services/product.service';
import { CheckboxModule } from 'primeng/checkbox';
import { CardModule } from 'primeng/card';
import { InputNumber, InputNumberModule } from 'primeng/inputnumber';
import { InputTextModule } from 'primeng/inputtext';
import { FileUploadModule } from 'primeng/fileupload';
import { ChipModule } from 'primeng/chip';
import { Select, SelectModule } from "primeng/select";
import { ToastModule } from 'primeng/toast';

@Component({
  selector: 'app-product-form',
  templateUrl: './product-form.html',
  imports: [
    ButtonModule,
    CheckboxModule,
    CardModule,
    InputNumberModule,
    InputTextModule,
    FileUploadModule,
    CommonModule,
    ChipModule,
    SelectModule,
    ReactiveFormsModule,
    ToastModule
],
})
export class ProductForm {
  productForm: FormGroup;
  categories: Category[] = [];
  productImages: string[] = [];
  isEditMode = false;
  productId: string | null = null;
  saving = false;

  constructor(
    private fb: FormBuilder,
    private productService: ProductService,
    private route: ActivatedRoute,
    private router: Router,
    private messageService: MessageService,
  ) {
    this.productForm = this.fb.group({
      name: ['', Validators.required],
      description: ['', Validators.required],
      price: [0, [Validators.required, Validators.min(0.01)]],
      categoryId: ['', Validators.required],
      stockQuantity: [0, [Validators.required, Validators.min(0)]],
      inStock: [true],
      featured: [false],
      active: [true],
      tags: [[]],
    });
  }

  ngOnInit() {
    this.loadCategories();

    this.route.params.subscribe((params) => {
      if (params['id']) {
        this.isEditMode = true;
        this.productId = params['id'];
        this.loadProduct(this.productId ?? '');
      }
    });
  }

  loadCategories() {
    this.productService.getCategories().subscribe((categories) => {
      this.categories = categories;
    });
  }

  loadProduct(id: string) {}

  onImageUpload(event: any) {
    // Handle image upload
    for (let file of event.files) {
      const reader = new FileReader();
      reader.onload = (e: any) => {
        this.productImages.push(e.target.result);
      };
      reader.readAsDataURL(file);
    }
  }

  onImageRemove(event: any) {
    // Handle image removal
  }

  removeImage(index: number) {
    this.productImages.splice(index, 1);
  }

  onSubmit() {
    if (this.productForm.valid) {
      this.saving = true;
      const formData = {
        ...this.productForm.value,
        images: this.productImages,
      };

      const operation = this.isEditMode
        ? this.productService.updateProduct(this.productId!, formData)
        : this.productService.createProduct(formData);

    //   operation.subscribe({
    //     next: (product) => {
    //       this.messageService.add({
    //         severity: 'success',
    //         summary: 'Success',
    //         detail: `Product ${this.isEditMode ? 'updated' : 'created'} successfully`,
    //       });
    //       this.router.navigate(['/admin/products']);
    //     },
    //     error: (error) => {
    //       console.error('Error saving product:', error);
    //       this.messageService.add({
    //         severity: 'error',
    //         summary: 'Error',
    //         detail: `Failed to ${this.isEditMode ? 'update' : 'create'} product`,
    //       });
    //       this.saving = false;
    //     },
    //   });
    }
  }

  saveAsDraft() {
    // Save as draft logic
    this.onSubmit();
  }
}
