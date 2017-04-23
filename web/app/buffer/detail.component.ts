import { Component, Input, Output, EventEmitter, NgZone, OnInit, OnDestroy } from '@angular/core'

import { BufferService, BufferItem } from './buffer.service'
import { Event, Item } from '../api.service'
import { ConfigService, ConfigBase } from '../config/config.service';
import { ApiService } from '../api.service';


declare var jQuery: any

@Component({
    selector: 'detail',
    templateUrl: './detail.component.html',
})


export class DetailComponent implements OnInit, OnDestroy {

	@Input() collection: string
    @Input() buffer: BufferItem
    @Output() close: EventEmitter<any> = new EventEmitter();
	private needNameDetails: boolean = false


    private action: any

	constructor(private zone: NgZone,
				private configService : ConfigService) {
	}

    ngOnInit() {

		let buffer = this.buffer;

        this.action = jQuery('div#detail').modal({
            complete: () => { this.close.emit() }
        })

		console.log("INIT", buffer)

        this.action.modal("open")
    }

    ngOnDestroy() {
        this.action.modal("close")
        this.close.emit()
    }

	getNameDetails() {
		this.zone.run(() => {
			this.needNameDetails = !this.needNameDetails;
		})
	}

	onUpdate(buffer: BufferItem) {
		this.zone.run(() => {
			console.log("DETAIL ZONE UPDATE?")
			this.buffer = buffer
		})
	}

	onChange(event) {
		this.configService.onChange(this.collection, event)
    }

    // Set the fields depending on the collection

    // Set the import list

    // Set the websites list

    // Open the detail modal
}
