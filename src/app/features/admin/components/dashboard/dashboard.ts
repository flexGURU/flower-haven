import { Component } from '@angular/core';
import { BadgeModule } from 'primeng/badge';
import { CardModule } from 'primeng/card';
import { TableModule } from 'primeng/table';
import { ChartModule } from 'primeng/chart';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.html',
  imports: [CardModule, BadgeModule, TableModule, ChartModule, CommonModule],
})
export class DashboardComponent {
  stats = {
    totalRevenue: 45230,
    totalOrders: 1234,
    activeSubscriptions: 456,
    totalCustomers: 2890,
  };

  recentOrders = [
    {
      id: '#12345',
      customerName: 'John Doe',
      total: 89.99,
      status: 'delivered',
    },
    {
      id: '#12346',
      customerName: 'Jane Smith',
      total: 156.5,
      status: 'processing',
    },
    {
      id: '#12347',
      customerName: 'Bob Johnson',
      total: 75.25,
      status: 'shipped',
    },
    {
      id: '#12348',
      customerName: 'Alice Brown',
      total: 234.0,
      status: 'pending',
    },
    {
      id: '#12349',
      customerName: 'Mike Wilson',
      total: 125.75,
      status: 'delivered',
    },
  ];

  topProducts = [
    {
      name: 'Rose Bouquet',
      sales: 145,
      revenue: 2890,
      image: '/assets/images/rose-bouquet.jpg',
    },
    {
      name: 'Tulip Arrangement',
      sales: 98,
      revenue: 1960,
      image: '/assets/images/tulip-arrangement.jpg',
    },
    {
      name: 'Mixed Flowers',
      sales: 87,
      revenue: 1740,
      image: '/assets/images/mixed-flowers.jpg',
    },
    {
      name: 'Orchid Plant',
      sales: 65,
      revenue: 1950,
      image: '/assets/images/orchid-plant.jpg',
    },
  ];

  revenueChartData = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
    datasets: [
      {
        label: 'Revenue',
        data: [12000, 19000, 15000, 25000, 22000, 30000],
        fill: false,
        borderColor: '#ec4899',
        backgroundColor: '#ec4899',
        tension: 0.4,
      },
    ],
  };

  ordersChartData = {
    labels: ['Delivered', 'Processing', 'Shipped', 'Pending', 'Cancelled'],
    datasets: [
      {
        data: [45, 25, 15, 10, 5],
        backgroundColor: [
          '#10b981',
          '#3b82f6',
          '#f59e0b',
          '#ef4444',
          '#6b7280',
        ],
      },
    ],
  };

  chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
    },
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  };

  doughnutOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom',
      },
    },
  };

  ngOnInit() {
    // Load dashboard data
  }

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
    if (quantity > 50) return 'In Stock';
    if (quantity > 10) return 'Low Stock';
    return 'Out of Stock';
  }

  getStockSeverity(
    quantity: number,
  ): 'info' | 'success' | 'warn' | 'danger' | 'secondary' | 'contrast' {
    if (quantity > 50) return 'success';
    if (quantity > 10) return 'warn';
    return 'danger';
  }
}
