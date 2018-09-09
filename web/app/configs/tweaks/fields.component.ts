import {
    Component, NgZone, Input, Output, EventEmitter, OnInit
} from '@angular/core'
import { NgSwitch } from '@angular/common'

import { Tweak } from './tweak'

declare var jQuery: any;

@Component({
    selector: 'tweaks-fields',
    templateUrl: './fields.component.html',
})

export class TweaksFieldsComponent implements OnInit {
	@Input() tweak: Tweak

    ngOnInit() {

		// Init selectors
		let selectors = jQuery('select')
		selectors.each((idx) => {
			selectors[idx].material_select();
		})
	}
}
