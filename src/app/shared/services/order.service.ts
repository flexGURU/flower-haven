import { inject, Injectable } from '@angular/core';
import { apiUrl } from '../../../environments/environment.development';
import { HttpClient } from '@angular/common/http';
import { catchError, map, Observable, throwError } from 'rxjs';
import { OrderPayload } from '../models/models';

@Injectable({
  providedIn: 'root',
})
export class OrderService {
  private readonly apiUrl = apiUrl;

  #http = inject(HttpClient);

  createOrder = (orderPayload: OrderPayload) => {
    return this.#http.post(`${this.apiUrl}/orders`, orderPayload).pipe(
      map((response: any) => response.data),
      catchError((error) => {
        console.error('Error creating order:', error);
        return throwError(() => new Error('Order creation failed'));
      }),
    );
  };

  getOrderById = (orderId: string) => {
    return this.#http.get(`${this.apiUrl}/orders/${orderId}`).pipe(
      map((response: any) => {
        return response.data;
      }),
      catchError((error) => {
        console.error('Error fetching order:', error);
        return throwError(() => new Error('Fetching order failed'));
      }),
    );
  };

  getOrders = (): Observable<OrderPayload[]> => {
    return this.#http
      .get<{ data: OrderPayload[] }>(`${this.apiUrl}/orders`)
      .pipe(
        map((response: any) => response.data),
        catchError((error) => {
          console.error('Error fetching orders:', error);
          return throwError(() => new Error('Fetching orders failed'));
        }),
      );
  };
}
