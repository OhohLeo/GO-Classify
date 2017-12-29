import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { Response } from '@angular/http';
import { BufferItem } from './item'

export class BufferEvent {

    constructor(public collection: string,
        public status: string,
        public buffer: BufferItem) { }
}

export class CollectionBuffer {

    public currentIndex: number = 0
    public items: BufferItem[] = []

    isEmpty(): boolean {
        return (this.items.length === 0)
    }

    hasBufferItem(idx: number): boolean {

        if (idx >= 0 && idx < this.items.length)
            return true

        return false
    }

    addBufferItem(item: BufferItem) {
        this.items.push(item)
    }

    removeBufferItem(idx: number) {

        if (this.hasBufferItem(idx)) {
            this.items.splice(idx, 1)
            return
        }

        console.error("Invalid item id", idx)
    }

    getBufferItem(idx: number): BufferItem {

        if (this.hasBufferItem(idx)) {
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
    buffersById: { [key: string]: BufferItem; } = {}

    private eventObservers: { [key: string]: BufferObserver } = {}

    private onEvent = {
        "create": (collection: string, item: BufferItem) => {
            this.addBufferItem(collection, item)
        },
        "update": (collection: string, item: BufferItem) => {
            this.updateBufferItem(collection, item)
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
    hasBufferItem(search: BufferItem): boolean {
        return this.buffersById[search.id] != undefined
    }

    // Check if item does exist in specified collection
    hasCollectionBufferItem(collection: string, id: string): number {

        let collectionBuffer = this.getCollection(collection)
        if (collectionBuffer === undefined)
            return -1

        return collectionBuffer.getIdx(id)
    }

    getCollectionBufferItem(collection: string, idx: number): BufferItem {

        let collectionBuffer = this.getCollection(collection)
        if (collectionBuffer === undefined)
            return undefined

        return collectionBuffer.getBufferItem(idx)
    }

    disableCache() {
        this.enableCache = false
    }

    // Add buffer item in specified collection
    addBufferItem(collection: string, item: BufferItem) {

        // Do not store already existing item
        if (this.buffersById[item.id] !== undefined)
            return

        // Store buffers by id
        this.buffersById[item.id] = item

        // Store buffers by collection
        if (this.buffersByCollection[collection] === undefined) {
            this.buffersByCollection[collection] = new CollectionBuffer()
        }

        this.buffersByCollection[collection].addBufferItem(item)
    }

    // Remove buffer item from specified collection
    removeBufferItem(collection: string, item: BufferItem) {

        let collectionBuffer = this.getCollection(collection)
        if (collectionBuffer === undefined)
            return

        let itemIdx = collectionBuffer.getIdx(item.id)
        if (itemIdx < 0)
            return

        collectionBuffer.removeBufferItem(itemIdx)
        delete this.buffersById[item.id]

        this.enableCache = false
    }

    getBufferItems(collection: string) {

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
                .subscribe((rsp: BufferItem[]) => {

                    for (let buffer of rsp) {
                        this.addBufferItem(collection, new BufferItem(buffer))
                    }

                    this.enableCache = true

                    observer.next(buffers.items)
                })
        })
    }

    updateBufferItem(collection: string, item: BufferItem) {

        if (this.hasBufferItem(item) == false) {
            this.addBufferItem(collection, item)
            return
        }

        let itemIdx = this.hasCollectionBufferItem(collection, item.id)
        if (itemIdx < 0) {
            console.error("BufferItem '" + item.id
                + "' not found in specified collection '"
                + collection + '"', item)
            return
        }

        // Store buffers by id
        this.buffersById[item.id] = item
        this.buffersByCollection[collection][itemIdx] = item
    }

    // Validate item from specified collection
    validateBufferItem(collection: string, item: BufferItem) {

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
                        this.removeBufferItem(collection, item)
                    }

                    observer.next(ok)
                })
        })

    }

    // Cancel item from specified collection
    cancelBufferItem(collection: string, item: BufferItem) {

        return new Observable<boolean>(observer => {

            // Validate for the current list
            this.apiService.delete(
                "collections/" + collection + "/buffers/" + item.id + "/validate")
                .subscribe((rsp: Response) => {

                    let ok = (rsp.status == 204)
                    if (ok)
                        this.removeBufferItem(collection, item)

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

    addEvent(collection: string, event: Event, item: BufferItem) {

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
