import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { AuthService } from '../../../../core/auth/auth.service';

@Component({
  selector: 'app-admin-header',
  templateUrl: './admin_header.html',
  imports: [CommonModule, ButtonModule],
})
export class AdminHeaderComponent {
  private authService = inject(AuthService);
  logOut() {
    this.authService.logout();
  }
}
