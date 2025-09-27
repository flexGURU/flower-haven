import { inject } from '@angular/core';
import { ProductService } from './product.service';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';

export const productQuery = () => {
  const productService = inject(ProductService);

  const query = injectQuery(() => ({
    queryKey: ['products', productService.productBaseApiUrl()],
    queryFn: () => lastValueFrom(productService.fetchProducts()),
    staleTime: 300 * 1000,
    refetchOnWindowFocus: true,
  }));

  return query;
};

export const categoryQuery = () => {
  const productService = inject(ProductService);

  const query = injectQuery(() => ({
    queryKey: ['categories'],
    queryFn: () => lastValueFrom(productService.fetchCategories()),
    staleTime: 300 * 1000,
    refetchOnWindowFocus: true,
  }));

  return query;
};

export const addonQuery = () => {
  const productService = inject(ProductService);

  const query = injectQuery(() => ({
    queryKey: ['addons'],
    queryFn: () => lastValueFrom(productService.fetchAddOns()),
    staleTime: 300 * 1000,
    refetchOnWindowFocus: true,
  }));

  return query;
};

export const messageCardQuery = () => {
  const productService = inject(ProductService);

  const query = injectQuery(() => ({
    queryKey: ['message-cards'],
    queryFn: () => lastValueFrom(productService.fetchMessageCards()),
    staleTime: 300 * 1000,
    refetchOnWindowFocus: true,
  }));

  return query;
};
