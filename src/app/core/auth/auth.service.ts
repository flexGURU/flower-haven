import { Injectable } from '@angular/core';
import { apiUrl } from '../../../environments/environment.development';
import { Observable, tap } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

interface UserResponse {
  auth: {
    access_token: string;
    refresh_token: string;
  };
  data: {
    id: number;
    name: string;
    email: string;
    phone_number: string;
    is_admin: boolean;
    is_active: boolean;
    created_at: string;
  };
}

interface LoginResponse {
  access_token: string;
  refresh_token: string;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly apiUrl = apiUrl;
  private static readonly accessToken = 'JWT_ACCESS_KEY';

  constructor(
    private http: HttpClient,
    private router: Router,
  ) {}
  get jwt(): string {
    return sessionStorage.getItem(AuthService.accessToken) ?? '';
  }

  private set jwt(value: string) {
    sessionStorage.setItem(AuthService.accessToken, value);
  }

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>(`${this.apiUrl}/user/login`, { email, password })
      .pipe(
        tap((resp) => {
          this.jwt = resp.access_token;
        }),
      );
  }

  isLoggedIn(): boolean {
    return !!this.jwt;
  }

  logout(): void {
    sessionStorage.removeItem(AuthService.accessToken);
    this.router.navigate(['/admin/login']);
  }
}
