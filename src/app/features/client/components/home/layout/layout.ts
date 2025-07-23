import { Component } from '@angular/core';
import { HomeComponent } from "../homepage/home";

@Component({
  selector: 'app-home-layout',
  templateUrl: './layout.html',
  imports: [HomeComponent],
})
export class HomeLayoutComponent {
  constructor() {}

  ngOnInit() {
  }
}
