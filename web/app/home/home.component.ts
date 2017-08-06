import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core';
import { HomeService, ItemEvent } from './home.service';
import { Event } from '../api.service';
import { Item } from '../item/item'

declare var jQuery: any;

@Component({
    selector: 'home',
    templateUrl: './home.component.html'
})

export class HomeComponent implements OnInit, OnDestroy {

    @Input() collection: string

    private items: Item[] = []
    private events: any

    constructor(private zone: NgZone,
        private homeService: HomeService) { }

    ngOnInit() {

        this.getItems()

        // Subscribe to item modification
        this.events = this.homeService.subscribeEvents(
            this.collection + "/home")
            .subscribe((event: ItemEvent) => {

                // Check if it is the expected collection
                if (event.collection != this.collection)
                    return;
            })
    }

    ngOnDestroy() {
        if (this.events != undefined) {
            this.events.unsubscribe()
            this.events = undefined
        }
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
        this.homeService.getItems(this.collection)
            .subscribe((items: Item[]) => {
                console.log("UPDATE", items)
                this.zone.run(() => {
                    this.items = items
                })
            })
    }
}
