import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'admin',
    loadChildren: () =>
      import('./features/admin/admin.routes').then((m) => m.adminRoutes),
  },
  {
    path: '',
    loadChildren: () =>
      import('./features/client/client.routes').then((m) => m.clientRoutes),
  },
  {
    path: 'login',
    loadComponent: () =>
      import('./features/admin/core/login/login').then((m) => m.LoginComponent),
  },
];
