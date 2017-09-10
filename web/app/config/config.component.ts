import { Component, NgZone, Input, OnInit, OnDestroy, Renderer } from '@angular/core';
import { ConfigService, ConfigBase } from './config.service';
import { StringListEvent } from '../tools/stringlist.component';

@Component({
    selector: 'config',
    templateUrl: './config.component.html'
})

export class ConfigComponent {
    @Input() collection: string
    public config: ConfigBase
    public currentRef: string
    public collectionConfig: any
    public importsConfig: any
    public exportsConfig: any

    constructor(private zone: NgZone,
        private render: Renderer,
        private configService: ConfigService) {
    }

    ngOnInit() {
        this.update()
    }

    update() {
        this.configService.getConfigs(this.collection)
            .subscribe((config: ConfigBase) => {
                this.zone.run(() => {
                    this.config = config
                })
            })
    }

    onRef(event: any, ref: string) {

        // Set collection-items as active
        event.preventDefault()

        for (let item of event.target.parentElement.children) {
            this.render.setElementClass(item, "active", false)
        }

        this.render.setElementClass(event.target, "active", true)

        this.zone.run(() => {
            this.currentRef = ref
        })
    }

    onChange(event) {

        switch (event.target.name) {
            case "'enableStore'":
                event.target.value = this.config.enableStore
                break;
            case "'enableBuffer'":
                event.target.value = this.config.enableBuffer
                break;
        }

        this.configService.onChange(this.collection, event)
    }
}
