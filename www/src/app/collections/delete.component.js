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
var classify_service_1 = require('../classify.service');
var DeleteCollectionComponent = (function () {
    function DeleteCollectionComponent(classifySercice) {
        this.classifySercice = classifySercice;
        // Set the name of the collection to delete
        this.title = classifySercice.collectionSelected.name;
    }
    DeleteCollectionComponent.prototype.onDelete = function () {
        console.log("DELETE", this.title);
        // Delete the collection
        this.classifySercice.deleteCollection(this.title)
            .subscribe(function (status) {
        });
    };
    DeleteCollectionComponent = __decorate([
        core_1.Component({
            selector: 'collection-delete',
            templateUrl: 'app/collections/delete.component.html',
        }), 
        __metadata('design:paramtypes', [classify_service_1.ClassifyService])
    ], DeleteCollectionComponent);
    return DeleteCollectionComponent;
}());
exports.DeleteCollectionComponent = DeleteCollectionComponent;
//# sourceMappingURL=delete.component.js.map