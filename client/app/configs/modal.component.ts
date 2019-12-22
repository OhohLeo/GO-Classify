import {
    Component, NgZone, Input, ViewChild, OnInit, OnDestroy
} from '@angular/core'
import { ConfigsService } from './configs.service'
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
    public validate: boolean = false
    
    private action: any

    constructor(private zone: NgZone,
		private configsService: ConfigsService) { }
    
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

    onUpdate() {
        this.zone.run(() => {
            this.validate = false
        })
    }

    onChange(event) {
        this.zone.run(() => {
            this.validate = true
        })
    }

    onSubmit(event) {
        event.preventDefault()
        this.configsService.setConfig(this.src, this.item.name)
            .subscribe((res) => {
                this.zone.run(() => {
                    this.validate = false
                })
            })
    }
}
