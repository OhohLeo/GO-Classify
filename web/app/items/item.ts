import { Observable } from 'rxjs/Rx'
import { Collection } from '../collections/collection'
import { Ref } from './ref'

export class ItemEvent {

    constructor(public collection: string,
        public status: string,
        public item: Item) { }
}

export class ItemObserver {
    public observable: Observable<ItemEvent>
    public observer: any
}

export class Item {

    public id: string
    public data: any
    public ref: Ref

    public name: string

    public hasIcon: boolean = false
    public icons: { [key: string]: string } = {}

    private eventObservers: { [key: string]: ItemObserver } = {}

    constructor(public collection: Collection, item: any) {

        this.id = item["id"]
        this.data = item["data"]
        this.name = this.data["name"]

	let icons = this.data["icons"]
	if (icons !== undefined) {
	    for (let size in icons) {  
		this.icons[size] = icons[size]["Name"]
	    }
	    this.hasIcon = true
        }
    }

    getIconUrl(size?: string): string {

	// If size does exist : take it
        if (size == undefined
	    || this.icons[size] == undefined) {

	    // Otherwise take 1st size found
            for (let size in this.icons) {
		return size
            }
	}

	return size
    }

    setRef(ref: Ref) {
        this.ref = ref
    }

    getRef() {
        return this.ref.name
    }

    match(criterious: any): boolean {

        if (criterious['hasAttribute'] !== undefined
            && this.ref.hasAttribute(criterious['hasAttribute']) == false) {
            return false
        }

        if (criterious['ref'] !== undefined
            && criterious['ref'] !== this.ref.name)
            return false

        return true
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
}
