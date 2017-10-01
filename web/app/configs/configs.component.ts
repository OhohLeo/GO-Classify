import { Component, NgZone } from '@angular/core'

@Component({
    selector: 'configs',
    templateUrl: './configs.component.html'
})

export class ConfigsComponent {

    constructor(private zone: NgZone) { }

    ngOnInit() {
        // this.update()
    }

    update() {
        // this.configService.getConfigs(this.collection)
        //     .subscribe((config: ConfigBase) => {
        //         this.zone.run(() => {
        //             this.config = config
        //         })
        //     })
    }

    onChange(event) {

        // switch (event.target.name) {
        //     case "'enableStore'":
        //         event.target.value = this.config.enableStore
        //         break;
        //     case "'enableBuffer'":
        //         event.target.value = this.config.enableBuffer
        //         break;
        // }

        // this.configService.onChange(this.collection, event)
    }
}
