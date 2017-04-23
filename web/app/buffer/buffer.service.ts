import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event, Item } from './../api.service';
import { CfgStringList } from '../config/stringlist.component'
import { Response } from '@angular/http';

export class BufferItem {

    public id: string
    public probability: number
    public name: string
    public image: string

	public cleanedName: string
	public banned: CfgStringList
	public separators: CfgStringList

    private bestMatch: any
    private imports: any[] = []

	public webQuery: string
    private websites: any[] = []

    constructor(public data: any) {
        this.id = data.id
        this.name = (data.name != undefined) ? data.name : "<unknown>"

		this.cleanedName = data.cleanedName
		this.banned = new CfgStringList(data.banned)
		this.separators = new CfgStringList(data.separators)

        this.bestMatch = data.bestMatch
        this.probability = data.probability

        if (data.imports != undefined) {
            for (let key in data.imports) {
                data.imports[key].forEach((data: any) => {
                    this.imports.push(data)
                })
            }
        }

		this.webQuery = data.webQuery
        if (data.websites != undefined) {
            for (let key in data.websites) {
                data.websites[key].forEach((data: any) => {
                    this.websites.push(data)
                })
            }
        }
    }

	public getName() : string {
		return (this.cleanedName === "") ? this.name : this.cleanedName;
	}

	public getImports() : any[] {
		return this.imports
	}

	public getWebsites() : any[] {
		return this.websites
	}
}

export class BufferEvent {

	constructor(public collection: string,
				public status: string,
				public buffer: BufferItem) {}
}

@Injectable()
export class BufferService {

    enableCache: boolean
    buffersByCollection: Map<string, BufferItem[]> =
		new Map<string, BufferItem[]>()
    buffersById: Map<string, BufferItem> =
		new Map<string, BufferItem>()

    private eventObservers = {}

	private onEvent = {
		"create": (collection: string, item: BufferItem) => {
			this.addBufferItem(collection, item)
		},
		"update": (collection: string, item: BufferItem) => {
			this.updateBufferItem(collection, item)
		},
	}

    constructor(private apiService: ApiService) { }

    // Check if item does exist
    hasBufferItem(search: BufferItem): boolean {
        return this.buffersById.get(search.id) != undefined
    }

	// Check if item does exist in specified collection
    hasCollectionBufferItem(collection: string, search: BufferItem): number {

		let itemList = this.buffersByCollection.get(collection)

        for (let idx in itemList) {
            if (itemList[idx].id === search.id) {
                return +idx
            }
        }

		return -1
    }

	disableCache() {
		this.enableCache = false
	}

    // Add buffer item in specified collection
    addBufferItem(collection: string, i: BufferItem) {

        // Store buffers by id
        this.buffersById.set(i.id, i)

        // Store buffers by collection
        if (this.buffersByCollection.get(collection) === undefined) {
            this.buffersByCollection.set(collection, [])
        }

        this.buffersByCollection.get(collection).push(i)

		console.log(this.buffersByCollection.get(collection))
    }

    // Delete buffer item from specified collection
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

            // // Returns the cache if the list should not have changed
            if (this.buffersByCollection && this.enableCache === true) {

				let buffers = this.buffersByCollection.get(collection)
				if (buffers === undefined) {
					buffers = [];
				}

                observer.next(buffers)
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

	updateBufferItem(collection: string, item: BufferItem) {

		if (this.hasBufferItem(item) == false)
		{
			console.error("Item '" + item.id + "' not found", item)
			return;
		}

		let idx = this.hasCollectionBufferItem(collection, item)
		if (idx < 0)
		{
			console.error("Item '" + item.id + "' not found in specified collection '"
						  + collection +'"', item)
			return;
		}

		// Store buffers by id
        this.buffersById.set(item.id, item)
        this.buffersByCollection.get(collection)[idx] = item;

		console.log("UPDATE ", item);
	}

    // Validate item from specified collection
    validateBufferItem(collection: string, item: BufferItem) {

		return new Observable<boolean>(observer => {

            // Validate for the current list
            this.apiService.put(
				"collections/" + collection + "/buffers/" + item.id + "/validate")
				.subscribe((rsp : Response) => {

					let ok = (rsp.status == 204)
					if (ok)
					{
						this.deleteBufferItem(collection, item)
						this.getBufferItems(collection)
					}

					console.log("VALIDATE ", rsp)
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
				.subscribe((rsp : Response) => {

					let ok = (rsp.status == 204)
					if (ok)
						this.deleteBufferItem(collection, item)

					console.log("DELETE ", rsp)
					observer.next(ok)
				})
        })
	}

    subscribeEvents(name: string): Observable<BufferEvent> {

        if (this.eventObservers[name] != undefined) {
            console.error("Already existing observer", name)
            return;
        }

        return Observable.create(observer => {
            this.eventObservers[name] = observer
			return () => {
				delete this.eventObservers[name]
			}
        })
    }

    addEvent(collection:string, event: Event) {

		let onEventCb = this.onEvent[event.status]
		if (onEventCb == undefined)
		{
			console.log("Unhandled event status '" + event.status + "'", event)
			return;
		}

		let buffer = new BufferItem(event.data)

		onEventCb(collection, buffer)

		let bufferEvent = new BufferEvent(collection, event.status, buffer);

        for (let name in this.eventObservers) {
            this.eventObservers[name].next(bufferEvent)
        }
    }
}
