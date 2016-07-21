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
var collection_1 = require('./collection');
var CreateCollectionComponent = (function () {
    function CreateCollectionComponent(classifySercice) {
        var _this = this;
        this.classifySercice = classifySercice;
        this.collection = new collection_1.Collection('', '');
        classifySercice.getReferences()
            .subscribe(function (references) {
            _this.websites = references["websites"];
            _this.types = references["types"];
            console.log(_this.websites, _this.types);
        });
    }
    CreateCollectionComponent.prototype.onSubmit = function () {
        // Check that a collection with same name doesn't already
        // exist
        var _this = this;
        console.log("CREATE", this.collection);
        // Create new collection
        this.classifySercice.newCollection(this.collection)
            .subscribe(function (status) {
            console.log(status);
            // Reset the collection
            _this.collection = new collection_1.Collection('', '');
        });
    };
    CreateCollectionComponent = __decorate([
        core_1.Component({
            selector: 'collection-create',
            templateUrl: 'app/collections/create.component.html',
        }), 
        __metadata('design:paramtypes', [classify_service_1.ClassifyService])
    ], CreateCollectionComponent);
    return CreateCollectionComponent;
}());
exports.CreateCollectionComponent = CreateCollectionComponent;
//# sourceMappingURL=create.component.js.map