import {
    Component, NgZone, Input, OnInit, OnDestroy
} from '@angular/core'
import {
    ControlValueAccessor, NG_VALUE_ACCESSOR
} from '@angular/forms'

import { ParamsService } from './params.service'

declare var jQuery: any

@Component({
    selector: 'params-path',
    templateUrl: './path.component.html',
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: ParamsPathComponent,
            multi: true
        }
    ],
})

export class ParamsPathComponent implements OnInit, ControlValueAccessor, OnDestroy {

    @Input() type: string
    @Input() name: string
    private path: string
    private initPath: string
    private displayPath: string

    public list: string[] = []
    public directories: string[] = []
    public files: string[] = []

    private action: any

    // the method set in registerOnChange, it is just
    // a placeholder for a method that takes one parameter,
    // we use it to emit changes back to the form
    private propagateChange = (_: any) => { };

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

    // this is the initial value set to the component
    public writeValue(value: any) {
        this.path = value;
    }

    // registers 'fn' that will be fired when changes are made
    // this is how we emit the changes back to the form
    public registerOnChange(fn: any) {
        this.propagateChange = fn;
    }
    // not used, used for touch input
    public registerOnTouched() { }

    // change events from the textarea
    private onChange(event) { }

    onModal() {

        if (this.action != undefined) {
            return
        }

        this.send(this.path)

        this.action = jQuery('div#path').modal({
            complete: () => {
                this.stop()
            }
        })

        this.action.modal("open")
    }

    send(path: string) {

        this.paramsService.actionParam(
            this.type, this.name, "path",
            {
                "directory": (path === undefined) ? "" : path
            })
            .subscribe((result) => {

                this.zone.run(() => {

                    this.path = result["current"]
                    this.displayPath = this.path

                    if (path == "") {
                        this.initPath = this.path
                    }

                    // Reset list
                    this.list = []

                    if (this.path != this.initPath) {
                        this.list = this.path.substr(
                            this.initPath.length).split("/")
                    }

                    this.list.unshift(this.initPath)

                    this.directories = result["directories"]
                    this.files = result["files"]
                })
            })
    }

    onSelectBack(inputIdx: number) {

        let path: string = ""
        for (let idx = 0; idx < inputIdx; idx++) {
            let name = this.list[idx]
            path += name + (name.slice(-1) != "/") ? "/" : ""
        }

        this.send(path + this.list[inputIdx])
    }

    onSelect(directory: string) {

        if (this.displayPath.slice(-1) != "/") {
            this.displayPath += "/"
        }

        this.send(this.displayPath + directory)
    }

    onValidate() {
        this.stop()
    }

    stop() {

        if (this.action == undefined) {
            return
        }

        this.action.modal('close')
        this.action = undefined

        // update the form
        this.propagateChange(this.path);
    }
}
