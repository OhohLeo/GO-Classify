import { Component, Input } from '@angular/core';
import { ApiService, CollectionStatus } from '../api.service';
import { Collection } from '../collections/collection';

@Component({
    selector: 'collections-list',
    templateUrl: './list.component.html',
})

export class ListCollectionsComponent {
    @Input() title: string

    public collectionStatus = CollectionStatus
    public collectionState = CollectionStatus.NONE
    public collections: Collection[] = []

    constructor(private apiService: ApiService) {

        apiService.getCollections().subscribe(
            (list) => {
                this.collections = list
                this.onChooseCollection(undefined)
            })
    }

    onNewCollection() {
        this.collectionState = CollectionStatus.CREATED
    }

    onChooseCollection(collection: Collection): boolean {

        // If no collection exists : ask to create new one
        if (this.collections.length === 0) {
            this.collectionState = CollectionStatus.CREATED
            return true
        }

        if (collection !== undefined) {
            this.onSelectCollection(collection, CollectionStatus.SELECTED)
            this.collectionState = CollectionStatus.NONE
            return true
        }

        // If 1 collection exists : we display this collection
        if (this.collections.length === 1) {
            this.onSelectCollection(this.collections[0], CollectionStatus.SELECTED)
            this.collectionState = CollectionStatus.NONE
            return false
        }

        this.collectionState = CollectionStatus.SELECTED
        return true
    }

    onSelectCollection(collection: Collection, status: CollectionStatus) {
        this.apiService.setCollection(collection, status)
        this.collectionState = CollectionStatus.NONE
    }

    onModifyCollection(collection: Collection) {
        this.onSelectCollection(collection, CollectionStatus.MODIFIED)
        this.collectionState = CollectionStatus.MODIFIED
    }

    onDeleteCollection(collection: Collection) {
        this.onSelectCollection(collection, CollectionStatus.DELETED)
        this.collectionState = CollectionStatus.DELETED
    }

    resetCollectionState() {

        // If no collection exists : ask to create new one
        if (this.collections.length === 0) {
            this.collectionState = CollectionStatus.CREATED
            return
        }

        if (this.apiService.collectionSelected == undefined) {
            this.collectionState = CollectionStatus.SELECTED
            return
        }

        this.collectionState = CollectionStatus.NONE
    }
}
