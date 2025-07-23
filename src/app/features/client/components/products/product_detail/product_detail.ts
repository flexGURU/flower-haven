import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { GalleriaModule } from 'primeng/galleria';
import { InputNumberModule } from 'primeng/inputnumber';
import { TabViewModule } from 'primeng/tabview';
import { RatingModule } from 'primeng/rating';
import { ActivatedRoute } from '@angular/router';
import { Product } from '../../../../../shared/models/models';
import { ProductService } from '../../../../../shared/services/product.service';
import { CartService } from '../../cart/cart.service';

@Component({
  selector: 'app-product-detail',
  templateUrl: './product_detail.html',
  imports: [
    GalleriaModule,
    ButtonModule,
    FormsModule,
    TabViewModule,
    RatingModule,
    InputNumberModule,
  ],
})
export class ProductDetailComponent {
  product: Product | null = null;
  quantity = 1;
  rating = 4.5;
  reviewCount = 24;
  averageRating = 4.5;

  galleryImages: any ;
  galleryResponsiveOptions = [
    {
      breakpoint: '1024px',
      numVisible: 3,
    },
    {
      breakpoint: '768px',
      numVisible: 2,
    },
    {
      breakpoint: '560px',
      numVisible: 1,
    },
  ];

  sampleReviews = [
    {
      name: 'Sarah M.',
      rating: 5,
      date: '2 days ago',
      comment:
        'Absolutely beautiful flowers! They lasted much longer than expected.',
    },
    {
      name: 'John D.',
      rating: 4,
      date: '1 week ago',
      comment: 'Great quality and fast delivery. Will order again.',
    },
  ];

  constructor(
    private productService: ProductService,
    private cartService: CartService,
    private route: ActivatedRoute
  ) {}

  ngOnInit() {
    this.route.params.subscribe((params) => {
      const productId = params['id'];
      if (productId) {
        this.loadProduct(productId);
      }
    });
  }

  loadProduct(id: string) {
    if (id) {
      this.product = this.productService.getProductById(id);
      this.setupGallery();
    }
  }

  setupGallery() {
    if (this.product) {
      this.galleryImages = this.product.image
    }
  }

  addToCart() {
    if (this.product) {
      this.cartService.addToCart(this.product, this.quantity);
    }
  }
}
