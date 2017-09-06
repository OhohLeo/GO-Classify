import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'

import { Item } from '../item'
import { Collection } from '../collection'


@Component({
	selector: 'simple-collection',
	templateUrl: './collection.component.html'
})

export class SimpleCollectionComponent implements OnInit, OnDestroy {

    @Input() collection: Collection
	@Input() items: Item[]

	refs: string[] = []
	itemsById: Map<string, Item> = new Map<string, Item>()
	itemsByRef: Map<string, Item[]> = new Map<string, Item[]>()

    constructor(private zone: NgZone) {}

    ngOnInit() {

		let needUpdate = false;
		for (let item of this.items) {
			if (this.addItem(item)) {
				needUpdate = true;
			}
		}

		if (needUpdate) {
			this.update()
		}
	}

	ngOnDestroy() {}

	addItem(i: Item) : boolean {

		// Already existing item
		if (this.itemsById.get(i.id) != undefined) {
			return false
		}

		// Store items by id
		this.itemsById[i.id] = i

        // Store items by ref
        if (this.itemsByRef.get(i.getRef()) === undefined) {
            this.itemsByRef.set(i.getRef(), [])
        }

        this.itemsByRef.get(i.getRef()).push(i)
		return true
	}

	removeItem(i: Item) : boolean {

		// Check not existing item
		if (this.itemsById.get(i.id) === undefined) {
			console.error("Can't remove unexisting item")
			return false
		}

		// Delete item by ref
		let items = this.itemsByRef.get(i.getRef())
		for (let idx in items) {
			if (items[idx] === i) {
				items.splice(+idx, 1)
				break
			}
		}

        // Remove refs with no items
        if (items.length == 0) {
            this.itemsByRef.delete(i.getRef())
        }

		// Delete item by Id
		this.itemsById.delete(i.id)

		return true
	}

	update() {
		this.zone.run(() => {
			this.refs = Array.from(this.itemsByRef.keys())
		})
	}
}
