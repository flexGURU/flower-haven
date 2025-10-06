import { inject, Signal } from '@angular/core';
import { OrderService } from './order.service';
import { lastValueFrom } from 'rxjs';
import { injectQuery } from '@tanstack/angular-query-experimental';

export const ordersQuery = () => {
  const orderService = inject(OrderService);

  const query = injectQuery(() => ({
    queryKey: ['orders'],
    queryFn: () => lastValueFrom(orderService.getOrders()),
    staleTime: 300 * 1000,
    refetchOnWindowFocus: true,
  }));

  return query;
};

export const orderByIdQuery = (orderId: Signal<string>) => {
  const orderService = inject(OrderService);

  const query = injectQuery(() => ({
    queryKey: ['order', orderId()],
    queryFn: () => lastValueFrom(orderService.getOrderById(orderId())),
    staleTime: 300 * 1000,
    refetchOnWindowFocus: true,
    enabled: !!orderId,
  }));

  return query;
};
