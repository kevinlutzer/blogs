Lazy loaded modules are a really useful angular feature that allow us to do better access control and decrease the initial bundle load time. Let first look at how we might setup a lazy loaded module. To make one, you really only need to do two things:

1. Import a child `RoutingModule` in the module that you want to lazy load. This might look something like: 

``` typescript
ROUTES = [
  {path: "manage-users", component: ManageUsersComponent},
  {path: "", redirectTo: "manage-users"},
]
@NgModule({
  imports: [
    RouterModule.forChild(ROUTES)
  ]
})
export class SuperAdminModule
```

2. Specify what module the children routes are located in, in your parent routing module. For this example lets just say that the parent routing module is the base app routing. 

``` typescript
ROUTES = [
  {path: "some-view", component: SomeViewComponent},
  {path: "superadmin", loadChildren: './path/to/super/admin/module/super-admin.module#SuperAdminModule'}
]
@NgModule({
  imports: [RoutingModule.forRoot(ROUTES)]
})
export class AppRoutingModule{}
```

During compilation the angular cli's `chunk optomizer` will create a new chunk that contains all the module depencies for `SuperAdminModule` defined by it's metadata. So what happens when a user tries to activate the superadmin route? Well, if that route has not been hit yet, the app will do a fetch for the corresponding chunk. 

Lets discuss how can we use lazy loaded modules to do access control. Looking at our example module, do we want every user to have access to the superadmin module code and resources? No, not really! So lets use a routing guard to prevent users for every being able to activate that route within the app and subsequently lazy load the  `SuperAdminModule`. A `canActivate` routing guard is essentially just a typescript class that implements the `CanActivate` interface from angular. Here is the code in the entire file we would need for our simple routing gaurd. 

``` typescript
import {ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot} from '@angular/router';
import {AuthService} from './auth';
import {Observable} from 'rxjs/Observable';

export class SuperAdminRoutingGuard implements CanActivate {
  constructor(private authService: AuthService) {}

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
    return this.authService.isSuperAdmin$
  }
}
```

Note that we can theoretically put any sort of service provider in the constructor parameters. For this example we just want to use a serivce that can get us information about the current user. Looking back on our parent routing module, lets update it to use the routing guard. The finally routes will look like this: 

``` typescript
ROUTES = [
  {path: "some-base-path", component: SomeBaseComponent},
  {
    path: "superadmin", 
    loadChildren: './path/to/super/admin/module/super-admin.module#SuperAdminModule',
    canActivate: [SuperAdminRoutingGuard]
  }
]
```

When A user tries to navigate to the super admin route through a router link or even just manually typing in the route, The angular router will call the `canActivate` method in `SuperAdminRoutingGuard`. If the result is equivalent to `Observable.of(true)`, `ManageUsersComponent` will be surfaced to the user. If not, the view will not change.