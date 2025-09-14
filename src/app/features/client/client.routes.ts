import { Routes } from '@angular/router';
import { ClientLayout } from './layout/client.layout';
import { HomeLayoutComponent } from './components/home/layout/layout';
import { ProductComponent } from './components/products/product_list/product';
import { ProductDetailComponent } from './components/products/product_detail/product_detail';
import { CartComponent } from './components/cart/cart';

export const clientRoutes: Routes = [
  {
    path: '',
    component: ClientLayout,
    children: [
      { path: '', component: HomeLayoutComponent },
      { path: 'products', component: ProductComponent },
      { path: 'product/:id', component: ProductDetailComponent },
      { path: 'cart', component: CartComponent },
      {
        path: 'contact',
        loadComponent: () =>
          import('./components/contact/contact.component').then(
            (m) => m.ContactComponent,
          ),
      },
       {
        path: 'about',
        loadComponent: () =>
          import('./components/about/about.component').then(
            (m) => m.AboutComponent,
          ),
      },
    ],
  },
];
