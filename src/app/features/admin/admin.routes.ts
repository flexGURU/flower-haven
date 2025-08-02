import { Routes } from '@angular/router';
import { AdminLayout } from './layout/admin.layout';
import { DashboardComponent } from './components/dashboard/dashboard';
import { ProductManagement } from './components/product_management/product_management';
import { CategoryComponent } from '../client/components/home/category/category';
import { CategoryManagement } from './components/category-management/category-management';
import { LoginComponent } from './core/login/login';

export const adminRoutes: Routes = [
  {
    path: '',
    component: AdminLayout,
    children: [
      {
        path: 'dashboard',
        component: DashboardComponent,
      },
      {
        path: 'products',
        component: ProductManagement,
      },
      {
        path: 'categories',
        component: CategoryManagement,
      },
    ],
  },
];
