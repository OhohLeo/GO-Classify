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
    private items: BufferItem[] = []
    private itemsById: Map<string, number> = new Map<string, number>()

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
        console.log("START")
        this.bufferService.getItems(this.collection)
            .subscribe((buffers) => {
                console.log("REQUEST")
                this.action.modal("open")
            })
    }

    // Check if item is displayed
    hasItem(id: string) {
        return this.itemsById.get(id) != undefined
    }

    add(item: BufferItem) {

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
        console.log("Validate", item)
    }
}
