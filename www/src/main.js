"use strict";
var platform_browser_dynamic_1 = require('@angular/platform-browser-dynamic');
var http_1 = require('@angular/http');
var app_component_1 = require('./app/app.component');
// Add all operators to Observable
require('rxjs/Rx');
/*if (process.env.ENV === 'production') {
  enableProdMode();
}*/
platform_browser_dynamic_1.bootstrap(app_component_1.AppComponent, [http_1.HTTP_PROVIDERS]);
//# sourceMappingURL=main.js.map