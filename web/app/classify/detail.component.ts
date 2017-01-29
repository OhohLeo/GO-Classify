import { Component, Input, Output, EventEmitter, NgZone, OnInit, OnDestroy } from '@angular/core'

import { ClassifyService, ClassifyItem } from './classify.service'
import { Event, Item } from '../api.service'

declare var jQuery: any

@Component({
    selector: 'detail',
    templateUrl: './detail.component.html',
})


export class DetailComponent implements OnInit, OnDestroy {

	@Input() item: ClassifyItem
	@Output() close: EventEmitter<any> = new EventEmitter();

	private action: any

    ngOnInit() {
		console.log(this.item);

		this.action = jQuery('div#detail').modal({
			complete : () => { this.close.emit() }
		})

		this.action.modal("open")
    }

    ngOnDestroy() {
        this.action.modal("close")
		this.close.emit()
    }

	// Set the fields depending on the collection

	// Set the import list

	// Set the websites list

	// Open the detail modal
}
