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
		this.configService.onChange(this.collection, event)
    }
}
