import { Routes } from '@angular/router';
import { AdminLayout } from './layout/admin.layout';
import { DashboardComponent } from './components/dashboard/dashboard';
import { ProductManagement } from './components/product_management/product_management';
import { CategoryComponent } from '../client/components/home/category/category';
import { CategoryManagement } from './components/category-management/category-management';
import { LoginComponent } from '../../core/components/login/login';
import { authGuard } from '../../core/auth/auth.guard';

export const adminRoutes: Routes = [
  { path: 'login', component: LoginComponent },
  {
    path: '',
    component: AdminLayout,
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      {
        path: 'dashboard',
        canActivate: [authGuard],
        loadComponent: () =>
          import('./components/dashboard/dashboard').then(
            (m) => m.DashboardComponent,
          ),
      },
      {
        path: 'products',
        canActivate: [authGuard],
        loadComponent: () =>
          import('./components/product_management/product_management').then(
            (m) => m.ProductManagement,
          ),
      },
      {
        path: 'categories',
        canActivate: [authGuard],
        loadComponent: () =>
          import('./components/category-management/category-management').then(
            (m) => m.CategoryManagement,
          ),
      },
    ],
  },
];
