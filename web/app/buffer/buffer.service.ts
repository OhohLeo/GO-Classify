import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event, Item } from './../api.service';
import { Response } from '@angular/http';

export class BufferItem {

    public id: string
    public probability: number
    public name: string
    public image: string

    private bestMatch: any
    private imports: any[] = []
    private websites: any[] = []

    constructor(public type: string, public data: any) {
        this.id = data.id
        this.name = (data.name != undefined) ? data.name : "<unknown>"
        this.bestMatch = data.bestMatch
        this.probability = data.probability

        if (data.imports != undefined) {
            for (let key in data.imports) {
                data.imports[key].forEach((data: any) => {
                    this.imports.push(data)
                })
            }
        }

        if (data.websites != undefined) {
            for (let key in data.websites) {
                data.websites[key].forEach((data: any) => {
                    this.websites.push(data)
                })
            }
        }
    }
}

@Injectable()
export class BufferService {

    enableCache: boolean
    itemsByCollection: Map<string, Item[]> = new Map<string, Item[]>()
    itemsById: Map<string, Item> = new Map<string, Item>()

    private eventObservers = {}

    constructor(private apiService: ApiService) { }

    // Check if item does exist
    hasItem(search: Item): boolean {
        return this.itemsById.get(search.id) != undefined
    }

    // Add item in all specified collection
    add(i: Item) {

        // Store items by id
        this.itemsById.set(i.id, i)

        // Store items by collection
        for (let collection of i.collections) {
            if (this.itemsByCollection.get(collection) === undefined) {
                this.itemsByCollection.set(collection, [])
            }

            this.itemsByCollection.get(collection).push(i)
        }
    }

    // Delete item from all collection
    deleteItem(i: Item) {
        for (let collection of i.collections) {
            this.deleteItemFromCollection(i, collection)
        }
    }

    // Delete item from specified collection
    deleteItemFromCollection(i: Item, collection: string) {

        // Detach item & collection
        for (let idx in i.collections) {
            if (i.collections[idx] == collection) {
                i.collections.splice(+idx, 1);
            }
        }

        // Delete item by type
        let itemList = this.itemsByCollection.get(collection)
        for (let idx in itemList) {
            let item = itemList[idx]
            if (item.id === i.id) {
                itemList.splice(+idx, 1)
                break;
            }
        }

        // No more collection are referenced by the item
        if (i.collections.length == 0) {
            // Delete item by id
            this.itemsById.delete(i.id)
        }
    }

    getItems(collection: string) {

        return new Observable(observer => {

            // Returns the cache if the list should not have changed
            // if (this.itemsByCollection && this.enableCache === true) {
            //     observer.next(this.itemsByCollection.get(collection))
            //     return
            // }

            // Ask for the current list
            this.apiService.get("collection/" + collection + "/buffers").subscribe(rsp => {

                // // Init the item lists
                // this.items = new Map<string, Item[]>()
                // this.itemsById = new Map<string, Item>()

                // for (let itemType in rsp) {

                //     for (let itemId in rsp[itemType]) {
                //         let i = convert(itemId, rsp[itemType][itemId])
                //         if (i === undefined)
                //             continue

                //         this.add(i)
                //     }
                // }

                // this.enableCache = true

                console.log(rsp)

                observer.next(rsp)
            })
        })
    }


    subscribeEvents(name: string): Observable<Event> {

        if (this.eventObservers[name] != undefined) {
            console.error("Already existing observer", name)
            return;
        }

        return Observable.create(observer => {

            // Initialisation de l'observer
            this.eventObservers[name] = observer

            return () => delete this.eventObservers[name]
        })
    }

    addEvent(event: Event) {
        for (let name in this.eventObservers) {
            this.eventObservers[name].next(event)
        }
    }
}
