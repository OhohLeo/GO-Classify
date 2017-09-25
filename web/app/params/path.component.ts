import {
    Component, NgZone, OnInit, OnDestroy
} from '@angular/core'
import { ParamsService } from './params.service'

declare var jQuery: any

@Component({
    selector: 'params-path',
    templateUrl: './path.component.html'
})

export class ParamsPathComponent implements OnInit, OnDestroy {

    public path: string = ""
    private displayPath: string

    public directories: string[] = []
    public files: string[] = []

    private action: any


    constructor(private zone: NgZone,
        private paramsService: ParamsService) { }

    ngOnInit() {

        let displayPath: string

        if (this.path === "") {
            displayPath = this.path
        }

        this.zone.run(() => {
            this.displayPath = "Select path"
        })
    }

    ngOnDestroy() {
        this.stop
    }

    onModal() {

        if (this.action != undefined) {
            return
        }

        this.paramsService.actionParam("imports", "directory", "path", this.path)
            .subscribe((result) => {
                console.log(result)
            })


        this.action = jQuery('div#path').modal({
            complete: () => {
                this.stop()
            }
        })

        this.action.modal("open")
    }

    stop() {

        if (this.action == undefined) {
            return
        }

        this.action.modal('close')
        this.action = undefined
    }
}
