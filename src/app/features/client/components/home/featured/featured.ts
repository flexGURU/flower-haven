import {
  Component,
  AfterViewInit,
  ChangeDetectionStrategy,
  inject,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ProductService } from '../../../../../shared/services/product.service';
import { Router } from '@angular/router';

// The App component encapsulates the entire service showcase section,
// including its template, custom styles, and scroll animation logic.
@Component({
  selector: 'app-featured',
  standalone: true,
  imports: [CommonModule],
  // Styles are included directly in the component metadata
  styles: [
    `
      /* --- Custom Styles from HTML/CSS --- */
      @import url('https://fonts.googleapis.com/css2?family=Inter:wght@100..900&display=swap');

      :host {
        display: block;
      }

      .app-root-container {
        font-family: 'Inter', sans-serif;
        min-height: 120vh; /* Ensures enough height for scrolling to test observer */
      }

      /* Base styles for the transition */
      .animate-start {
        transition: all 1s cubic-bezier(0.25, 0.46, 0.45, 0.94); /* Custom ease for a smooth feel */
        opacity: 0;
      }

      /* Slide from Left Start State */
      .slide-left-start {
        transform: translateX(-100px);
      }
      /* Slide from Right Start State */
      .slide-right-start {
        transform: translateX(100px);
      }

      /* End State (Target) - Triggers the animation */
      .in-view {
        opacity: 1;
        transform: translateX(0) translateY(0);
      }

      /* Custom PButton styling for the demo */
      .p-button {
        padding: 0.75rem 1.5rem;
        font-weight: 600;
        border-radius: 0.5rem;
        color: #ffffff;
        background-color: #d82d63; /* Custom pink 600 */
        border: 1px solid #d82d63;
        transition:
          background-color 0.3s,
          transform 0.2s;
        cursor: pointer;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
      }
      .p-button:hover {
        background-color: #c31b50; /* Slightly darker on hover */
        transform: translateY(-2px);
        box-shadow: 0 6px 12px rgba(216, 45, 99, 0.3);
      }

      /* Utility classes for the dots using the custom pink palette */
      .bg-pink-600-custom {
        background-color: #d82d63;
      }
      .bg-pink-400-custom {
        background-color: #f58ea9;
      }
      .bg-pink-200-custom {
        background-color: #fbcfd0;
      }
    `,
  ],
  // The HTML structure is included directly in the template
  template: `
    <div class="app-root-container">
      <section class="services-showcase py-16 bg-white">
        <!-- Dot indicators using custom colors -->
        <div class="flex justify-center mb-4">
          <span
            class="inline-block w-2 h-2 bg-pink-600-custom rounded-full mx-1"
          ></span>
          <span
            class="inline-block w-2 h-2 bg-pink-400-custom rounded-full mx-1"
          ></span>
          <span
            class="inline-block w-2 h-2 bg-pink-200-custom rounded-full mx-1"
          ></span>
        </div>

        <h1 class="text-3xl font-extrabold text-center mb-12 text-gray-800">
          Celebrate Life's Special Moments
        </h1>

        <div class="container mx-auto px-4 max-w-6xl">
          <!-- Service Block 1: Image Left, Text Right -->
          <div class="mb-20">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
              <!-- Image - Slides from Left -->
              <div
                class="relative animate-start slide-left-start"
                data-animate-target
              >
                <img
                  src="https://bouqs.com/blog/wp-content/uploads/2018/05/graduation-flower-bouquet.jpg"
                  alt="Graduation Bouquets"
                  class="w-full h-80 object-cover rounded-xl shadow-2xl"
                  onerror="this.onerror=null;this.src='https://placehold.co/600x400/F58EA9/ffffff?text=Graduation+Bouquets'"
                />
              </div>

              <!-- Text Content - Slides from Right -->
              <div
                class="space-y-4 animate-start slide-right-start"
                data-animate-target
              >
                <h2 class="text-3xl font-bold text-gray-800">
                  Graduation Bouquets & Gifts
                </h2>
                <p class="text-gray-600 text-lg leading-relaxed">
                  Congratulations on your graduation! Your hard work and
                  dedication have paid off, this is just the beginning of your
                  success. Best wishes for a bright future ahead!
                </p>
                <div class="pt-4">
                  <button (click)="navigateToCategory()" class="p-button">
                    Shop Graduation Flowers
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Service Block 2: Text Left, Image Right (Order Reversed) -->
          <div class="mb-20">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
              <!-- Text Content - Slides from Left (Order 2 on large screens) -->
              <div
                class="space-y-4 order-2 lg:order-1 animate-start slide-left-start"
                data-animate-target
              >
                <h3 class="text-lg text-pink-600-custom font-medium">
                  Floral artistry for your most unforgettable moments
                </h3>
                <h2 class="text-3xl font-bold text-gray-800">
                  Weddings & Events
                </h2>
                <p class="text-gray-600 text-lg leading-relaxed">
                  From intimate elopements to grand celebrations, Floral Haven
                  brings your floral dreams to life. We specialise in bespoke
                  wedding and event arrangements designed to reflect your style,
                  theme, and story.
                </p>
                <div class="pt-4">
                  <button class="p-button">Explore</button>
                </div>
              </div>

              <!-- Image - Slides from Right (Order 1 on large screens) -->
              <div
                class="relative order-1 lg:order-2 animate-start slide-right-start"
                data-animate-target
              >
                <img
                  src="https://images.unsplash.com/photo-1519225421980-715cb0215aed?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80"
                  alt="Wedding and Event Arrangements"
                  class="w-full h-80 object-cover rounded-xl shadow-2xl"
                  onerror="this.onerror=null;this.src='https://placehold.co/600x400/E94F83/ffffff?text=Weddings+and+Events'"
                />
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FeaturedComponent {
  // Lifecycle hook called after Angular initializes the component's views and child views.
  // This is where we safely interact with the DOM to set up the Intersection Observer.
  ngAfterViewInit(): void {
    this.setupIntersectionObserver();
  }

  setupIntersectionObserver(): void {
    // Select all elements tagged for animation
    const targets = document.querySelectorAll('[data-animate-target]');

    // Check for browser support
    if ('IntersectionObserver' in window) {
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              // Add the 'in-view' class to trigger the CSS transition
              entry.target.classList.add('in-view');
              // Once animated, stop observing the element
              observer.unobserve(entry.target);
            }
          });
        },
        {
          // Options
          threshold: 0.1, // Trigger when 10% of the item is visible
          rootMargin: '0px 0px -50px 0px', // Start slightly before the item is fully on screen
        },
      );

      // Start observing each target element
      targets.forEach((target) => {
        observer.observe(target);
      });
    } else {
      // Fallback: If Intersection Observer isn't supported, show content immediately
      targets.forEach((target) => target.classList.add('in-view'));
    }
  }

  #productService = inject(ProductService);
  #router = inject(Router);

  navigateToCategory() {
    this.#productService.categoryId.set(['16']);
    this.#router.navigate(['/products']);
  }
}
