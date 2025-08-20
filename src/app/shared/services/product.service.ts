import { Injectable, signal } from '@angular/core';
import { Category, Product } from '../models/models';
import { BehaviorSubject, catchError, Observable, tap, throwError } from 'rxjs';
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

  constructor(private http: HttpClient) {
    this.fecthCategories().subscribe();
  }

  products$ = this.productsSubject.asObservable();

  createProduct(product: Product) {
    const products = this.productsSubject.value;
    this.productsSubject.next([...products, product]);
  }

  updateProduct(id: string, product: Product) {
    const products = this.productsSubject.value;
    const index = products.findIndex((p) => p.id === product.id);
    if (index !== -1) {
      products[index] = product;
      this.productsSubject.next([...products]);
    }
  }

  getProductById(id: string) {
    return this.productsSubject.value[0];
  }

  deleteProduct(id: string) {
    const products = this.productsSubject.value.filter((p) => p.id !== id);
    this.productsSubject.next(products);
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
      .patch<{
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
          console.log(`Category with id ${id} deleted successfully.`);
          this.fecthCategories().subscribe();
        }),
      );
  }
}
