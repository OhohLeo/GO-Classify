import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService } from './../api.service';
import { Item } from './item'
import { Collection } from './collection'
import { Event } from '../api.service'

export class ItemEvent {

    constructor(public collection: string,
        public status: string,
        public item: Item) { }
}

export class CollectionItem {

    public items: Item[] = []

    isEmpty(): boolean {
        return (this.items.length === 0)
    }

    hasItem(idx: number): boolean {

        if (idx >= 0 && idx < this.items.length)
            return true

        return false
    }

    addItem(item: Item) {
        this.items.push(item)
    }

    removeItem(idx: number) {

        if (this.hasItem(idx)) {
            this.items.splice(idx, 1)
            return
        }

        console.error("Invalid item id", idx)
    }

    getItem(idx: number): Item {

        if (this.hasItem(idx)) {
            return this.items[idx]
        }

        return undefined
    }

    getIdx(id: string): number {

        for (let idx in this.items) {
            if (this.items[idx].id === id) {
                return +idx
            }
        }

        console.error("invalid buffer item with id", id)
        return -1
    }

}

export class ItemObserver {
    public observable: Observable<ItemEvent>
    public observer: any
}

@Injectable()
export class CollectionsService {

    enableCache: boolean
    itemsByCollection: { [collectionName: string]: CollectionItem; } = {}
    itemsById: { [key: string]: Item; } = {}

    private eventObservers: { [key: string]: ItemObserver } = {}

    private onEvent = {
        "add": (collection: string, item: Item) => {
            this.addItem(collection, item)
        },
        "update": (collection: string, item: Item) => {
            this.updateItem(collection, item)
        }
    }

    constructor(private apiService: ApiService) { }

    // Check if item does exist in specified collection
    hasCollectionItem(collection: string, id: string): number {

        let collectionItem = this.getCollection(collection)
        if (collectionItem === undefined)
            return -1

        return collectionItem.getIdx(id)
    }

    getCollection(collection: string): CollectionItem {

        let collectionItem = this.itemsByCollection[collection]
        if (collectionItem === undefined) {
            console.error("collection item '" + collection + "' not found")
            return undefined
        }

        return collectionItem
    }

    // Check if item does exist
    hasItem(search: Item): boolean {
        return this.itemsById[search.id] != undefined
    }

    addItem(collection: string, item: Item) {

        // Do not store already existing item
        if (this.itemsById[item.id] !== undefined)
            return

        // Store items by id
        this.itemsById[item.id] = item

        // Store items by collection
        if (this.itemsByCollection[collection] === undefined) {
            this.itemsByCollection[collection] = new CollectionItem()
        }

        this.itemsByCollection[collection].addItem(item)
    }

    deleteItem(collection: string, item: Item) {

        let collectionItem = this.getCollection(collection)
        if (collectionItem === undefined)
            return

        let itemIdx = collectionItem.getIdx(item.id)
        if (itemIdx < 0)
            return

        collectionItem.removeItem(itemIdx)
        delete this.itemsById[item.id]

        this.enableCache = false
    }

    // Ask for current collection list
    getItems(collection: string) {

        if (this.itemsByCollection[collection] === undefined) {
            this.itemsByCollection[collection] = new CollectionItem()
        }

        return new Observable(observer => {

            let items = this.itemsByCollection[collection]

            // Returns the cache if the list should not have changed
            if (items && this.enableCache === true) {

                observer.next(items.items)
                return
            }

            // Ask for the current list
            this.apiService.get("collections/" + collection + "/items")
                .subscribe((rsp: Item[]) => {

                    for (let item of rsp) {
                        this.addItem(collection, new Item(item))
                    }

                    this.enableCache = true

                    observer.next(items.items)
                })
        })
    }

    updateItem(collection: string, item: Item) {

        if (this.hasItem(item) == false) {
            this.addItem(collection, item)
            return
        }

        let itemIdx = this.hasCollectionItem(collection, item.id)
        if (itemIdx < 0) {
            console.error("Item '" + item.id
                + "' not found in specified collection '"
                + collection + '"', item)
            return
        }

        // Store items by id
        this.itemsById[item.id] = item
        this.itemsByCollection[collection][itemIdx] = item
    }

    subscribeEvents(name: string): Observable<ItemEvent> {

        if (this.eventObservers[name] != undefined) {
            console.error("Already existing observer", name)
            return this.eventObservers[name].observable;
        }

        this.eventObservers[name] = new ItemObserver()

        let observable = Observable.create(observer => {
            this.eventObservers[name].observer = observer
            return () => delete this.eventObservers[name]
        })

        this.eventObservers[name].observable = observable

        return observable
    }


    addEvent(collection: string, event: Event, item: Item) {

        let onEventCb = this.onEvent[event.status]
        if (onEventCb == undefined) {
            console.log("Unhandled event status '" + event.status + "'", event)
            return;
        }

        onEventCb(collection, item)

        let itemEvent = new ItemEvent(collection, event.status, item);

        for (let name in this.eventObservers) {
            this.eventObservers[name].observer.next(itemEvent)
        }
    }
}
