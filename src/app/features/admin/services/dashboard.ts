import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { apiUrl } from '../../../../environments/environment.development';
import { map, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class DashboardService {
  #http = inject(HttpClient);

  readonly #apiUrl = apiUrl;

  getDashboardStats(): Observable<any> {
    return this.#http
      .get<any>(`${this.#apiUrl}/dashboard`)
      .pipe(map((response) => response.data));
  }
}
