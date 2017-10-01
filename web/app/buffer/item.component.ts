import { Component, Input, Output, EventEmitter, NgZone, OnInit, OnDestroy } from '@angular/core'

import { BufferService } from './buffer.service'
import { ApiService, Event } from '../api.service'
import { BufferItem } from './item'

declare var jQuery: any

@Component({
    selector: 'buffer-item',
    templateUrl: './item.component.html',
})

export class BufferItemComponent implements OnInit, OnDestroy {

    @Input() collection: string
    @Input() bufferItem: BufferItem
    @Output() close: EventEmitter<any> = new EventEmitter()

    private needNameDetails: boolean = false
    private match: any
    private bufferIdx: number

    private selectColor: { [key: string]: string } = {}

    constructor(private zone: NgZone,
        private bufferService: BufferService) {
    }

    ngOnInit() {
        this.onUpdate(this.bufferItem)
    }

    ngOnDestroy() {
        this.close.emit()
    }

    getNameDetails() {
        this.zone.run(() => {
            this.needNameDetails = !this.needNameDetails;
        })
    }

    onPrevious() {
        console.log("previous")
        this.displayBufferItem(false)
    }

    onNext() {
        console.log("next")
        this.displayBufferItem(true)
    }

    displayBufferItem(isNext: boolean) {

        if (this.bufferIdx < 0) {
            console.error("invalid current item idx", this.bufferIdx)
            return
        }

        let collectionBuffer = this.bufferService.getCollection(
            this.collection)
        if (collectionBuffer === undefined)
            return

        if (collectionBuffer.isEmpty()) {
            this.close.emit()
            return
        }

        let idx = this.bufferIdx + (isNext ? 1 : -1)
        if (idx < 0)
            idx += collectionBuffer.items.length

        let bufferItem = collectionBuffer.getBufferItem(
            idx % collectionBuffer.items.length)
        if (bufferItem === undefined) {
            console.error("invalid bufferItem with idx", idx)
            return
        }

        this.onUpdate(bufferItem)
    }

    onValidate() {
        this.bufferService.validateBufferItem(this.collection, this.bufferItem)
            .subscribe((ok: boolean) => {
                if (ok) {
                    this.onNext()
                }
            })
    }

    onUpdate(bufferItem: BufferItem) {
        this.zone.run(() => {
            this.bufferItem = bufferItem
            this.bufferIdx = this.bufferService.hasCollectionBufferItem(
                this.collection, bufferItem.id)
        })

        this.onSelect(bufferItem.matchId)
    }

    onChange(event) {
        // this.configService.onChange(this.collection, event)
    }

    onSearch(search: string) {
        console.log("SEARCH!", search)
    }

    onOver(id: string) {
        this.setItemColor(id, "blue")
    }

    onLeave(id: string) {
        this.setItemColor(id, "")
    }

    setItemColor(id: string, color: string) {

        if (this.selectColor[id] === "red")
            return

        this.zone.run(() => {
            this.selectColor[id] = color
        })
    }

    onSelect(id: string) {

        // Select best match item
        this.bufferItem.setMatch(id)

        // Reset other selection
        for (let id in this.selectColor) {
            this.selectColor[id] = ""
        }

        this.zone.run(() => {
            this.match = this.bufferItem.getMatch()
            this.selectColor[id] = "red"
        })
    }

    // Set the fields depending on the collection

    // Set the import list

    // Set the websites list
}
