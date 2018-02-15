import { Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core'

import { Item } from './item'

declare var jQuery: any

@Component({
    selector: 'item',
    templateUrl: './item.component.html'
})

export class ItemComponent implements OnInit, OnDestroy {

    @Input() item: Item
    @Output() close = new EventEmitter()

    private action: any
    
    ngOnInit() {
	
        this.action = jQuery('div#item').modal({
            complete: () => {
                this.stop()
            }
        })

	this.start()

	console.log(this.item)
    }
    
    ngOnDestroy() {}

    start() {
        // If has buffer items : open modal
	this.action.modal("open")	
    }

    stop() {
	this.close.emit()
    }
}
