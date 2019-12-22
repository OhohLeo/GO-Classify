import { Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core'
import { Collection } from '../collection'
import { Items } from '../../items/items'
import { Item } from '../../items/item'

declare var jQuery: any

@Component({
    selector: 'world-collection',
    templateUrl: './world.component.html',
})

export class WorldCollectionComponent implements OnInit, OnDestroy {


    @Input() collection: Collection
    @Input() items: Items
    @Output() open: EventEmitter<Item> = new EventEmitter<Item>()

    constructor() { }

    ngOnInit() {
        jQuery('#map').vectorMap({ map: 'world_mill' });
    }

    ngOnDestroy() { }

    onItem(item: Item) {
	this.open.emit(item)
    }
}
