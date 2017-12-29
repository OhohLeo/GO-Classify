import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'

import { Item } from '../../../../items/item'


@Component({
    selector: 'simple-items',
    templateUrl: './items.component.html'
})

export class SimpleItemsComponent implements OnInit, OnDestroy {

    @Input() items: Item[]

    ngOnInit() {
	console.log("SIMPLE ITEMS", this.items)
    }
    
    ngOnDestroy() { }
}
