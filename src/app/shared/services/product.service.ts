import { Injectable, signal } from '@angular/core';
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
  private productsSubject = new BehaviorSubject<Product[]>([]);
  private categorySubject = new BehaviorSubject<Category[]>([]);
  private addOnsSubject = new BehaviorSubject<Product[]>([]);
  private messageCardSubject = new BehaviorSubject<Product[]>([]);

  constructor(private http: HttpClient) {
    this.fecthCategories().subscribe();
    this.fecthProducts().subscribe();
    this.fetchAddOns().subscribe();
    this.fetchMessageCards().subscribe();
  }

  products$ = this.productsSubject.asObservable();

  getProducts() {
    return this.productsSubject.asObservable();
  }

  getAddOns() {
    return this.addOnsSubject.asObservable();
  }

  getMessageCards() {
    return this.messageCardSubject.asObservable();
  }

  fecthProducts(): Observable<{ data: Product[] }> {
    return this.http.get<{ data: Product[] }>(this.productApiUrl).pipe(
      tap((response) => {
        this.productsSubject.next(response.data);
      }),
      catchError((error) => {
        console.error('Error fetching products:', error);
        return throwError(() => new Error('Failed to fetch products.'));
      }),
    );
  }

  addProduct(product: Product): Observable<{ data: Product }> {
    return this.http.post<{ data: Product }>(this.productApiUrl, product).pipe(
      tap(() => {
        this.fecthProducts().subscribe();
      }),
      catchError((error) => {
        console.error('Error creating product:', error);
        return throwError(() => new Error('Failed to create product.'));
      }),
    );
  }

  updateProduct(product: Product, id: string): Observable<{ data: Product }> {
    console.log('dddd', product);

    return this.http
      .put<{ data: Product }>(`${this.productApiUrl}/${id}`, product)
      .pipe(
        tap(() => {
          this.fecthProducts().subscribe();
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
          this.fecthProducts().subscribe();
        }),
      );
  }

  getCategories() {
    return this.categorySubject.asObservable();
  }

  fecthCategories(): Observable<{ data: Category[] }> {
    return this.http.get<{ data: Category[] }>(this.categoryApiUrl).pipe(
      tap((response) => {
        this.categorySubject.next(response.data);
      }),
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
          this.fecthCategories().subscribe();
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
          this.fecthCategories().subscribe();
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
          this.fecthCategories().subscribe();
        }),
      );
  }
  fetchAddOns(): Observable<Product[]> {
    return this.http.get<{ data: Product[] }>(`${this.productApiUrl}`).pipe(
      tap((response) => {
        const addons = response.data.filter(
          (product) => product.category_data?.name === 'Add-Ons',
        );
        this.addOnsSubject.next(addons);
      }),
      map((response) =>
        response.data.filter(
          (product) => product.category_data?.name === 'Add-Ons',
        ),
      ),
    );
  }

  fetchMessageCards(): Observable<{ data: Product[] }> {
    return this.http.get<{ data: Product[] }>(`${this.productApiUrl}`).pipe(
      tap((response) => {
        const messageCards = response.data.filter(
          (product) => product.category_data?.name === 'Message Cards',
        );
        this.messageCardSubject.next(messageCards);
      }),
    );
  }
}
