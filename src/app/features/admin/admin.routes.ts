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
        component: DashboardComponent,
        canActivate: [authGuard],
      },
      {
        path: 'products',
        component: ProductManagement,
        canActivate: [authGuard],
      },
      {
        path: 'categories',
        component: CategoryManagement,
        canActivate: [authGuard],
      },
    ],
  },
];
