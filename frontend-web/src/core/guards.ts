import { CanActivate, ActivatedRouteSnapshot, Router, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

import { AuthService } from './auth.service';
import { TokenService } from './token.service';

// Revalidates the access token when the user is not authenticated, then
// lets the navigation continue. Fails the route when the refresh fails.
export function authCanActivate(
  tokenService: TokenService,
  auth: AuthService
): CanActivate {
  return {
    canActivate: (child: ActivatedRouteSnapshot, state: RouterStateSnapshot, nav: any) => {
      return new Observable((observer) => {
        const token = auth.getToken();
        if (token) {
          observer.next(nav);
          observer.complete();
          return;
        }
        tokenService
          .refresh()
          .pipe(tap(() => {}))
          .subscribe({
            next: (token) => {
              auth.setToken(token);
              observer.next(nav);
              observer.complete();
            },
            error: () => observer.error(new Error('unauthorized')),
            complete: () => observer.complete(),
          });
      });
    },
  };
}
