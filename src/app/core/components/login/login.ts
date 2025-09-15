import { Component, effect, inject, signal } from '@angular/core';
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
import { Toast, ToastModule } from 'primeng/toast';
import { PasswordModule } from 'primeng/password';
import { Message } from 'primeng/message';

@Component({
  selector: 'app-login',
  templateUrl: './login.html',
  imports: [
    ButtonModule,
    ReactiveFormsModule,
    InputTextModule,
    PasswordModule,
    ToastModule,
    Message,
  ],
  providers: [MessageService],
})
export class LoginComponent {
  loginForm: FormGroup;
  private authService = inject(AuthService);
  returnUrl: string | null = null;
  loading = signal(false);
  severity = signal('');
  message = signal('');

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required],
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
      this.loading.set(true);
      const { email, password } = this.loginForm.value;
      this.authService.login(email, password).subscribe({
        next: () => {
          this.severity.set('success');
          this.message.set('Success');
          this.loading.set(false);
          if (this.returnUrl) {
            this.router.navigateByUrl(this.returnUrl);
          }
        },
        error: () => {
          this.loading.set(false);
          this.severity.set('error');
          this.message.set('Invalid password or email');
        },
      });
    }
  }
}
