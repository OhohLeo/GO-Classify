import { Component, NgZone, Input, OnInit, OnDestroy, Renderer } from '@angular/core'
import { CollectionService, ItemEvent } from './collection.service'
import { Event } from '../api.service'
import { Item } from './item'
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
    private items: Item[] = []
    private events: any

    constructor(private zone: NgZone,
        private render: Renderer,
        private collectionService: CollectionService) { }

    ngOnInit() {

        this.getItems()

        // Subscribe to item modification
        this.events = this.collectionService.subscribeEvents(
            this.collection + "/items")
            .subscribe((event: ItemEvent) => {

                // Check if it is the expected collection
                if (event.collection != this.collection.name)
                    return;
            })
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

    // Check if item is displayed
    hasItem(id: string): number {

        for (let idx in this.items) {
            if (id === this.items[idx].id)
                return +idx
        }

        return -1
    }

    add(item: Item) {

        let id = item.id

        // Check if item is already displayed
        if (this.hasItem(id) >= 0) {
            console.error("Add Item with id '" + id + "' already displayed")
            return
        }

        // Add & refresh display
        this.zone.run(() => {
            this.items.push(item)
        })
    }

    remove(item: Item) {

        let id = item.id

        // Check if item does exist
        let idx = this.hasItem(id)
        if (idx < 0) {
            console.error("Remove Item with id '" + id + "' not found")
            return;
        }

        // Delete & refresh display
        this.zone.run(() => {
            this.items.splice(idx, 1)
        })
    }

    update(item: Item) {
        let id = item.id

        // Check if item does exist
        let idx = this.hasItem(id)
        if (idx < 0) {
            console.error("Update item with id '" + id + "' not found")
            return;
        }

        this.zone.run(() => {
            this.items[idx] = item
        })
    }


    getItems() {
        this.collectionService.getItems(this.collection.name)
            .subscribe((items: Item[]) => {
                console.log("UPDATE", items)
                this.zone.run(() => {
                    this.items = items
                })
            })
    }
}
