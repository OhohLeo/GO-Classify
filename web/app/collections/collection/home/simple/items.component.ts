import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'

import { Item } from '../../../item'
import { Collection } from '../../../collection'


@Component({
    selector: 'simple-items',
    templateUrl: './items.component.html'
})

export class SimpleItemsComponent implements OnInit, OnDestroy {
    @Input() items: Item[]

    ngOnInit() { }
    ngOnDestroy() { }
}
