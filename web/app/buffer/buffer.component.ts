import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'

import { BufferService, BufferItem } from './buffer.service'
import { Event, Item } from '../api.service'

declare var jQuery: any

@Component({
    selector: 'buffer',
    templateUrl: './buffer.component.html',
})

export class BufferComponent implements OnInit, OnDestroy {

    @Input() collection: string

    private action: any
    private events: any
    private buffers: BufferItem[] = []

    private detailItem: BufferItem

    constructor(private zone: NgZone,
        private bufferService: BufferService) {

        this.events = bufferService.subscribeEvents("status")
            .subscribe((e: Event) => {
                this.add(new BufferItem(e.event, e.data))
            })
    }

    ngOnInit() {
        this.action = jQuery('div#buffer').modal()
    }

    ngOnDestroy() {
        this.action.modal("close")
    }

    start() {
		this.bufferService.getBufferItems(this.collection)
            .subscribe((buffers: BufferItem[]) => {

				for (let buffer of buffers) {
					this.add(buffer);
				}

				// If has buffer items : open modal
				if (buffers.length > 0)
					this.action.modal("open")
            })
    }

    // Check if item is displayed
    hasItem(id: string) : number {

		for (let idx in this.buffers) {
			if (id === this.buffers[idx].id)
				return +idx
		}

        return -1
    }

    add(item: BufferItem) {

        let id = item.id

        // Check if item is already displayed
        if (this.hasItem(id) >= 0) {
            console.error("Item with id '" + id + "' already displayed")
            return
        }

        // Add & refresh display
        this.zone.run(() => {
            this.buffers.push(item)
        })
    }

    remove(item: BufferItem) {

        let id = item.id

		// Check if item does exist
		let idx = this.hasItem(id)
        if (idx < 0) {
            console.error("Item with id '" + id + "' not found")
            return;
        }

        // Delete & refresh display
        this.zone.run(() => {
            this.buffers.splice(idx, 1)
        })
    }

    getDetails(item: BufferItem) {
        this.zone.run(() => {
            this.detailItem = item
        })
    }

    onCloseDetail() {
        this.zone.run(() => {
            this.detailItem = undefined
        })
    }

    validate(item: BufferItem) {
		this.remove(item)
        console.log("Validate", item)

		// If no more buffer items : close modal
		if (this.buffers.length <= 0)
			this.action.modal("close")
    }
}
