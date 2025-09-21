import { CommonModule } from '@angular/common';
import {
  Component,
  effect,
  EventEmitter,
  inject,
  Input,
  Output,
  signal,
  ViewChild,
  viewChild,
} from '@angular/core';
import {
  FormGroup,
  FormBuilder,
  Validators,
  ReactiveFormsModule,
  FormArray,
  FormsModule,
} from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { Category, Product, Stem } from '../../../../shared/models/models';
import { Checkbox, CheckboxModule } from 'primeng/checkbox';
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
import { TableModule } from 'primeng/table';
import { categoryQuery } from '../../../../shared/services/product.query';

@Component({
  selector: 'app-product-form',
  templateUrl: './product-form.html',
  imports: [
    ButtonModule,
    FormsModule,
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
    TableModule,
    Checkbox,
  ],
  providers: [MessageService],
})
export class ProductFormComponent {
  @ViewChild('fileUpload') fileUpload: any;
  @Input() productData: Product | null = null;
  @Input() isEditMode: boolean = false;
  @Output() onSave = new EventEmitter<Product>();
  @Output() onCancelForm = new EventEmitter();
  checked = signal(false);

  productForm!: FormGroup;
  saving = signal(false);
  categories: Category[] = []; // Load your categories here
  currentImageUrl: string[] | null = null;
  selectedFiles = signal<File[]>([]);
  selectedImagePreview = signal<string[]>([]);


  categoryQueryData = categoryQuery();

  private productService = inject(ProductService);
  private firebaseService = inject(FirebaseService);
  private messageService = inject(MessageService);

  constructor(private fb: FormBuilder) {
    effect(() => {

      this.categories = this.categoryQueryData.data() ?? [];
    });
  }

  ngOnInit() {
    this.initializeForm();

    if (this.isEditMode && this.productData) {
      this.populateForm();
    }
  }

  initializeForm() {
    this.productForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      description: ['', [Validators.required, Validators.minLength(10)]],
      price: [0, [Validators.required, Validators.min(0.01)]],
      categoryId: ['', Validators.required],
      stock: [0, [Validators.required, Validators.min(0)]],
      is_message_card: [false],
      is_add_on: [false],
      has_stems: [false],
      stems: this.fb.array([]),
    });
  }

  get stems() {
    return this.productForm.get('stems') as FormArray;
  }

  populateForm() {
    if (this.productData) {
      if (this.productData && this.productData.stems) {
        this.productData.stems.forEach((stem) => {
          this.stems.push(
            this.fb.group({
              stem_count: [stem.stem_count, Validators.required],
              price: [stem.price, [Validators.required, Validators.min(0.01)]],
            }),
          );
        });
      }

      this.productForm.patchValue({
        name: this.productData.name,
        description: this.productData.description,
        price: this.productData.price,
        categoryId: this.productData.category_id,
        stock: this.productData.stock_quantity,
        is_message_card: this.productData.is_message_card,
        is_add_on: this.productData.is_add_on,
        has_stems: this.productData.has_stems,
        stems: this.productData.stems,
      });
      this.currentImageUrl = this.productData.image_url;
    }
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

  removeImagePreview() {
    this.selectedImagePreview.set([]);
    this.selectedFiles.set([]);
    this.fileUpload.clear();
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
      this.saving.set(true);

      const productData: Product = {
        name: this.productForm.value.name,
        description: this.productForm.value.description,
        price: this.productForm.value.price,
        image_url: this.isEditMode ? this.productData!.image_url : [],
        category_id: this.productForm.value.categoryId,
        stock_quantity: this.productForm.value.stock,
        is_message_card: this.productForm.value.is_message_card,
        is_add_on: this.productForm.value.is_add_on,
        has_stems: this.productForm.value.has_stems,
        stems: this.stems.value,
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
            this.saving.set(false); // Set to false on success
            const message = this.isEditMode
              ? 'Product updated successfully!'
              : 'Product added successfully!';
            this.onSave.emit(response.data);
            this.productForm.reset();
          },
          error: (err) => {
            this.saving.set(false); // Set to false on error
            this.messageService.add({
              severity: 'error',
              summary: 'Operation failed',
              detail: err.message,
            });
            console.error('API Error:', err);
          },
        });
      } catch (error) {
        console.error('Error uploading image:', error);
        this.saving.set(false); // Set to false if image upload fails
      }
    }
  }

  onCancel() {
    this.onCancelForm.emit();
  }

  addStemRow() {
    this.stems.push(
      this.fb.group({
        stem_count: [null, Validators.required],
        price: [null, [Validators.required, Validators.min(0.01)]],
      }),
    );
  }

  removeStemRow(index: number) {
    this.stems.removeAt(index);
  }
}
