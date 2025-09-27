import { Component, computed, effect, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ButtonModule } from 'primeng/button';
import { GalleriaModule } from 'primeng/galleria';
import { InputNumberModule } from 'primeng/inputnumber';
import { TabViewModule } from 'primeng/tabview';
import { RatingModule } from 'primeng/rating';
import { ActivatedRoute, Router } from '@angular/router';
import { initialProd, Product } from '../../../../../shared/models/models';
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
import { Breadcrumb, BreadcrumbModule } from 'primeng/breadcrumb';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { messageCardQuery } from '../../../../../shared/services/product.query';

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
    BreadcrumbModule,
    ProgressSpinnerModule,
  ],
})
export class ProductDetailComponent {
  product = signal<Product>(initialProd);
  quantity = 1;
  rating = 4.5;
  reviewCount = 24;
  averageRating = 4.5;
  addonsVisible = false;
  messageCardQuery = messageCardQuery();
  loading = signal(false);

  home = { icon: 'pi pi-home', routerLink: '/' };

  items = signal([
    { label: 'Products', icon: 'pi pi-fw pi-list', routerLink: '/products' },
  ]);

  selectedStemSize: any = null;
  stemSizes = [
    { label: '12 Stems', value: 12, price: 0 },
    { label: '24 Stems', value: 24, price: 500 },
    { label: '36 Stems', value: 36, price: 1000 },
    { label: '48 Stems', value: 48, price: 1500 },
    { label: '60 Stems', value: 60, price: 2000 },
    { label: '72 Stems', value: 72, price: 2500 },
  ];

  includeMessageCard = signal(false);
  selectedMessageCard = signal<Product | null>(null);

  giftMessage = '';
  purchaseType = 'onetime';

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
    private router: Router,
  ) {
    effect(() => {
      this.addMessageCardToCart();

      this.items.set([
        {
          label: 'Products',
          icon: '',
          routerLink: '/products',
        },
        {
          label: this.product().name,
          icon: '',
          routerLink: `/product/${this.product().id}`,
        },
      ]);
    });
  }

  private productService = inject(ProductService);
  private cartService = inject(CartService);
  private messageService = inject(MessageService);

  messageCards = computed<Product[]>(() => this.messageCardQuery.data() ?? []);

  ngOnInit() {
    this.route.params.subscribe((params) => {
      this.loading.set(true);

      const productId = params['id'];
      if (productId) {
        this.loadProduct(productId);
        this.loading.set(false);
      }
    });

    this.selectedStemSize = this.stemSizes[0];
  }

  loadProduct(id: string) {
    if (id) {
      this.productService.getProductById(id).subscribe((response) => {
        this.product.set(response);
      });
      this.setupGallery();
    }
  }

  setupGallery() {
    if (this.product) {
      this.galleryImages = this.product().image_url;
    }
  }

  addToCart() {
    if (this.product()) {
      const cartItem = {
        product: this.product(),
        quantity: this.quantity,
        stemSize: this.selectedStemSize,
        messageCard: this.includeMessageCard()
          ? this.selectedMessageCard
          : null,
        giftMessage: this.giftMessage,
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

    let total = this.product().price * this.quantity;

    if (this.selectedStemSize) {
      total += this.selectedStemSize.price;
    }

    if (this.includeMessageCard() && this.selectedMessageCard()) {
      total += this.selectedMessageCard()?.price ?? 0;
    }

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

  addMessageCardToCart() {
    const card = this.selectedMessageCard();
    if (card) {
      this.cartService.addToCart(card);
    }
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
