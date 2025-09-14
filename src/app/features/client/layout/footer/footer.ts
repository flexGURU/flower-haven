import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Component } from '@angular/core';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.html',
  imports: [CommonModule, ReactiveFormsModule, FormsModule],
})
export class FoooterComponent {
  newsletterEmail = '';
  subscribeNewsletter() {}


  date = new Date();
}
