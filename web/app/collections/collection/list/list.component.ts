import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'
import { Items } from '../../../items/items'
import { Item } from '../../../items/item'

@Component({
    selector: 'list-collection',
    templateUrl: './list.component.html',
})

export class ListCollectionComponent implements OnInit, OnDestroy {

    @Input() items: Items

    constructor(private zone: NgZone) { }

    ngOnInit() { }

    ngOnDestroy() { }

}
