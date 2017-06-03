import { Component, Input, Output, EventEmitter, NgZone, OnInit, OnDestroy } from '@angular/core'

import { BufferService } from '../buffer/buffer.service'
import { Event } from '../api.service'
import { ConfigService, ConfigBase } from '../config/config.service';
import { ApiService } from '../api.service';
import { Item } from './item'

declare var jQuery: any

@Component({
    selector: 'item',
    templateUrl: './item.component.html',
})

export class ItemComponent implements OnInit, OnDestroy {

    @Input() collection: string
    @Input() item: Item

    @Output() close: EventEmitter<any> = new EventEmitter()
    private needNameDetails: boolean = false
    private bestMatch: any
    private bufferIdx: number

    private selectColor: { [key: string]: string } = {}

    constructor(private zone: NgZone,
        private configService: ConfigService,
        private bufferService: BufferService) {
    }

    ngOnInit() {
        this.onUpdate(this.item)
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
        this.displayItem(false)
    }

    onNext() {
        console.log("next")
        this.displayItem(true)
    }

    displayItem(isNext: boolean) {

        if (this.bufferIdx < 0) {
            console.error("invalid current item idx", this.bufferIdx)
            return
        }

        let collectionBuffer = this.bufferService.getCollection(
            this.collection)
        if (collectionBuffer === undefined)
            return

        let idx = this.bufferIdx + (isNext ? 1 : -1)
        if (idx < 0)
            idx += collectionBuffer.items.length

        let item = collectionBuffer.getItem(
            idx % collectionBuffer.items.length)
        if (item === undefined) {
            console.error("invalid item with idx", idx)
            return
        }

        this.onUpdate(item)
    }

    onValidate() {
        console.log("validate")
    }

    onUpdate(item: Item) {
        this.zone.run(() => {
            this.item = item
            this.bufferIdx = this.bufferService.hasCollectionItem(
                this.collection, item.id)
            console.log("BUFFER IDX", this.bufferIdx, item)
        })

        this.onSelect(item.bestMatchId)
    }

    onChange(event) {
        this.configService.onChange(this.collection, event)
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
        this.item.setBestMatch(id)

        // Reset other selection
        for (let id in this.selectColor) {
            this.selectColor[id] = ""
        }

        this.zone.run(() => {
            this.bestMatch = this.item.getBestMatch()
            this.selectColor[id] = "red"
        })
    }

    // Set the fields depending on the collection

    // Set the import list

    // Set the websites list
}
