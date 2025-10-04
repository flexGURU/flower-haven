import { Component, effect, input, Input } from '@angular/core';
import { orderByIdQuery } from '../../../../../shared/services/order.query';
import { CommonModule } from '@angular/common';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { MessageModule } from 'primeng/message';
import { TagModule } from 'primeng/tag';
import { TableModule } from 'primeng/table';

@Component({
  selector: 'app-order-detail',
  imports: [
    CommonModule,
    ProgressSpinnerModule,
    MessageModule,
    TagModule,
    TableModule,
  ],
  templateUrl: './order-detail.component.html',
  styleUrl: './order-detail.component.css',
})
export class OrderDetailComponent {
  orderId = input.required<string>({ alias: 'id' });

  orderDetail = orderByIdQuery(this.orderId);
}
