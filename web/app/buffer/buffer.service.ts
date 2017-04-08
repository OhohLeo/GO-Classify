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
    buffersByCollection: Map<string, BufferItem[]> =
		new Map<string, BufferItem[]>()
    buffersById: Map<string, BufferItem> =
		new Map<string, BufferItem>()

    private eventObservers = {}

    constructor(private apiService: ApiService) { }

    // Check if item does exist
    hasBufferItem(search: BufferItem): boolean {
        return this.buffersById.get(search.id) != undefined
    }

    // Add item in all specified collection
    addBufferItem(collection: string, i: BufferItem) {

        // Store buffers by id
        this.buffersById.set(i.id, i)

        // Store buffers by collection
        if (this.buffersByCollection.get(collection) === undefined) {
            this.buffersByCollection.set(collection, [])
        }

        this.buffersByCollection.get(collection).push(i)
    }

    // Delete item from all collection
    deleteBufferItem(collection: string, i: BufferItem) {

        // Delete item by type
        let itemList = this.buffersByCollection.get(collection)
        for (let idx in itemList) {
            let item = itemList[idx]
            if (item.id === i.id) {
                itemList.splice(+idx, 1)
                break;
            }
        }

        this.buffersById.delete(i.id)

		this.enableCache = false
    }

    getBufferItems(collection: string) {

        return new Observable(observer => {

            // Returns the cache if the list should not have changed
            if (this.buffersByCollection && this.enableCache === true) {
                observer.next(this.buffersByCollection.get(collection))
                return
            }

            // Ask for the current list
            this.apiService.get("collections/" + collection + "/buffers")
				.subscribe((rsp: BufferItem[]) => {

					for (let buffer of rsp) {
						this.addBufferItem(collection, buffer)
					}

					this.enableCache = true

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
