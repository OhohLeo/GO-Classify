import { Collection } from '../collections/collection'
import { Item } from './item'
import { Ref } from './ref'

export class Items {

    public totalCount: number = 0
    public itemsCount: number = 0

    public refList: Ref[] = []
    public refs: { [ref: string]: Ref } = {}

    public items: { [id: string]: Item } = {}

    constructor(public collection: Collection) {
    }
    
    init(items: any) {

        this.items = {}

        for (let item of items) {
            this.addItem(new Item(this.collection, item))
        }
    }

    isUpToDate(): boolean {
        return this.totalCount == this.itemsCount
    }

    getList(criterious?: any): Item[] {

        let items: Item[] = []

        for (let id in this.items) {
            let item = this.items[id]

            if (criterious != undefined
                && item.match(criterious)) {
                items.push(item)
            }
        }

        return items
    }

    hasItem(item: Item): boolean {
        return (this.items[item.id] != undefined)
    }

    hasItemById(id: string): boolean {
        return (this.items[id]    != undefined)
    }

    addItem(data: any) {
	
        let item = new Item(this.collection, data)

        // Do not store already existing item
        if (this.hasItem(item))
            return
	    
        // Check
        let refType = data["ref"]
        if (refType == undefined) {
            console.error("No ref type specified", data)
            return
        }

        // Create ref if it is unhandled
        let ref = this.refs[refType]
        if (ref == undefined) {
            ref = new Ref(refType, item.data)
            this.refs[refType] = ref
            this.refList.push(ref)
        }

        // Attached ref to item
        item.setRef(ref)

        // Store new item
        this.items[item.id] = item
        this.itemsCount++
    }

    removeItem(item: Item): boolean {

        if (this.hasItem(item)) {
            delete this.items[item.id]
            this.itemsCount--
            return true
        }

        console.error("Invalid item id", item.id)
        return false
    }

    updateItem(item: Item): boolean {

        if (this.hasItem(item)) {

            // TODO: copy event handler?

            this.items[item.id] = item
            return true
        }

        console.error("Invalid item id", item.id)
        return false
    }

    getItem(id: string): Item {
        return this.items[id]
    }
}
