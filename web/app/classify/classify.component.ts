import { Component, OnInit, OnDestroy } from '@angular/core';

import { ClassifyService } from './classify.service';
import { Event, Item } from '../api.service';

declare var jQuery: any;

@Component({
    selector: 'classify',
    templateUrl: './classify.component.html',
})

export class ClassifyComponent implements OnInit, OnDestroy {

	private action: any
    private events: any

    constructor(private classifyService: ClassifyService) {

        this.events = classifyService.subscribeEvents("status")
            .subscribe((e: Event) => {
				console.log("TEST", e)
			})
	}

    ngOnInit() {
		this.action = jQuery('div#classify').modal()
    }

    startModal() {
        this.action.modal("open")
    }

    ngOnDestroy() {
        this.action.modal("close")
    }
}
