import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'

import { Item } from '../../item/item'
import { Collection } from '../collection'


@Component({
	selector: 'collection-simple',
	templateUrl: './collection.component.html'
})

export class SimpleCollectionComponent implements OnInit, OnDestroy {

    @Input() collection: Collection
	@Input() items: Item[]

    ngOnInit() {
		console.log(this.items)
	}

	ngOnDestroy() {}
}
