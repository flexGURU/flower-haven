import { inject, Injectable } from '@angular/core';
import { payStackConfig } from '../../../../environments/environment';
import { apiUrl } from '../../../../environments/environment.development';
import { catchError, map, Observable, throwError } from 'rxjs';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class PaystackService {
  private readonly apiUrl = apiUrl;
  #http = inject(HttpClient);

  constructor() {}

  initializePayment(email: string, amount: number): Observable<any> {
    return this.#http
      .post(`${this.apiUrl}/paystack/initialize`, {
        email: email,
        amount: amount,
      })
      .pipe(
        map((response: any) => {
          console.log('response', response);
          return response;
        }),
        catchError((error) => {
          console.error('Payment initialization error:', error);
          return throwError(() => new Error('Payment initialization failed'));
        }),
      );
  }

  verifyPayment(reference: string): Observable<any> {
    return this.#http
      .get<any>(`${this.apiUrl}/paystack/payments/${reference}`)
      .pipe(
        map((response) => {
          return response.data;
        }),
      );
  }
}
