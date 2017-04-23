import { Component, NgZone, Input, OnInit, OnDestroy, ViewChild } from '@angular/core'
import { Response } from '@angular/http';
import { BufferService, BufferItem, BufferEvent } from './buffer.service'
import { Event, Item } from '../api.service'
import { DetailComponent } from './detail.component'

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

	@ViewChild(DetailComponent) detail: DetailComponent
    private detailBuffer: BufferItem

    constructor(private zone: NgZone,
        private bufferService: BufferService) {
    }

    ngOnInit() {
        this.action = jQuery('div#buffer').modal({
		    complete: () => {

				if (this.events == undefined)
					return

				this.events.unsubscribe()
				this.events = undefined
			}
		})
    }

    ngOnDestroy() {
        this.action.modal("close")
    }

    start() {

		// Get actual buffer items
		this.bufferService.getBufferItems(this.collection)
            .subscribe((buffers: BufferItem[]) => {

				if (buffers.length <= 0)
					return;

				for (let buffer of buffers) {
					this.add(new BufferItem(buffer));
				}

				// If has buffer items : open modal
				this.action.modal("open")
            })

		// Subscribe to buffer modification
		this.events = this.bufferService.subscribeEvents(this.collection + "/buffer")
			.subscribe((event: BufferEvent) => {

				// Check if it is the expected collection
				if (event.collection != this.collection)
					return;

				console.log("GET BUFFER EVENT", event)

				if (event.status === "create")
				{
					this.add(event.buffer)
				}
				else if (event.status === "update")
				{
					this.update(event.buffer)
				}
				else
				{
					console.error("Unhandled buffer event status '" + status + "'")
				}

			})
    }

    // Check if item is displayed
    hasBuffer(id: string) : number {

		for (let idx in this.buffers) {
			if (id === this.buffers[idx].id)
				return +idx
		}

        return -1
    }

    add(buffer: BufferItem) {

        let id = buffer.id

        // Check if buffer is already displayed
        if (this.hasBuffer(id) >= 0) {
            console.error("Add Buffer with id '" + id + "' already displayed")
            return
        }

        // Add & refresh display
        this.zone.run(() => {
            this.buffers.push(buffer)
        })
    }

    remove(buffer: BufferItem) {

        let id = buffer.id

		// Check if buffer does exist
		let idx = this.hasBuffer(id)
        if (idx < 0) {
            console.error("Remove Buffer with id '" + id + "' not found")
            return;
        }

        // Delete & refresh display
        this.zone.run(() => {
            this.buffers.splice(idx, 1)
        })
    }

	update(buffer: BufferItem)
	{
        let id = buffer.id

		// Check if buffer does exist
		let idx = this.hasBuffer(id)
		if (idx < 0) {
			console.error("Update buffer with id '" + id + "' not found")
            return;
		}

		if (this.detailBuffer != undefined
			&& this.detailBuffer.id == id)
		{
			this.detail.onUpdate(buffer)
			this.detailBuffer = buffer
		}

        this.zone.run(() => {
            this.buffers[idx] = buffer
        })
	}

    getDetails(buffer: BufferItem) {
        this.zone.run(() => {
            this.detailBuffer = buffer
        })
    }

    onCloseDetail() {
        this.zone.run(() => {
            this.detailBuffer = undefined
        })
    }

    validate(buffer: BufferItem) {

		this.bufferService.validateBufferItem(this.collection, buffer)
            .subscribe((ok : boolean) => {

				if (ok == false) {
					return
				}

				console.log("Validate", buffer, ok)
				this.remove(buffer)

				// If no more buffer items : close modal
				if (this.buffers.length <= 0)
					this.action.modal("close")
			})
	}
}
