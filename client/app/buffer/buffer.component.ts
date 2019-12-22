import {
    Component, NgZone, Input, Output,
    EventEmitter, OnInit, OnDestroy
} from '@angular/core'

import { Response } from '@angular/http';
import { BufferService, BufferEvent } from './buffer.service'
import { Event } from '../api.service'
import { BufferItemComponent } from './item.component'
import { BufferItem } from './item'

declare var jQuery: any

@Component({
    selector: 'buffer',
    templateUrl: './buffer.component.html',
})

export class BufferComponent implements OnInit, OnDestroy {

    @Input() collection: string
    @Output() onBufferItemSelected = new EventEmitter();

    private action: any
    private events: any
    private buffers: BufferItem[] = []

    constructor(private zone: NgZone,
        private bufferService: BufferService) {
    }

    ngOnInit() {
        this.action = jQuery('div#buffer').modal({
            complete: () => {
                this.stop();
            }
        })
    }

    ngOnDestroy() {
        this.stop()
    }

    start() {

        // Get actual buffer items
        this.bufferService.getBufferItems(this.collection)
            .subscribe((buffers: BufferItem[]) => {

                if (buffers.length <= 0)
                    return

                this.zone.run(() => {
                    this.buffers = buffers
                })

                // Subscribe to buffer modification
                this.events = this.bufferService.subscribeEvents(
                    this.collection + "/buffer")
                    .subscribe((event: BufferEvent) => {

                        // Check if it is the expected collection
                        if (event.collection != this.collection)
                            return;

                        if (event.status === "create") {
                            this.add(event.buffer)
                        }
                        else if (event.status === "update") {
                            this.update(event.buffer)
                        }
                        else {
                            console.error("Unhandled buffer event status '"
                                + status + "'")
                        }

                    })

                // If has buffer items : open modal
                this.action.modal("open")
            })
    }


    stop() {
        this.events.unsubscribe()
        this.events = undefined
    }

    // Check if item is displayed
    hasBuffer(id: string): number {

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

    update(buffer: BufferItem) {
        let id = buffer.id

        // Check if buffer does exist
        let idx = this.hasBuffer(id)
        if (idx < 0) {
            console.error("Update buffer with id '" + id + "' not found")
            return;
        }

        this.zone.run(() => {
            this.buffers[idx] = buffer
        })
    }

    onBufferItem(item: BufferItem) {
        this.onBufferItemSelected.emit(item)
        this.action.modal('close')
    }

    validate(item: BufferItem) {

        this.bufferService.validateBufferItem(this.collection, item)
            .subscribe((ok: boolean) => {

                if (ok == false) {
                    return
                }

                this.remove(item)

                // If no more buffer items : close modal
                if (this.buffers.length <= 0)
                    this.action.modal("close")
            })
    }
}
