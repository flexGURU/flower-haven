import { Component, inject } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { AuthService } from '../../auth/auth.service';
import { MessageService } from 'primeng/api';
import { Toast } from 'primeng/toast';

@Component({
  selector: 'app-login',
  templateUrl: './login.html',
  imports: [ButtonModule, ReactiveFormsModule, InputTextModule, Toast],
  providers: [MessageService],
})
export class LoginComponent {
  loginForm: FormGroup;
  private authService = inject(AuthService);
  private messageService = inject(MessageService);
  returnUrl: string | null = null;

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.loginForm = this.fb.group({
      email: ['admin1@admin.com', [Validators.required, Validators.email]],
      password: ['admin', Validators.required],
    });
  }

  ngOnInit() {
    this.route.queryParams.subscribe((params) => {
      this.returnUrl = params['returnUrl'] || '/admin';
    });
  }

  get email() {
    return this.loginForm.get('email')!;
  }

  get password() {
    return this.loginForm.get('password')!;
  }

  login() {
    if (this.loginForm.valid) {
      const { email, password } = this.loginForm.value;
      this.authService.login(email, password).subscribe({
        next: () => {
          this.messageService.add({
            severity: 'success',
            summary: 'Login Successful',
          });

          if (this.returnUrl) {
            this.router.navigateByUrl(this.returnUrl);
          }
        },
      });
    }
  }
}
