import { Component } from '@angular/core';
import { MenubarModule } from 'primeng/menubar';

@Component({
  selector: 'app-nav',
  templateUrl: './nav.html',
  imports: [MenubarModule],
})
export class NavComponent {
  items!: any;
  ngOnInit() {
    this.items = [
      {
        label: 'Dashboard',
        icon: 'pi pi-home',
        routerLink: 'dashboard',
      },
      {
        label: 'Category',
        icon: 'pi pi-folder',
        routerLink: 'categories',
      },

      {
        label: 'Product list',
        icon: 'pi pi-tags',
        routerLink: 'products',
      },
    ];
  }
}
