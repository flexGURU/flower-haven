import {
  Component,
  effect,
  EventEmitter,
  inject,
  Input,
  Output,
  signal,
  ViewChild,
} from '@angular/core';
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
import { FileUpload } from 'primeng/fileupload';
import { FirebaseService } from '../../services/firebase';
import { ProgressSpinner } from 'primeng/progressspinner';
import { Message } from 'primeng/message';

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
    FileUpload,
    ProgressSpinner,
    Message,
  ],
  providers: [MessageService],
})
export class CategoryForm {
  @ViewChild('fileUpload') fileUpload: any;

  @Input() categoryData: Category | null = null;
  @Input() isEditMode: boolean = false;
  @Output() onSave = new EventEmitter<Category>();
  @Output() onCancel = new EventEmitter<void>();
  statusMessage: string | null = null;
  imageUrl = signal('');
  selectedFiles = signal<File[]>([]);
  selectedImagePreview = signal<string[]>([]);
  categoryForm!: FormGroup;
  saving = false;
  currentImageUrl: string[] | null = null;
  uploadInProgress = signal(false);
  imageUploadStatus = signal<'success' | 'error' | ''>('');

  private productService = inject(ProductService);
  private firebaseService = inject(FirebaseService);

  constructor(
    private fb: FormBuilder,
    private messageService: MessageService,
  ) {
    effect(() => {
      if (this.selectedFiles().length > 0) {
        this.uploadImage(this.selectedFiles()).then((url) => {
          if (url) {
            this.imageUrl.set(url);
            this.categoryForm.patchValue({ imageUrl: url });
          }
        });
      }
    });
  }

  ngOnInit() {
    this.initializeForm();

    if (this.isEditMode && this.categoryData) {
      this.populateForm();
    }
  }

  async uploadImage(file: File[]): Promise<string | null> {
    try {
      this.uploadInProgress.set(true);
      const imageUrl = await this.firebaseService.uploadImage(file);
      this.uploadInProgress.set(false);
      this.imageUploadStatus.set('success');
      return imageUrl;
    } catch (error) {
      this.uploadInProgress.set(false);
      this.imageUploadStatus.set('error');
      this.messageService.add({
        severity: 'error',
        summary: 'Image Upload Failed',
        detail: 'There was an error uploading the image. Please try again.',
      });
      console.error('Image upload failed:', error);
      return null;
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
      this.imageUrl.set(this.categoryData.image_url?.[0] || '');
      this.categoryForm.patchValue({
        name: this.categoryData.name,
        description: this.categoryData.description,
        imageUrl: this.categoryData.image_url?.[0] || '',
      });
      this.currentImageUrl = this.categoryData.image_url || null;
    }
  }

  async onSubmit() {
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
}
