import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-admin-header',
  templateUrl: './admin_header.html',
  imports: [CommonModule, ButtonModule],
})
export class AdminHeaderComponent {
  logOut() {}
}
