import { CommonModule } from '@angular/common';
import {
  Component,
  EventEmitter,
  inject,
  Input,
  Output,
  signal,
} from '@angular/core';
import {
  FormGroup,
  FormBuilder,
  Validators,
  ReactiveFormsModule,
} from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { Category, Product } from '../../../../shared/models/models';
import { CheckboxModule } from 'primeng/checkbox';
import { CardModule } from 'primeng/card';
import { InputNumberModule } from 'primeng/inputnumber';
import { InputTextModule } from 'primeng/inputtext';
import { FileUploadModule } from 'primeng/fileupload';
import { ChipModule } from 'primeng/chip';
import { SelectModule } from 'primeng/select';
import { ToastModule } from 'primeng/toast';
import { BadgeModule } from 'primeng/badge';
import { TextareaModule } from 'primeng/textarea';
import { ProductService } from '../../../../shared/services/product.service';
import { FirebaseService } from '../../services/firebase';
import { MessageService } from 'primeng/api';

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
    ToastModule,
    BadgeModule,
    TextareaModule,
  ],
  providers: [MessageService],
})
export class ProductFormComponent {
  @Input() productData: Product | null = null;
  @Input() isEditMode: boolean = false;
  @Output() onSave = new EventEmitter<Product>();
  @Output() onCancelForm = new EventEmitter<void>();

  productForm!: FormGroup;
  saving = false;
  categories: Category[] = []; // Load your categories here
  currentImageUrl: string[] | null = null;
  selectedFiles = signal<File[]>([]);
  selectedImagePreview = signal<string[]>([]);

  private productService = inject(ProductService);
  private firebaseService = inject(FirebaseService);
  private messageService = inject(MessageService);

  constructor(private fb: FormBuilder) {}

  ngOnInit() {
    this.initializeForm();
    this.loadCategories();

    if (this.isEditMode && this.productData) {
      this.populateForm();
      console.log('sss', this.productData);
    }
  }

  initializeForm() {
    this.productForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      description: ['', [Validators.required, Validators.minLength(10)]],
      price: [0, [Validators.required, Validators.min(0.01)]],
      categoryId: ['', Validators.required],
      stock: [0, [Validators.required, Validators.min(0)]],
    });
  }

  populateForm() {
    if (this.productData) {
      this.productForm.patchValue({
        name: this.productData.name,
        description: this.productData.description,
        price: this.productData.price,
        categoryId: this.productData.category_id,
        stock: this.productData.stock_quantity,
      });
      this.currentImageUrl = this.productData.image_url;
    }
  }

  loadCategories() {
    this.productService.getCategories().subscribe((reponse) => {
      this.categories = reponse;
    });
  }

  onImageSelect(event: any) {
    if (event.files && event.files.length) {
      const newFiles: File[] = [...event.files]; // convert FileList to array

      this.selectedFiles.update((files) => [...files, ...newFiles]);

      for (let file of newFiles) {
        const reader = new FileReader();
        reader.onload = () => {
          this.selectedImagePreview.update((previews) => [
            ...previews,
            reader.result as string,
          ]);
        };
        reader.readAsDataURL(file);
      }
    }
  }

  removeCurrentImage() {
    this.currentImageUrl = null;
  }

  getStockStatus(): string {
    const stockValue = this.productForm.get('stock')?.value || 0;
    if (stockValue === 0) return 'Out of Stock';
    if (stockValue <= 10) return 'Low Stock';
    return 'In Stock';
  }

  getStockSeverity():
    | 'info'
    | 'success'
    | 'warn'
    | 'danger'
    | 'secondary'
    | 'contrast' {
    const stockValue = this.productForm.get('stock')?.value || 0;
    if (stockValue === 0) return 'danger';
    if (stockValue <= 10) return 'warn';
    return 'success';
  }

  async onSubmit() {
    if (this.productForm.valid) {
      this.saving = true;

      const productData: Product = {
        name: this.productForm.value.name,
        description: this.productForm.value.description,
        price: this.productForm.value.price,
        image_url: this.isEditMode ? this.productData!.image_url : [],
        category_id: this.productForm.value.categoryId,
        stock_quantity: this.productForm.value.stock,
        createdAt: this.isEditMode ? this.productData!.createdAt : new Date(),
      };

      try {
        if (this.selectedFiles().length > 0) {
          const imageUrl = await this.firebaseService.uploadImage(
            this.selectedFiles(),
          );
          productData.image_url = [imageUrl];
        }

        let apiCall$;

        if (this.isEditMode) {
          if (this.productData?.id) {
            apiCall$ = this.productService.updateProduct(
              productData,
              this.productData.id,
            );
          }
        } else {
          apiCall$ = this.productService.addProduct(productData);
        }

        apiCall$?.subscribe({
          next: (response) => {
            const message = this.isEditMode
              ? 'Product updated successfully!'
              : 'Product added successfully!';
            this.onSave.emit(response.data);
            this.productForm.reset();
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
            this.saving = false;
          },
        });
      } catch (error) {
        console.error('Error uploading image:', error);
      } finally {
        this.saving = false;
      }
    }
  }

  async uploadImage(file: File): Promise<string> {
    // Implement your image upload logic here
    // Return the uploaded image URL
    return '';
  }

  onCancel() {
    this.onCancelForm.emit();
  }
}
