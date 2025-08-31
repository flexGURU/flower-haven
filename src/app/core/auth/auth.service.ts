import { Injectable } from '@angular/core';
import { apiUrl } from '../../../environments/environment.development';
import { map, Observable, tap } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { User } from '../../shared/models/models';

interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly apiUrl = apiUrl;
  private static readonly accessToken = 'JWT_ACCESS_KEY';
  private static readonly userRole = 'USER_ROLE';

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

  private set role(value: string) {
    sessionStorage.setItem(AuthService.userRole, value);
  }

  get role(): string {
    return sessionStorage.getItem(AuthService.userRole) ?? '';
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

  isAdmin(): boolean {
    if (this.role) {
      return true;
    } else {
      return false;
    }
  }

  logout(): void {
    sessionStorage.removeItem(AuthService.accessToken);
    this.router.navigate(['/admin/login']);
  }

  signup = (user: User): Observable<User> => {
    return this.http.post<{ data: User }>(`${apiUrl}/users`, user).pipe(
      map((response) => {
        return response.data;
      }),
    );
  };
}
