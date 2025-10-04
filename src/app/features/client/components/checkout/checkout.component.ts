import { CommonModule } from '@angular/common';
import { Component, computed, inject, Input, signal } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { DropdownModule } from 'primeng/dropdown';
import { InputTextModule } from 'primeng/inputtext';
import { PaginatorModule } from 'primeng/paginator';
import { DatePicker } from 'primeng/datepicker';
import { SelectModule } from 'primeng/select';
import { Cart } from '../cart/cart.model';
import PaystackPop from '@paystack/inline-js';
import { PaystackService } from '../../services/paystack.service';
import { OrderPayload } from '../../../../shared/models/models';
import { CartSignalService } from '../cart/cart.signal.service';
import { OrderService } from '../../../../shared/services/order.service';
import { MessageService } from 'primeng/api';
import { Router } from '@angular/router';

@Component({
  selector: 'app-checkout',
  imports: [
    FormsModule,
    CommonModule,
    InputTextModule,
    DropdownModule,
    ButtonModule,
    PaginatorModule,
    DatePicker,
    InputTextModule,
    SelectModule,
    ReactiveFormsModule,
  ],
  templateUrl: './checkout.component.html',
  styleUrl: './checkout.component.css',
})
export class CheckoutComponent {
  minDate: Date;
  timeSlots: any[] = [];
  private readonly phoneNumber = '254794663008';
  today = new Date();
  loading = signal(false);

  orderDetails!: FormGroup;
  #payStackService = inject(PaystackService);
  #cartSignalService = inject(CartSignalService);
  #orderService = inject(OrderService);
  #messageService = inject(MessageService);
  #router = inject(Router);

  constructor(private fb: FormBuilder) {
    this.minDate = new Date();

    this.orderDetails = this.fb.group({
      fullName: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      phoneNumber: ['', Validators.required],
      location: ['', Validators.required],
      address: ['', Validators.required],
      expectedDeliveryDate: ['', Validators.required],
      selectedTimeSlot: ['', Validators.required],
    });
  }

  ngOnInit() {
    this.timeSlots = this.generateTimeSlots();
  }

  total = computed(() => this.#cartSignalService.cartTotal());

  generateTimeSlots(): any[] {
    const slots = [];
    for (let hour = 9; hour <= 17; hour++) {
      slots.push({
        label: `${this.formatHour(hour)}:00 - ${this.formatHour(hour + 1)}:00`,
        value: `${hour}:00`,
      });
    }
    return slots;
  }

  formatHour(hour: number): string {
    return hour > 12 ? `${hour - 12} PM` : `${hour} AM`;
  }

  confirmOrder() {
    const orderForm = this.orderDetails.getRawValue();
    const total = this.total();
    const items = this.#cartSignalService.cart();
    const date = new Date(orderForm.expectedDeliveryDate);

    const orderPayload = {
      'Full Name': orderForm.fullName,
      Email: orderForm.email,
      'Phone Number': orderForm.phoneNumber,
      Location: orderForm.location,
      Address: orderForm.address,
      'Expected Delivery Date': date.toISOString().split('T')[0],
      'Selected Time Slot': orderForm.selectedTimeSlot.label,
      'Total Amount': total,
      Items: this.#cartSignalService.cart().map((item) => {
        return {
          'Product ID': item.product.id,
          'Product Name': item.product.name,
          Quantity: item.quantity,
          Amount: item.amount,
        };
      }),
    };

    const message = encodeURIComponent(JSON.stringify(orderPayload, null, 2));
    const whatsappUrl = `https://wa.me/${this.phoneNumber}?text=${message}`;

    window.open(whatsappUrl, '_blank');
  }
  submitOrder() {
    this.loading.set(true);
    const orderForm = this.orderDetails.getRawValue();

    this.#payStackService
      .initializePayment(orderForm.email, this.total())
      .subscribe({
        next: (response) => {
          const popup = new PaystackPop();
          this.loading.set(false);
          popup.resumeTransaction(response.access_code, {
            onSuccess: (txn) => {
              const orderPayload = this.buildOrderPayload(txn.reference);
              this.saveOrder(orderPayload);
            },
            onError: (error) => {
              console.log('Transaction error:', error);
            },
          });
        },
      });
  }

  private saveOrder = (order: OrderPayload) => {
    this.#orderService.createOrder(order).subscribe({
      next: () => {
        this.#cartSignalService.clearCart();
        this.loading.set(false);
        this.handleTransactionStatus(true);
        this.#router.navigate(['/order-success']);
      },
      error: () => {
        this.loading.set(false);
        this.handleTransactionStatus(false);
      },
    });
  };

  private buildOrderPayload(tranxRef: string): OrderPayload {
    const orderForm = this.orderDetails.getRawValue();
    const date = new Date(orderForm.expectedDeliveryDate);
    let order: OrderPayload = {
      user_name: orderForm.fullName,
      user_phone_number: orderForm.phoneNumber,
      payment_status: false,
      status: 'active',
      delivery_date: date.toISOString().split('T')[0],
      time_slot: orderForm.selectedTimeSlot.label,
      total: this.total(),
      reference: tranxRef,
      items: this.#cartSignalService.cart().map((item) => ({
        product_id: item.product.id!,
        quantity: item.quantity,
        amount: item.amount,
        payment_method: 'one_time',
        frequency: 'weekly',
      })),
    };

    return order;
  }

  private handleTransactionStatus = (status: boolean) => {
    if (status) {
      this.#messageService.add({
        severity: 'success',
        summary: 'Payment Successful',
        detail: 'Your payment was successful and your order has been placed.',
      });
    } else {
      this.#messageService.add({
        severity: 'error',
        summary: 'Payment Failed',
        detail: 'There was an issue with your payment. Please try again.',
      });
    }
  };
}
