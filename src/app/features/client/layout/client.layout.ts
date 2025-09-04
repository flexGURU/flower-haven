import { Component, HostListener, signal } from '@angular/core';
import { HeaderComponent } from './header/header';
import { RouterOutlet } from '@angular/router';
import { FoooterComponent } from './footer/footer';
@Component({
  selector: 'app-client-layout',
  templateUrl: './client.layout.html',
  imports: [HeaderComponent, RouterOutlet, FoooterComponent],
})
export class ClientLayout {
  isHeaderVisible = signal(true);
  private lastScrollTop = 0;
  private ticking = false;

  @HostListener('window:scroll', [])
  onWindowScroll() {
    if (!this.ticking) {
      window.requestAnimationFrame(() => {
        const scrollTop = window.scrollY || document.documentElement.scrollTop;

        if (scrollTop > this.lastScrollTop + 10) {
          // scrolling down
          this.isHeaderVisible.set(false);
        } else if (scrollTop < this.lastScrollTop - 10) {
          // scrolling up
          this.isHeaderVisible.set(true);
        }

        this.lastScrollTop = scrollTop <= 0 ? 0 : scrollTop;
        this.ticking = false;
      });

      this.ticking = true;
    }
  }
}
