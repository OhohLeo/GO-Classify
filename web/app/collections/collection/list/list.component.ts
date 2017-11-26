import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'
import { Collection } from '../../collection'
import { Item } from '../../item'

@Component({
    selector: 'list-collection',
    templateUrl: './list.component.html',
})

export class ListCollectionComponent implements OnInit, OnDestroy {


    @Input() collection: Collection
    @Input() items: Item[]

    constructor() { }

    ngOnInit() {
        console.log("HERE!")
    }

    ngOnDestroy() { }

}
