import {
    Component, NgZone, Input, Output, EventEmitter, AfterViewInit
} from '@angular/core'
import { NgSwitch } from '@angular/common'

import { Tweak } from './tweak'

declare var jQuery: any;

@Component({
    selector: 'tweaks-fields',
    templateUrl: './fields.component.html',
})

export class TweaksFieldsComponent implements AfterViewInit {

	@Input() tweak: Tweak
	@Output() update = new EventEmitter()

	public inputValue: string = ""

    ngAfterViewInit() {

		// Init selector
		// jQuery('select').material_select()
	}

	onChange(event) {
		this.update.emit()
	}
}
