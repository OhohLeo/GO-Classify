import { Injectable } from '@angular/core'
import { Response } from '@angular/http'

import { Observable } from 'rxjs/Rx'

import { ApiService, Event } from '../../api.service'
import { ImportsService } from '../../imports/imports.service'
import { BaseElement } from '../../base'

@Injectable()
export class TweaksService {

    constructor(private apiService: ApiService,
		private importsService: ImportsService) { }

    // Returns inputs & outputs references depending on item type.
    // If type is "imports" :
    //  - inputs are import items
    //  - outputs are collection items
    // If type is "exports" :
    //  - inputs are collection items
    //  - outputs are export items
    getReferences(item: BaseElement) {

        return new Observable(observer => {

	    switch (item.getType()) {

	    case "imports":
		Observable.combineLatest(
		    this.importsService.getReferences(item),
		    this.apiService.getCollectionReferences()
		).subscribe(([inputs, outputs]) => {
		    observer.next([inputs, outputs])
		})
		break

	    case "exports":
		Observable.combineLatest(
		    this.apiService.getCollectionReferences()
		    //this.exportsService.getReferences(item),
		).subscribe(([inputs, outputs]) => {
		    observer.next([inputs, outputs])
		})
		break

	    default:
		console.error("Tweak item not possible on '" + item.getType() + "'")
	    }
	})
    }
}
