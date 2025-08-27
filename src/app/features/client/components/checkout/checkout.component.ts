import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
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
  private readonly phoneNumber = '254799335366';

  @Input() sumOrderDetails!: Cart;

  orderDetails!: FormGroup;

  constructor(private fb: FormBuilder) {
    this.minDate = new Date();

    this.orderDetails = this.fb.group({
      fullName: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      phoneNumber: ['', Validators.required],
      location: ['', Validators.required],
      address: ['', Validators.required],
      rangeDates: [null, Validators.required],
      selectedTimeSlot: [null, Validators.required],
      total: 0,
    });
  }

  ngOnInit() {
    this.timeSlots = this.generateTimeSlots();
    this.orderDetails.patchValue({ total: this.sumOrderDetails.total });
  }

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
  submitOrder() {
    const orderForm = this.orderDetails.getRawValue();
    const cart = this.sumOrderDetails;

    let message = `🛒 *New Order Request* \n\n`;

    message += `👤 *Customer Details*:\n`;
    message += `Full Name: ${orderForm.fullName}\n`;
    message += `Email: ${orderForm.email}\n`;
    message += `Phone: ${orderForm.phoneNumber}\n`;
    message += `Location: ${orderForm.location}\n`;
    message += `Address: ${orderForm.address}\n`;
    message += `Delivery Range: ${orderForm.rangeDates}\n`;
    message += `Time Slot: ${orderForm.selectedTimeSlot}\n\n`;

    message += `📦 *Cart Details*:\n`;
    cart.items.forEach((item: any, index: number) => {
      message += `${index + 1}. ${item.name} (x${item.quantity}) - ${item.price}\n`;
    });

    message += `\n *Total*: ${cart.total}`;

    const encodedMessage = encodeURIComponent(message);

    // ✅ Directly open WhatsApp app if available
    window.location.href = `https://wa.me/${this.phoneNumber}?text=${encodedMessage}`;
  }
}
