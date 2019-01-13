import {
    Component, NgZone, Input, ViewChild, OnInit, OnDestroy
} from '@angular/core'
import { ConfigsComponent } from './configs.component'
import { BaseElement } from '../base'

declare var jQuery: any

@Component({
    selector: 'config-modal',
    templateUrl: './modal.component.html',
})

export class ConfigModalComponent implements OnInit, OnDestroy {

    @ViewChild(ConfigsComponent) configs;
    
    public src: string
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

    start(src: string, item: BaseElement) {
	this.src = src
	this.item = item
	this.configs.forceInit(src, item)
        this.action.modal("open")
    }

    stop() {}

}
