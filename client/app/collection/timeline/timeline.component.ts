import { Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core'
import { Items } from '../../items/items'
import { Item } from '../../items/item'
import * as Vis from 'vis'

@Component({
    selector: 'timeline-collection',
    templateUrl: './timeline.component.html',
})

export class TimeLineCollectionComponent implements OnInit, OnDestroy {

    @Input() items: Items
    @Output() open: EventEmitter<Item> = new EventEmitter<Item>()

    constructor() { }

    ngOnInit() {

        // DOM element where the Timeline will be attached
        let container = document.getElementById('timeline')

        let values = []
        let id = 0
        for (let item of this.items.getList({ 'hasAttribute': 'date' })) {

            // Ignore no named items
            let name = item.data['name']
            if (name == undefined)
                continue

            values.push({
                id: id++,
                content: name,
                start: item.data['date'],
            })
        }

        // Create a DataSet (allows two way data-binding)
        let items = new Vis.DataSet(values)

        // Configuration for the Timeline
        let options = {}

        // Create a Timeline
        let timeline = new Vis.Timeline(container, items, options)
    }

    ngOnDestroy() { }

    onItem(item: Item) {
	this.open.emit(item)
    }
}
