import { Routes } from '@angular/router';
import { AdminLayout } from './layout/admin.layout';
import { DashboardComponent } from './components/dashboard/dashboard';

export const adminRoutes: Routes = [
  {
    path: '',
    component: AdminLayout,
    children: [
      {
        path: 'dashboard',
        component: DashboardComponent
      },
    ],
  },
];
