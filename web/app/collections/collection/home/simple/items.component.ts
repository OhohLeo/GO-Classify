import { Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core'
import { Item } from '../../../../items/item'

@Component({
    selector: 'simple-items',
    templateUrl: './items.component.html'
})

export class SimpleItemsComponent implements OnInit, OnDestroy {

    @Input() items: Item[]
    @Output() open: EventEmitter<Item> = new EventEmitter<Item>()

    ngOnInit() { }
    
    ngOnDestroy() { }

    onItem(item: Item) {
	this.open.emit(item)
    }
}
