import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'
import { Items } from '../../../../items/items'
import { Item } from '../../../../items/item'

@Component({
    selector: 'simple-collection',
    templateUrl: './collection.component.html'
})

export class SimpleCollectionComponent implements OnInit, OnDestroy {

    @Input() items: Items

    constructor(private zone: NgZone) { }

    ngOnInit() { }

    ngOnDestroy() { }
}
