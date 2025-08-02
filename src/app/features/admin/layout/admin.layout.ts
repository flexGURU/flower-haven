import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AdminHeaderComponent } from "./header/admin_header";

@Component({
  selector: 'app-admin-layout',
  templateUrl: './admin.layout.html',
  imports: [RouterOutlet, AdminHeaderComponent],
})
export class AdminLayout {}
