import { Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core'
import { Items } from '../../../items/items'
import { Item } from '../../../items/item'

@Component({
    selector: 'simple-collection',
    templateUrl: './collection.component.html'
})

export class SimpleCollectionComponent implements OnInit, OnDestroy {

    @Input() items: Items
    @Output() open: EventEmitter<Item> = new EventEmitter<Item>()

    constructor(private zone: NgZone) { }

    ngOnInit() { }

    ngOnDestroy() { }

    onItem(item: Item) {
	this.open.emit(item)
    }
}
