import { Component, Input, Output, EventEmitter, OnInit, ViewChild,
	 NgZone, Renderer } from '@angular/core'
import { ConfigsService } from './configs.service'
import { ConfigMultiComponent } from './multi.component'
import { ConfigRef } from './config_ref'
import { ConfigBase } from './config_base'
import { BaseElement } from '../base'

@Component({
    selector: 'configs',
    templateUrl: './configs.component.html'
})

export class ConfigsComponent implements OnInit {

    @Input() src: string
    @Input() item: BaseElement
    @Input() init: boolean
    @Output() onUpdate = new EventEmitter()
    @Output() onChange = new EventEmitter()
    
    @ViewChild(ConfigMultiComponent) multi;

    public mainConfigNames: string[] = []
    public refMulti: ConfigRef
    constructor(private zone: NgZone,
		private render: Renderer,
		private configsService: ConfigsService) { }

    ngOnInit() {
	if (this.init) {
            this.update()
	}
    }

    forceInit(src: string, item: BaseElement) {
	this.src = src
	this.item = item
	this.update()
    }
    
    update() {

        this.configsService.getConfig(this.src, this.item)
            .subscribe((cfg: ConfigBase) => {

                let refs = cfg.getRefs();
                let refsMainList: string[] = []

                for (let idx in refs) {

                    let ref = refs[idx]

                    switch (ref.type) {
                        case "struct":
                        case "map":
                            refsMainList.push(ref.name)
                            break
                    }

                }

                this.zone.run(() => {
                    this.mainConfigNames = refsMainList
                    this.onUpdate.emit()
                })
            })
    }


    onClick(event: any, refSelected: string) {

        // Set collection-items as active
        event.preventDefault()

        for (let item of event.target.parentElement.children) {
            this.render.setElementClass(item, "active", false)
        }

        this.render.setElementClass(event.target, "active", true)
	
        this.configsService.getConfig(this.src, this.item)
            .subscribe((cfg: ConfigBase) => {

                let ref = cfg.getRef(refSelected)

                this.zone.run(() => {
                    this.multi.onUpdate(ref)
                    this.onUpdate.emit()
                })
            })
    }

    change(event) {
        this.onChange.emit()
    }
}
