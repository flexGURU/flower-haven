import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-hero',
  templateUrl: './hero.html',
  imports: [RouterLink, ButtonModule],
})
export class HeroComponent {}
