import { computed, effect, inject, Injectable, signal } from '@angular/core';
import { Category, Product } from '../models/models';
import { catchError, map, Observable, tap, throwError } from 'rxjs';
import { apiUrl } from '../../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class ProductService {
  private readonly categoryApiUrl = `${apiUrl}/categories`;
  private readonly productApiUrl = `${apiUrl}/products`;
  private router = inject(Router);

  page = signal(1);
  limit = signal(15);
  search = signal('');
  priceFrom = signal<number | null>(null);
  priceTo = signal(0);
  categoryId = signal<string[] | []>([]);
  totalProducts = signal(0);
  totalAddOns = signal(0);
  is_add_on = signal(false);
  is_message_card = signal(false);

  initialProductFilters = {
    page: this.page(),
    limit: this.limit(),
    search: this.search(),
    priceFrom: this.priceFrom(),
    priceTo: this.priceTo(),
    categoryId: this.categoryId(),
    is_add_on: this.is_add_on(),
    is_message_card: this.is_message_card(),
  };

  productBaseApiUrl = computed(() => {
    const params = new URLSearchParams();

    if (this.page()) params.set('page', this.page().toString());
    if (this.limit()) params.set('limit', this.limit().toString());
    if (this.search()) params.set('search', this.search());

    if (this.categoryId()) {
      this.categoryId()!.forEach((id) => params.append('category_id', id));
    }

    params.set('is_add_on', this.is_add_on()!.toString());
    params.set('is_message_card', this.is_message_card()!.toString());

    return `${this.productApiUrl}?${params.toString()}`;
  });

  constructor(private http: HttpClient) {
    this.fetchMessageCards().subscribe();
    effect(() => {});
    const currentRoute = signal(this.router.url);
  }

  fetchProducts(): Observable<Product[]> {
    return this.http
      .get<{
        data: Product[];
        pagination: { total: number };
      }>(`${this.productBaseApiUrl()}`)
      .pipe(
        tap((response) => {
          this.totalProducts.set(response.pagination.total);
        }),
        map((response) => response.data),
        catchError((error) => {
          console.error('Error fetching products:', error);
          return throwError(() => new Error('Failed to fetch products.'));
        }),
      );
  }

  fetchAddOns(): Observable<Product[]> {
    return this.http
      .get<{ data: Product[] }>(`${this.productApiUrl}/add-ons`)
      .pipe(
        map((response) => {
          return response.data;
        }),
        catchError((error) => {
          console.error('Error fetching add-ons:', error);
          return throwError(() => new Error('Failed to fetch add-ons.'));
        }),
      );
  }

  fetchMessageCards(): Observable<Product[]> {
    return this.http
      .get<{ data: Product[] }>(`${this.productApiUrl}/message-cards`)
      .pipe(
        map((response) => {
          return response.data;
        }),
        catchError((error) => {
          console.error('Error fetching message cards:', error);
          return throwError(() => new Error('Failed to fetch message cards.'));
        }),
      );
  }

  addProduct(product: Product): Observable<{ data: Product }> {
    return this.http.post<{ data: Product }>(this.productApiUrl, product).pipe(
      tap(() => {
        this.fetchProducts().subscribe();
      }),
      catchError((error) => {
        console.error('Error creating product:', error);
        return throwError(() => new Error('Failed to create product.'));
      }),
    );
  }

  updateProduct(product: Product, id: string): Observable<{ data: Product }> {
    return this.http
      .put<{ data: Product }>(`${this.productApiUrl}/${id}`, product)
      .pipe(
        tap(() => {
          this.fetchProducts().subscribe();
        }),
        catchError((error) => {
          console.error('Error updating product:', error);
          return throwError(() => new Error('Failed to update product.'));
        }),
      );
  }

  getProductById(id: string): Observable<Product> {
    return this.http
      .get<{ data: Product }>(`${this.productApiUrl}/${id}`)
      .pipe(map((response) => response.data));
  }

  deleteProduct(id: string): Observable<{ message: string }> {
    return this.http
      .delete<{ message: string }>(`${this.productApiUrl}/${id}`)
      .pipe(
        tap(() => {
          this.fetchProducts().subscribe();
        }),
      );
  }

  fetchCategories(): Observable<Category[]> {
    return this.http.get<{ data: Category[] }>(this.categoryApiUrl).pipe(
      map((response) => response.data),
      catchError((error) => {
        console.error('Error fetching categories:', error);
        return throwError(() => new Error('Failed to fetch categories.'));
      }),
    );
  }

  addCategory(category: Category): Observable<{ data: Category }> {
    return this.http
      .post<{ data: Category }>(this.categoryApiUrl, category)
      .pipe(
        tap(() => {
          this.fetchCategories().subscribe();
        }),
      );
  }

  updateCategory(category: Category): Observable<{ data: Category }> {
    return this.http
      .put<{
        data: Category;
      }>(`${this.categoryApiUrl}/${category.id}`, category)
      .pipe(
        tap(() => {
          this.fetchCategories().subscribe();
        }),
        catchError((error) => {
          console.error('Error updating category:', error);
          return throwError(() => new Error('Failed to update category.'));
        }),
      );
  }

  deleteCategory(id: string): Observable<{ message: string }> {
    return this.http
      .delete<{ message: string }>(`${this.categoryApiUrl}/${id}`)
      .pipe(
        tap(() => {
          this.fetchCategories().subscribe();
        }),
      );
  }

  getCategoryById(): Observable<Category> {
    return this.http
      .get<{ data: Category }>(`${this.categoryApiUrl}/${this.categoryId()}`)
      .pipe(
        map((response) => response.data),
        catchError((error) => {
          console.error('Error fetching category by ID:', error);
          return throwError(() => new Error('Failed to fetch category by ID.'));
        }),
      );
  }
}
