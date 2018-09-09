import {
    Component, NgZone, Input, Output, EventEmitter, ViewChild, OnInit, OnDestroy
} from '@angular/core'
import { TweaksComponent } from './tweaks/tweaks.component'
import { BaseElement } from '../base'

declare var jQuery: any

@Component({
    selector: 'config',
    templateUrl: './config.component.html',
})

export class ConfigComponent implements OnInit, OnDestroy {

	@ViewChild(TweaksComponent) tweaks;

	public item: BaseElement

    private action: any

	ngOnInit() {
        this.action = jQuery('div#config').modal({
            complete: () => {
                this.stop();
            }
        })
    }

    ngOnDestroy() {
        this.stop()
    }

    start(item: BaseElement) {
		this.item = item
		this.tweaks.start(item)
        this.action.modal("open")
    }


    stop() {}

}
