import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AdminHeaderComponent } from "./header/admin_header";
import { NavComponent } from "./nav/nav";

@Component({
  selector: 'app-admin-layout',
  templateUrl: './admin.layout.html',
  imports: [RouterOutlet, AdminHeaderComponent, NavComponent],
})
export class AdminLayout {}
