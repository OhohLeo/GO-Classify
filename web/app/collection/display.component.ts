import { Component, NgZone, Input, OnInit, OnDestroy, Renderer } from '@angular/core'
import { CollectionsService } from '../collections/collections.service'
import { Event } from '../api.service'
import { Items } from '../items/items'
import { Item, ItemEvent } from '../items/item'
import { Collection } from './collection'

enum ModeStatus {
    HOME = 0,
    LIST,
    WORLD,
    HISTORY,
}

@Component({
    selector: 'collection',
    templateUrl: './display.component.html'
})

export class DisplayCollectionComponent implements OnInit, OnDestroy {

    @Input() collection: Collection

    public modeStatus = ModeStatus
    private modes: string[] = ["star", "list", "language", "history"]
    private currentMode: ModeStatus
    private items: Items
    private item: Item
    private events: any

    constructor(private zone: NgZone,
		private render: Renderer,
		private collectionsService: CollectionsService) { }

    ngOnInit() {

        this.items = this.collection.items

        // Update items list
        this.initItems()

        // // Subscribe to item modification
        // this.events = this.collectionsService.subscribeEvents(
        //     this.collection + "/items")
        //     .subscribe((event: ItemEvent) => {

        //         // Check if it is the expected collection
        //         if (event.collection != this.collection.name)
        //             return;
        //     })
    }

    ngOnDestroy() {
        if (this.events != undefined) {
            this.events.unsubscribe()
            this.events = undefined
        }
    }


    onMode(event: any, mode: ModeStatus) {

        // Set collection-items as active
        event.preventDefault()

        let target
        if (event.target.tagName === "I") {
            target = event.target.parentElement
        } else {
            target = event.target
        }

        for (let item of target.parentElement.children) {
            this.render.setElementClass(item, "active", false)
        }

        this.render.setElementClass(target, "active", true)

        this.zone.run(() => {
            this.currentMode = mode
        })
    }

    onItem(item: Item) {
	this.zone.run(() => {
            this.item = item
        })
    }

    onCloseItem() {
	this.zone.run(() => {
            this.item = undefined
        })
    }

    initItems() {
        this.collectionsService.getItems(this.collection.name)
            .subscribe((items: Items) => {
                this.zone.run(() => {
                    console.log("UPDATE", items)
                })
            })
    }
}
