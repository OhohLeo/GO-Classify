import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core';
import { ConfigService, ConfigBase } from './config.service';
import { StringListEvent } from './stringlist.component';

@Component({
    selector: 'config',
    templateUrl: './config.component.html'
})

export class ConfigComponent {
    @Input() collection: string
    public config: ConfigBase

    constructor(private zone: NgZone,
        private configService: ConfigService) {
    }

    ngOnInit() {
        this.update()
    }

    update() {
        this.configService.getConfigs(this.collection)
            .subscribe((config: ConfigBase) => {
                this.zone.run(() => {
                    console.log("GET CONFIG", config)
                    this.config = config
                })
            })
    }

    onChange(event) {

        let name, action, value

        if (event instanceof StringListEvent) {
            name = event.name
            action = event.action
            value = event.list
        } else {
            name = event.target.name
            switch (event.target.type) {
                case "number":
                    value = Number(event.target.value)
                default:
                    value = event.target.value
            }

            console.log(event.target.type, value)
        }

        let observable = this.configService.setConfig(
            this.collection, name, action, value)
        if (observable != undefined) {
            observable.subscribe((status) => {
            })
        }
    }

    onApply(event) {
        console.log("ON APPLY")
    }
}
