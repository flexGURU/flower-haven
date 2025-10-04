import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { ProductService } from '../../../../../shared/services/product.service';

@Component({
  selector: 'app-hero',
  templateUrl: './hero.html',
  imports: [RouterLink, ButtonModule],
})
export class HeroComponent {
  #productService = inject(ProductService);

  browseProducts() {
    this.#productService.is_add_on.set(false);
    this.#productService.is_message_card.set(false);
  }
}
