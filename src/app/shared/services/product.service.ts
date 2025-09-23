import { computed, effect, Injectable, signal } from '@angular/core';
import { Category, Product } from '../models/models';
import {
  BehaviorSubject,
  catchError,
  map,
  Observable,
  tap,
  throwError,
} from 'rxjs';
import { apiUrl } from '../../../environments/environment';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class ProductService {
  private readonly categoryApiUrl = `${apiUrl}/categories`;
  private readonly productApiUrl = `${apiUrl}/products`;
  private addOnsSubject = new BehaviorSubject<Product[]>([]);
  private messageCardSubject = new BehaviorSubject<Product[]>([]);

  page = signal(1);
  limit = signal(15);
  search = signal('');
  priceFrom = signal<number | null>(null);
  priceTo = signal(0);
  categoryId = signal<string[] | []>([]);
  totalProducts = signal(0);
  totalAddOns = signal(0);

  initialProductFilters = {
    page: this.page(),
    limit: this.limit(),
    search: this.search(),
    priceFrom: this.priceFrom(),
    priceTo: this.priceTo(),
    categoryId: this.categoryId(),
  };

  productBaseApiUrl = computed(() => {
    const params = new URLSearchParams();

    if (this.page()) params.set('page', this.page().toString());
    if (this.limit()) params.set('limit', this.limit().toString());
    if (this.search()) params.set('search', this.search());
    if (this.priceFrom())
      params.set('price_from', this.priceFrom()!.toString());
    if (this.priceTo()) params.set('price_to', this.priceTo()!.toString());
    if (this.categoryId()) {
      this.categoryId()!.forEach((id) => params.append('category_id', id));
    }

    return `${this.productApiUrl}?${params.toString()}`;
  });

  constructor(private http: HttpClient) {
    this.fetchMessageCards().subscribe();
    effect(() => {
    });
  }

  getAddOns() {
    return this.addOnsSubject.asObservable();
  }

  getMessageCards() {
    return this.messageCardSubject.asObservable();
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
          const addOns = response.data.filter(product => product.is_add_on);
          this.totalAddOns.set(addOns.length);
        }),
        map((response) => response.data),
        catchError((error) => {
          console.error('Error fetching products:', error);
          return throwError(() => new Error('Failed to fetch products.'));
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


  fetchMessageCards(): Observable<{ data: Product[] }> {
    return this.http.get<{ data: Product[] }>(`${this.productApiUrl}`).pipe(
      tap((response) => {
        const messageCards = response.data.filter(
          (product) => product.is_message_card === true,
        );

        this.messageCardSubject.next(messageCards);
      }),
    );
  }
}
