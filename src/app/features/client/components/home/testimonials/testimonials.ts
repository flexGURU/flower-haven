import { Component } from '@angular/core';
import { CarouselModule } from 'primeng/carousel';

@Component({
  selector: 'app-testimonials',
  templateUrl: './testimonials.html',
  imports: [CarouselModule],
})
export class TestimonialsComponent {
  testimonials = [
    {
      text: 'The flowers were absolutely beautiful and lasted for weeks! Will definitely order again.',
      name: 'Wanjiku Muthui',
      location: 'Nairobi, Kenya',
    },
    {
      text: 'Amazing service and the subscription plan is so convenient. Fresh flowers every week!',
      name: 'Harry Odongo',
      location: 'Kisumu, Kenya',
    },
    {
      text: 'Perfect for surprising my wife. The arrangement was exactly what I wanted.',
      name: 'Eunice Wanjiru',
      location: 'Mombasa, Kenya',
    },
    {
      text: 'Such a delightful experience. The bouquet was fresh and elegantly arranged.',
      name: 'Dennis Omari',
      location: 'Nakuru, Kenya',
    },
  ];
  carouselResponsiveOptions = [
    {
      breakpoint: '1024px',
      numVisible: 2,
      numScroll: 1,
    },
    {
      breakpoint: '768px',
      numVisible: 1,
      numScroll: 1,
    },
  ];
}
