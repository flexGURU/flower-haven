import { Component } from "@angular/core";
import { HeaderComponent } from "./header/header";
import { RouterOutlet } from "@angular/router";
import { FoooterComponent } from "./footer/footer";
@Component({
  selector: "app-client-layout",
  templateUrl: "./client.layout.html",
  imports: [HeaderComponent, RouterOutlet, FoooterComponent],
})
export class ClientLayout {}