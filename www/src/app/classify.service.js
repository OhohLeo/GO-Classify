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
var http_1 = require('@angular/http');
var Rx_1 = require('rxjs/Rx');
var ClassifyService = (function () {
    function ClassifyService(http) {
        this.http = http;
        this.url = "http://localhost:8080/";
    }
    ClassifyService.prototype.setOnChanges = function (changes) {
        this.onChanges = changes;
    };
    ClassifyService.prototype.selectCollection = function (collection) {
        this.collectionSelected = collection;
    };
    ClassifyService.prototype.getOptions = function () {
        return new http_1.RequestOptions({
            headers: new http_1.Headers({ 'Content-Type': 'application/json' })
        });
    };
    // Create a new collection
    ClassifyService.prototype.newCollection = function (collection) {
        var _this = this;
        return this.http.post(this.url + "collections", JSON.stringify(collection), this.getOptions())
            .map(function (res) {
            if (res.status != 204) {
                throw new Error('Impossible to create new collection: ' + res.status);
            }
            // Ajoute la collection nouvellement cr��e
            _this.collections.push(collection);
            _this.onChanges(collection);
        })
            .catch(this.handleError);
    };
    // Modify an existing collection
    ClassifyService.prototype.modifyCollection = function (name, collection) {
        var _this = this;
        return this.http.patch(this.url + "collections/" + name, JSON.stringify(collection), this.getOptions())
            .map(function (res) {
            if (res.status != 204) {
                throw new Error('Impossible to modify collection '
                    + name + ': ' + res.status);
            }
            // Replace the collection from the list
            for (var i in _this.collections) {
                if (_this.collections[i].name === name) {
                    _this.collections[i] = collection;
                    break;
                }
            }
            // Remove the selected collection
            _this.collectionSelected = collection;
            _this.onChanges(collection);
        })
            .catch(this.handleError);
    };
    // Delete an existing collection
    ClassifyService.prototype.deleteCollection = function (name) {
        var _this = this;
        return this.http.delete(this.url + "collections/" + name, this.getOptions())
            .map(function (res) {
            if (res.status != 204) {
                throw new Error('Impossible to modify collection '
                    + name + ': ' + res.status);
            }
            // Remove the collection from the list
            for (var i = 0; i < _this.collections.length; i++) {
                if (_this.collections[i].name === name) {
                    _this.collections.splice(i, 1);
                    break;
                }
            }
            // Reset the selected collection
            _this.collectionSelected = undefined;
            _this.onChanges(undefined);
        })
            .catch(this.handleError);
    };
    // Get the collections list
    ClassifyService.prototype.getAll = function () {
        var _this = this;
        return new Rx_1.Observable(function (observer) {
            if (_this.collections) {
                observer.next(_this.collections);
                return;
            }
            var request = _this.http.get(_this.url + "collections", _this.getOptions())
                .map(_this.extractData)
                .catch(_this.handleError);
            request.subscribe(function (collections) {
                if (collections) {
                    _this.collections = collections;
                    observer.next(collections);
                }
            });
        });
    };
    // Get the collections references
    ClassifyService.prototype.getReferences = function () {
        var _this = this;
        // Setup cache on the references
        return new Rx_1.Observable(function (observer) {
            if (_this.references) {
                observer.next(_this.references);
                return;
            }
            var request = _this.http.get(_this.url + "references")
                .map(_this.extractData)
                .catch(_this.handleError);
            request.subscribe(function (references) {
                _this.references = references;
                observer.next(references);
            });
        });
    };
    ClassifyService.prototype.extractData = function (res) {
        if (res.status < 200 || res.status >= 300) {
            throw new Error('Bad response status: ' + res.status);
        }
        // No content to return
        if (res.status === 204) {
            return true;
        }
        return res.json();
    };
    ClassifyService.prototype.handleError = function (error) {
        var errMsg = error.message || 'Server error';
        console.error(errMsg);
        return Rx_1.Observable.throw(errMsg);
    };
    ClassifyService = __decorate([
        core_1.Injectable(), 
        __metadata('design:paramtypes', [http_1.Http])
    ], ClassifyService);
    return ClassifyService;
}());
exports.ClassifyService = ClassifyService;
//# sourceMappingURL=classify.service.js.map