import { Component, NgZone, Input } from '@angular/core'

import { ApiService, CollectionStatus } from '../api.service'
import { CollectionsService } from './collections.service'

import { Collection } from '../collection/collection'

@Component({
    selector: 'collections-list',
    templateUrl: './list.component.html',
})

export class ListCollectionsComponent {
    @Input() title: string

    public collectionStatus = CollectionStatus
    public collectionState = CollectionStatus.NONE
    public collection: Collection
    public collections: Collection[] = []

    constructor(private zone: NgZone,
				private apiService: ApiService,
				private collectionsService: CollectionsService) {
		this.refresh(true)
    }

    refresh(display: boolean) {

        this.collectionsService.getCollections().subscribe(
            (list) => {

                if (list) {
                    this.collections = list
                }

				if (display) {
                    this.onChooseCollection(undefined)
				} else {
					this.zone.run(() => {
						this.collectionState = CollectionStatus.NONE
					})
				}
            })
    }

    nb() : number {
		return this.collections.length
    }

    getCollections(avoid?: Collection): Collection[] {

		let collections: Collection[] = []

		for (let collection of this.collections)
		{
			if (collection != avoid) {
				collections.push(collection)
			}
		}

		return collections
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
        this.zone.run(() => {
			this.collection = collection
            this.collectionState = CollectionStatus.MODIFIED
		})
    }

    onDeleteCollection(collection: Collection) {
		console.log(collection)

		this.zone.run(() => {
			this.collection = collection
			this.collectionState = CollectionStatus.DELETED
		})
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
