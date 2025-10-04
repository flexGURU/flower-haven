import { Component, effect, inject, signal } from '@angular/core';
import { TableModule } from 'primeng/table';
import { ordersQuery } from '../../../../../shared/services/order.query';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { MessageModule } from 'primeng/message';
import { CommonModule } from '@angular/common';
import { OrderPayload } from '../../../../../shared/models/models';
import { Router } from '@angular/router';

@Component({
  selector: 'app-order-management',
  imports: [TableModule, ProgressSpinnerModule, MessageModule, CommonModule],
  templateUrl: './order-management.component.html',
  styleUrl: './order-management.component.css',
})
export class OrderManagementComponent {
  ordersData = ordersQuery();

  #router = inject(Router);

  onOrderSelect(event: any) {
    this.#router.navigate(['/admin/order-detail', event.data.id]);
  }
}
