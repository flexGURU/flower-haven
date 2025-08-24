import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { DropdownModule } from 'primeng/dropdown';
import { InputTextModule } from 'primeng/inputtext';
import { PaginatorModule } from 'primeng/paginator';
import { DatePicker } from 'primeng/datepicker';
import { SelectModule } from 'primeng/select';

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
  ],
  templateUrl: './checkout.component.html',
  styleUrl: './checkout.component.css',
})
export class CheckoutComponent {
  location: string = '';
  address: string = '';
  rangeDates: Date | null = null;
  minDate: Date;
  timeSlots: any[] = [];
  selectedTimeSlot: any | null = null;
  expectedDeliveryDate!: Date;

  constructor() {
    this.minDate = new Date();
  }

  ngOnInit() {
    this.timeSlots = this.generateTimeSlots();
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
    const orderDetails = {
      location: this.location,
      address: this.address,
      deliveryDate: this.expectedDeliveryDate,
      timeSlot: this.selectedTimeSlot,
    };
    console.log('Order details submitted:', orderDetails);
    // Here, you would send this data to your backend
  }
}
