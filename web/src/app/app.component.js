"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var core_1 = require('@angular/core');
var classify_service_1 = require('./classify.service');
var collections_component_1 = require('./collections/collections.component');
var directory_component_1 = require('./imports/directory.component');
var AppStatus;
(function (AppStatus) {
    AppStatus[AppStatus["HOME"] = 1] = "HOME";
    AppStatus[AppStatus["IMPORT"] = 2] = "IMPORT";
    AppStatus[AppStatus["EXPORT"] = 3] = "EXPORT";
    AppStatus[AppStatus["CONFIG"] = 4] = "CONFIG";
})(AppStatus || (AppStatus = {}));
var AppComponent = (function () {
    function AppComponent(classifyService) {
        this.classifyService = classifyService;
        this.status = AppStatus;
        this.display = AppStatus.HOME;
        this.title = "Classify";
    }
    AppComponent.prototype.ngOnInit = function () {
        jQuery(".button-collapse").sideNav();
    };
    AppComponent.prototype.onHome = function () {
        this.resetCollectionState();
        this.display = AppStatus.HOME;
    };
    AppComponent.prototype.onImport = function () {
        this.resetCollectionState();
        this.display = AppStatus.IMPORT;
    };
    AppComponent.prototype.onExport = function () {
        this.resetCollectionState();
        this.display = AppStatus.EXPORT;
    };
    AppComponent.prototype.onConfig = function () {
        this.resetCollectionState();
        this.display = AppStatus.CONFIG;
    };
    AppComponent.prototype.onNewCollection = function () {
        this.collections.onNewCollection();
    };
    AppComponent.prototype.onSelectCollection = function () {
        this.display = AppStatus.HOME;
    };
    AppComponent.prototype.onCollectionChoosed = function (collection) {
        this.display = AppStatus.HOME;
        console.log("COLLECTION", collection);
    };
    AppComponent.prototype.resetCollectionState = function () {
        this.collections.resetCollectionState();
    };
    __decorate([
        core_1.ViewChild(collections_component_1.CollectionsComponent), 
        __metadata('design:type', collections_component_1.CollectionsComponent)
    ], AppComponent.prototype, "collections", void 0);
    AppComponent = __decorate([
        core_1.Component({
            selector: 'classify',
            templateUrl: 'app/app.component.html',
            providers: [classify_service_1.ClassifyService],
            directives: [collections_component_1.CollectionsComponent,
                directory_component_1.ImportDirectoryComponent]
        }), 
        __metadata('design:paramtypes', [classify_service_1.ClassifyService])
    ], AppComponent);
    return AppComponent;
}());
exports.AppComponent = AppComponent;
//# sourceMappingURL=app.component.js.map