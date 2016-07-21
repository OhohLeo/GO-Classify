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
var classify_service_1 = require('./../classify.service');
var ImportDirectoryComponent = (function () {
    function ImportDirectoryComponent(classifyService) {
        this.classifyService = classifyService;
        this.newDirectory = new Directory();
    }
    ImportDirectoryComponent.prototype.onSubmit = function () {
        console.log(this.newDirectory);
    };
    ImportDirectoryComponent = __decorate([
        core_1.Component({
            selector: 'import-directory',
            templateUrl: 'app/imports/directory.component.html'
        }), 
        __metadata('design:paramtypes', [classify_service_1.ClassifyService])
    ], ImportDirectoryComponent);
    return ImportDirectoryComponent;
}());
exports.ImportDirectoryComponent = ImportDirectoryComponent;
var Directory = (function () {
    function Directory() {
    }
    return Directory;
}());
//# sourceMappingURL=directory.component.js.map