import { Collection } from '../collection/collection'
import { Item, ItemEvent } from './item'

export class CollectionItemEvent extends ItemEvent {

    constructor(public collection: string,
				public status: string,
				public item: Item) {
		super(status, item)
	}
}

export class CollectionItem extends Item {

	constructor(public collection: Collection, item: any) {
		super(item)
	}


    getPath() : string {
		return "collections/" + this.collection.name + "/items/" + this.id
    }

	delete() {
		this.collection.deleteItem(this)
    }
}
