import { Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core'
import { Items } from '../../../items/items'
import { Item } from '../../../items/item'

@Component({
    selector: 'list-collection',
    templateUrl: './list.component.html',
})

export class ListCollectionComponent implements OnInit, OnDestroy {

    @Input() items: Items
    @Output() open: EventEmitter<Item> = new EventEmitter<Item>()

    constructor(private zone: NgZone) { }

    ngOnInit() { }

    ngOnDestroy() { }

    onOpen(item: Item) {
	this.open.emit(item)
    }
}
