import { Injectable } from '@angular/core'
import { Response } from '@angular/http'
import { Observable } from 'rxjs/Rx'
import { ApiService, CollectionStatus } from './../api.service'
import { Items } from '../items/items'
import { Item, ItemObserver, ItemEvent } from '../items/item'
import { Collection } from '../collection/collection'
import { Event } from '../api.service'

@Injectable()
export class CollectionsService {

    private enableCache: boolean
    private collections: Map<string, Collection> = new Map<string, Collection>();

    constructor(private apiService: ApiService) { }

    // Create a new collection
    newCollection(collection: Collection) {

        return this.apiService.post("collections", collection.toApi())
            .map((rsp: Response) => {
                if (rsp.status != 204) {
                    throw new Error('Impossible to create new collection: ' + rsp.status);
                }

                // Ajoute la collection nouvellement créée
                this.collections[collection.name] = collection

                // Choisit automatiquement la nouvelle collection
                this.apiService.setCollection(collection, CollectionStatus.CREATED)
            })
    }

    // Modify an existing collection
    modifyCollection(name: string, collection: Collection) {

        return this.apiService.patch("collections/" + name, collection.toApi())
            .map((rsp: Response) => {
                if (rsp.status != 204) {
                    throw new Error('Impossible to modify collection '
                        + name + ': ' + rsp.status);
                }

                // Replace the collection from the list
                this.collections[name] = collection

                // Choisit automatiquement la nouvelle collection
                this.apiService.setCollection(collection, CollectionStatus.MODIFIED)
            })
    }

    // Delete an existing collection
    deleteCollection(name: string) {

        return this.apiService.delete("collections/" + name)
            .map((rsp: Response) => {
                if (rsp.status != 204) {
                    throw new Error('Impossible to delete collection '
                        + name + ': ' + rsp.status);
                }

                // Remove the collection from the list
                delete this.collections[name]

		console.log("DELETED!")

                // Warn about deleted collection
                this.apiService.setCollection(undefined, CollectionStatus.DELETED)
            })
    }


    // Get the collections list
    getCollections() {

        return new Observable<Collection[]>(observer => {

            if (this.collections.size > 0) {
                observer.next(this.getCollectionsArray())
                return
            }

            let request = this.apiService.get("collections")
                .subscribe(collections => {

                    for (let collection of collections) {
                        this.collections[collection.name] =
                            new Collection(collection.name, collection.ref)
                    }

                    observer.next(this.getCollectionsArray())
                })
        });
    }

    getCollectionsArray(): Collection[] {

        let collections: Collection[] = []

        for (let name in this.collections) {
            collections.push(this.collections[name])
        }

        return collections
    }

    getCollection(name: string): Collection {

        let collection = this.collections[name]
        if (collection === undefined) {
            console.error("collection '" + name + "' not found")
            return undefined
        }

        return collection
    }

    // Ask for current collection list
    getItems(name: string): Observable<Items> {

        return new Observable(observer => {

            let collection = this.getCollection(name)
            if (collection == undefined)
                return

            // Returns the cache if the list should not have changed
            let items = collection.getItems()
            if (collection.enableCache
                && items.isUpToDate()) {
                observer.next(items)
                return
            }

            // Ask for the current list
            this.apiService.get("collections/" + name + "/items")
                .subscribe((rsp: Item[]) => {

		    console.log(rsp)

                    for (let item of rsp) {
                        collection.addItem(item)
                    }

                    collection.enableCache = true
                    observer.next(items)
                })
        })
    }
}
