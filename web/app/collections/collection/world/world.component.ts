import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'
import { Collection } from '../../collection'
import { Item } from '../../item'

declare var jQuery: any

@Component({
    selector: 'world-collection',
    templateUrl: './world.component.html',
})

export class WorldCollectionComponent implements OnInit, OnDestroy {


    @Input() collection: Collection
    @Input() items: Item[]

    constructor() { }

    ngOnInit() {
        jQuery('#map').vectorMap({ map: 'world_mill' });
    }

    ngOnDestroy() { }

}
