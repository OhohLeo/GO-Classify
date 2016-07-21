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
var create_component_1 = require('./create.component');
var modify_component_1 = require('./modify.component');
var delete_component_1 = require('./delete.component');
var CollectionStatus;
(function (CollectionStatus) {
    CollectionStatus[CollectionStatus["NONE"] = 0] = "NONE";
    CollectionStatus[CollectionStatus["CREATE"] = 1] = "CREATE";
    CollectionStatus[CollectionStatus["CHOOSE"] = 2] = "CHOOSE";
    CollectionStatus[CollectionStatus["MODIFY"] = 3] = "MODIFY";
    CollectionStatus[CollectionStatus["DELETE"] = 4] = "DELETE";
})(CollectionStatus || (CollectionStatus = {}));
var CollectionsComponent = (function () {
    function CollectionsComponent(classifyService) {
        var _this = this;
        this.classifyService = classifyService;
        this.collectionStatus = CollectionStatus;
        this.collectionState = CollectionStatus.NONE;
        this.collections = [];
        classifyService.setOnChanges(function (collection) {
            _this.onChooseCollection(collection);
        });
        classifyService.getAll().subscribe(function (list) {
            console.log(list);
            _this.collections = list;
            _this.onChooseCollection(undefined);
        });
    }
    CollectionsComponent.prototype.onNewCollection = function () {
        this.collectionState = CollectionStatus.CREATE;
    };
    CollectionsComponent.prototype.onChooseCollection = function (collection) {
        // If no collection exists : ask to create new one
        if (this.collections.length === 0) {
            this.collectionState = CollectionStatus.CREATE;
            return;
        }
        if (collection !== undefined) {
            this.onSelectCollection(collection);
        }
        this.collectionState = CollectionStatus.CHOOSE;
    };
    CollectionsComponent.prototype.onSelectCollection = function (collection) {
        this.classifyService.selectCollection(collection);
        this.collectionState = CollectionStatus.NONE;
    };
    CollectionsComponent.prototype.onModifyCollection = function (collection) {
        this.onSelectCollection(collection);
        this.collectionState = CollectionStatus.MODIFY;
    };
    CollectionsComponent.prototype.onDeleteCollection = function (collection) {
        this.onSelectCollection(collection);
        this.collectionState = CollectionStatus.DELETE;
    };
    CollectionsComponent.prototype.resetCollectionState = function () {
        // If no collection exists : ask to create new one
        if (this.collections.length === 0) {
            this.collectionState = CollectionStatus.CREATE;
            return;
        }
        if (this.classifyService.collectionSelected == undefined) {
            this.collectionState = CollectionStatus.CHOOSE;
            return;
        }
        this.collectionState = CollectionStatus.NONE;
    };
    __decorate([
        core_1.Input(), 
        __metadata('design:type', String)
    ], CollectionsComponent.prototype, "title", void 0);
    CollectionsComponent = __decorate([
        core_1.Component({
            selector: 'collections',
            templateUrl: 'app/collections/collections.component.html',
            directives: [create_component_1.CreateCollectionComponent,
                modify_component_1.ModifyCollectionComponent,
                delete_component_1.DeleteCollectionComponent]
        }), 
        __metadata('design:paramtypes', [classify_service_1.ClassifyService])
    ], CollectionsComponent);
    return CollectionsComponent;
}());
exports.CollectionsComponent = CollectionsComponent;
//# sourceMappingURL=collections.component.js.map