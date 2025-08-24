import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { GalleriaModule } from 'primeng/galleria';
import { InputNumberModule } from 'primeng/inputnumber';
import { TabViewModule } from 'primeng/tabview';
import { RatingModule } from 'primeng/rating';
import { ActivatedRoute, Router } from '@angular/router';
import { Product } from '../../../../../shared/models/models';
import { ProductService } from '../../../../../shared/services/product.service';
import { CartService } from '../../cart/cart.service';
import { Image } from 'primeng/image';
import { ProductAddonsComponent } from '../product-addons/product-addons';
import { DialogModule } from 'primeng/dialog';
import { DropdownModule } from 'primeng/dropdown';
import { CheckboxModule } from 'primeng/checkbox';
import { TextareaModule } from 'primeng/textarea';
import { RadioButtonModule } from 'primeng/radiobutton';
import { CardModule } from 'primeng/card';
import { CommonModule } from '@angular/common';
import { SelectModule } from 'primeng/select';
import { CartItem } from '../../cart/cart.model';
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-product-detail',
  templateUrl: './product_detail.html',
  imports: [
    CommonModule,
    GalleriaModule,
    ButtonModule,
    FormsModule,
    TabViewModule,
    RatingModule,
    InputNumberModule,
    Image,
    ProductAddonsComponent,
    DialogModule,
    SelectModule,
    CheckboxModule,
    TextareaModule,
    RadioButtonModule,
    CardModule,
  ],
})
export class ProductDetailComponent {
  product: Product | null = null;
  quantity = 1;
  rating = 4.5;
  reviewCount = 24;
  averageRating = 4.5;
  addonsVisible = false;

  selectedStemSize: any = null;
  stemSizes = [
    { label: '12 Stems', value: 12, price: 0 },
    { label: '24 Stems', value: 24, price: 500 },
    { label: '36 Stems', value: 36, price: 1000 },
    { label: '48 Stems', value: 48, price: 1500 },
    { label: '60 Stems', value: 60, price: 2000 },
    { label: '72 Stems', value: 72, price: 2500 },
  ];

  includeMessageCard = false;
  selectedMessageCard: any = null;
  messageCards = [
    { label: 'Birthday Card', value: 'birthday', price: 100 },
    { label: 'Anniversary Card', value: 'anniversary', price: 100 },
    { label: 'Thank You Card', value: 'thankyou', price: 100 },
    { label: 'Love Card', value: 'love', price: 100 },
    { label: 'Get Well Soon Card', value: 'getwell', price: 100 },
  ];

  giftMessage = '';
  purchaseType = 'onetime';

  addons = [
    {
      id: 1,
      name: 'Premium Red Wine',
      price: 2500,
      image: '/red-wine-bottle.png',
      selected: false,
    },
    {
      id: 2,
      name: 'Chocolate Box',
      price: 800,
      image: '/assorted-chocolates-gift-box.png',
      selected: false,
    },
    {
      id: 3,
      name: 'Teddy Bear',
      price: 1200,
      image: '/cozy-teddy-bear.png',
      selected: false,
    },
    {
      id: 4,
      name: 'Champagne',
      price: 3500,
      image: '/champagne-bottle.png',
      selected: false,
    },
  ];

  galleryImages: any;
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

  constructor(
    private route: ActivatedRoute,
    private router: Router, // Added router for navigation
  ) {}

  private productService = inject(ProductService);
  private cartService = inject(CartService);
  private messageService = inject(MessageService);

  ngOnInit() {
    this.route.params.subscribe((params) => {
      const productId = params['id'];
      if (productId) {
        this.loadProduct(productId);
      }
    });

    this.selectedStemSize = this.stemSizes[1]; // Default to 24 stems
  }

  loadProduct(id: string) {
    if (id) {
      this.productService.getProductById(id).subscribe((response) => {
        this.product = response;
      });
      this.setupGallery();
    }
  }

  setupGallery() {
    if (this.product) {
      this.galleryImages = this.product.image_url;
    }
  }

  addToCart() {
    if (this.product) {
      const cartItem = {
        product: this.product,
        quantity: this.quantity,
        stemSize: this.selectedStemSize,
        messageCard: this.includeMessageCard ? this.selectedMessageCard : null,
        giftMessage: this.giftMessage,
        addons: this.addons.filter((addon) => addon.selected),
        purchaseType: this.purchaseType,
      };

      this.cartService.addToCart(cartItem.product, cartItem.quantity);
      this.messageService.add({
        severity: 'success',
        summary: 'Info',
        detail: `${cartItem.product.name} added to cart`,
      });
    }
  }

  getTotalPrice(): number {
    if (!this.product) return 0;

    let total = this.product.price * this.quantity;

    if (this.selectedStemSize) {
      total += this.selectedStemSize.price;
    }

    if (this.includeMessageCard && this.selectedMessageCard) {
      total += this.selectedMessageCard.price;
    }

    const addonTotal = this.addons
      .filter((addon) => addon.selected)
      .reduce((sum, addon) => sum + addon.price, 0);
    total += addonTotal;

    return total;
  }

  handleSubscription() {
    if (this.purchaseType === 'subscription') {
      this.router.navigate(['/login'], {
        queryParams: { returnUrl: '/subscription-plans' },
      });
    }
  }

  toggleAddon(addon: any) {
    addon.selected = !addon.selected;
  }

  isFlowerProduct(): boolean {
    return true;
  }

  onAddonAddedToCart(product: Product) {
    if (product) {
      this.cartService.addToCart(product);
      this.messageService.add({
        severity: 'success',
        summary: 'Info',
        detail: `${product.name} added to cart`,
      });
    }
  }
  onAddonRemovedFromCart(product: Product) {
    if (product.id) {
      this.cartService.removeFromCart(product.id);
      this.messageService.add({
        severity: 'info',
        summary: 'Info',
        detail: `${product.name} removed from cart`,
      });
    }
  }
}
