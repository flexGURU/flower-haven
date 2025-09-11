import { inject } from '@angular/core';
import { ProductService } from './product.service';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';

export const productQuery = () => {
  const productService = inject(ProductService);
  const PRODUCTSERVICEQUERYKEY = ['products'];

  const query = injectQuery(() => ({
    queryKey: PRODUCTSERVICEQUERYKEY,
    queryFn: () => lastValueFrom(productService.fetchProducts()),
    staleTime: 300 * 1000,
    refetchOnWindowFocus: true
  }));

  return query;
};
