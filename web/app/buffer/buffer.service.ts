import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { Response } from '@angular/http';
import { Item } from '../item/item'

export class BufferEvent {

    constructor(public collection: string,
        public status: string,
        public buffer: Item) { }
}

export class CollectionBuffer {

    public currentIndex: number = 0
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

export class BufferObserver {
    public observable: Observable<BufferEvent>
    public observer: any
}

@Injectable()
export class BufferService {

    enableCache: boolean

    buffersByCollection: { [collectionName: string]: CollectionBuffer; } = {}
    buffersById: { [key: string]: Item; } = {}

    private eventObservers: { [key: string]: BufferObserver } = {}

    private onEvent = {
        "create": (collection: string, item: Item) => {
            this.addItem(collection, item)
        },
        "update": (collection: string, item: Item) => {
            this.updateItem(collection, item)
        }
    }

    constructor(private apiService: ApiService) { }

    getCollection(collection: string): CollectionBuffer {

        let collectionBuffer = this.buffersByCollection[collection]
        if (collectionBuffer === undefined) {
            console.error("collection buffer '" + collection + "' not found")
            return undefined
        }

        return collectionBuffer
    }

    // Check if item does exist
    hasItem(search: Item): boolean {
        return this.buffersById[search.id] != undefined
    }

    // Check if item does exist in specified collection
    hasCollectionItem(collection: string, id: string): number {

        let collectionBuffer = this.getCollection(collection)
        if (collectionBuffer === undefined)
            return -1

        return collectionBuffer.getIdx(id)
    }

    getCollectionItem(collection: string, idx: number): Item {

        let collectionBuffer = this.getCollection(collection)
        if (collectionBuffer === undefined)
            return undefined

        return collectionBuffer.getItem(idx)
    }

    disableCache() {
        this.enableCache = false
    }

    // Add buffer item in specified collection
    addItem(collection: string, item: Item) {

        // Do not store already existing item
        if (this.buffersById[item.id] !== undefined)
            return

        // Store buffers by id
        this.buffersById[item.id] = item

        // Store buffers by collection
        if (this.buffersByCollection[collection] === undefined) {
            this.buffersByCollection[collection] = new CollectionBuffer()
        }

        this.buffersByCollection[collection].addItem(item)
    }

    // Remove buffer item from specified collection
    removeItem(collection: string, item: Item) {

        let collectionBuffer = this.getCollection(collection)
        if (collectionBuffer === undefined)
            return

        let itemIdx = collectionBuffer.getIdx(item.id)
        if (itemIdx < 0)
            return

        collectionBuffer.removeItem(itemIdx)
        delete this.buffersById[item.id]

        this.enableCache = false
    }

    getItems(collection: string) {

        if (this.buffersByCollection[collection] === undefined) {
            this.buffersByCollection[collection] = new CollectionBuffer()
        }

        return new Observable(observer => {

            let buffers = this.buffersByCollection[collection]

            // Returns the cache if the list should not have changed
            if (buffers && this.enableCache === true) {

                observer.next(buffers.items)
                return
            }

            // Ask for the current list
            this.apiService.get("collections/" + collection + "/buffers")
                .subscribe((rsp: Item[]) => {

                    for (let buffer of rsp) {
                        this.addItem(collection, new Item(buffer))
                    }

                    this.enableCache = true

                    observer.next(buffers.items)
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

        // Store buffers by id
        this.buffersById[item.id] = item
        this.buffersByCollection[collection][itemIdx] = item
    }

    // Validate item from specified collection
    validateItem(collection: string, item: Item) {

        return new Observable<boolean>(observer => {

            let match = item.getMatch()

            // Validate for the current list
            this.apiService.post(
                "collections/" + collection + "/buffers/" + item.id + "/validate",
                match)
                .subscribe((rsp: Response) => {

                    let ok = (rsp.status == 204)
                    if (ok) {
                        // Remove item from buffer
                        this.removeItem(collection, item)
                    }

                    observer.next(ok)
                })
        })

    }

    // Cancel item from specified collection
    cancelItem(collection: string, item: Item) {

        return new Observable<boolean>(observer => {

            // Validate for the current list
            this.apiService.delete(
                "collections/" + collection + "/buffers/" + item.id + "/validate")
                .subscribe((rsp: Response) => {

                    let ok = (rsp.status == 204)
                    if (ok)
                        this.removeItem(collection, item)

                    observer.next(ok)
                })
        })
    }

    subscribeEvents(name: string): Observable<BufferEvent> {

        if (this.eventObservers[name] != undefined) {
            console.error("Already existing observer", name)
            return this.eventObservers[name].observable;
        }

        this.eventObservers[name] = new BufferObserver()

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

        let bufferEvent = new BufferEvent(collection, event.status, item);

        for (let name in this.eventObservers) {
            this.eventObservers[name].observer.next(bufferEvent)
        }
    }
}
