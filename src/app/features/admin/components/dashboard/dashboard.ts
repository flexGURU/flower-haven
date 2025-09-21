import { Component } from '@angular/core';
import { BadgeModule } from 'primeng/badge';
import { CardModule } from 'primeng/card';
import { TableModule } from 'primeng/table';
import { ChartModule } from 'primeng/chart';
import { CommonModule } from '@angular/common';
import { dashboardQuery } from '../../services/dashboard.query';
import { ButtonModule } from 'primeng/button';
import { ProgressSpinnerModule } from 'primeng/progressspinner';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.html',
  imports: [
    CardModule,
    BadgeModule,
    TableModule,
    ChartModule,
    CommonModule,
    ButtonModule,
    ProgressSpinnerModule,
  ],
})
export class DashboardComponent {
  stats = dashboardQuery();

  ngOnInit() {}

  getOrderStatusSeverity(
    status: string,
  ): 'info' | 'success' | 'warn' | 'danger' | 'secondary' | 'contrast' {
    switch (status) {
      case 'delivered':
        return 'success';
      case 'processing':
        return 'info';
      case 'shipped':
        return 'warn';
      case 'pending':
        return 'secondary';
      case 'cancelled':
        return 'danger';
      default:
        return 'secondary';
    }
  }
  stockSummary = [
    { name: 'Floral Vase', sku: 'FLR123', quantity: 120 },
    { name: 'Wedding Bouquet', sku: 'WED456', quantity: 18 },
    { name: 'Succulent Pack', sku: 'SUC789', quantity: 3 },
  ];

  getStockStatus(quantity: number): string {
    if (quantity > 20) return 'In Stock';
    if (quantity > 10) return 'Low Stock';
    return 'Out of Stock';
  }

  getStockSeverity(
    quantity: number,
  ): 'info' | 'success' | 'warn' | 'danger' | 'secondary' | 'contrast' {
    if (quantity > 20) return 'success';
    if (quantity > 10) return 'warn';
    return 'danger';
  }
}
