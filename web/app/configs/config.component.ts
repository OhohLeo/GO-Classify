import {
    Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy
} from '@angular/core'

declare var jQuery: any

@Component({
    selector: 'config',
    templateUrl: './config.component.html',
})

export class ConfigComponent implements OnInit, OnDestroy {

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

    start() {
        this.action.modal("open")
    }


    stop() {}

}
