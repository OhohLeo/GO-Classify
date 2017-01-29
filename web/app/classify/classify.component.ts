import { Component, NgZone, OnInit, OnDestroy } from '@angular/core'

import { ClassifyService, ClassifyItem } from './classify.service'
import { Event, Item } from '../api.service'

declare var jQuery: any

@Component({
    selector: 'classify',
    templateUrl: './classify.component.html',
})

export class ClassifyComponent implements OnInit, OnDestroy {

	private action: any
    private events: any
	private items: ClassifyItem[] = []
	private itemsById: Map<string, number> = new Map<string, number>()

	private detailItem: ClassifyItem

    constructor(private zone: NgZone,
				private classifyService: ClassifyService) {

        this.events = classifyService.subscribeEvents("status")
            .subscribe((e: Event) => {
				this.add(new ClassifyItem(e.event, e.data))
			})
	}

    ngOnInit() {
		this.action = jQuery('div#classify').modal()
    }

    ngOnDestroy() {
        this.action.modal("close")
    }

    start() {
        this.action.modal("open")
    }

	// Check if item is displayed
	hasItem(id: string) {
		return this.itemsById.get(id) != undefined
	}

	add(item: ClassifyItem) {

		let id = item.id

		// Check if item is already displayed
		if (this.hasItem(id)) {
			console.error("Item with id '" + id + "' already displayed")
			return;
		}

		// Add & refresh display
		this.zone.run(() => {
			this.items.push(item)
			this.itemsById.set(id, this.items.length - 1)
		})
	}

	remove(id: string) {

		// Check if item does exist
		if (this.hasItem(id) == false) {
			console.error("Item with id '" + id + "' not found")
			return;
		}

		// Delete & refresh display
		this.zone.run(() => {
			this.items.splice(this.itemsById.get(id), 1)
			this.itemsById.delete(id)
		})
	}

    getDetails(item: ClassifyItem) {
		this.zone.run(() => {
			this.detailItem = item
		})
    }

	onCloseDetail()	{
		this.zone.run(() => {
			this.detailItem = undefined
		})
	}

	validate(item: ClassifyItem) {
		console.log("Validate", item)
	}
}
