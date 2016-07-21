import {bootstrap}    from '@angular/platform-browser-dynamic';
import {HTTP_PROVIDERS} from '@angular/http';
import {AppComponent} from './app/app.component';


// Add all operators to Observable
import 'rxjs/Rx';

/*if (process.env.ENV === 'production') {
  enableProdMode();
}*/

bootstrap(AppComponent, [HTTP_PROVIDERS]);
