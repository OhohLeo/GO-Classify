import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'
import { Collection } from '../collection'
import { Item } from '../item'
import * as Vis from 'vis'

@Component({
    selector: 'timeline-collection',
    templateUrl: './timeline.component.html',
})

export class TimeLineCollectionComponent implements OnInit, OnDestroy {

    @Input() collection: Collection
    @Input() items: Item[]

    constructor() { }

    ngOnInit() {

        // DOM element where the Timeline will be attached
        var container = document.getElementById('timeline');

        // Create a DataSet (allows two way data-binding)
        var items = new Vis.DataSet([
            { id: 1, content: 'item 1', start: '2013-04-20' },
            { id: 2, content: 'item 2', start: '2013-04-14' },
            { id: 3, content: 'item 3', start: '2013-04-18' },
            { id: 4, content: 'item 4', start: '2013-04-16', end: '2013-04-19' },
            { id: 5, content: 'item 5', start: '2013-04-25' },
            { id: 6, content: 'item 6', start: '2013-04-27' }
        ]);

        // Configuration for the Timeline
        var options = {};

        // Create a Timeline
        var timeline = new Vis.Timeline(container, items, options);
    }

    ngOnDestroy() { }

}
