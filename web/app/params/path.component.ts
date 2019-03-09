import {
    Component, NgZone, Input, Output, EventEmitter, OnInit
} from '@angular/core'
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms'
import { ParamsService } from './params.service'

export enum PathMode {
    DIRECTORY_ONLY = 0,
    FILES_ONLY,
    DIRECTORY_AND_FILES,
}

@Component({
    selector: 'params-path',
    templateUrl: './path.component.html',
    providers: [{
        provide: NG_VALUE_ACCESSOR,
        useExisting: ParamsPathComponent,
        multi: true
    }
    ],
})

export class ParamsPathComponent implements OnInit, ControlValueAccessor {

    @Input() type: string
    @Input() name: string
    @Input() mode: PathMode
    @Output() change = new EventEmitter<string>()

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

    // Initial value set to the component
    public writeValue(value: any) {

        if (value == undefined || value == "")
            return

        this.path = value;

        this.zone.run(() => {
            this.displayPath = value
        })
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

    start() {
        this.send(this.path)
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

                    if (path === undefined || path == "") {
                        this.initPath = this.path
                    }

                    // Reset list
                    this.list = []

                    if (this.path !== this.initPath) {

                        let path = this.path

                        if (this.initPath != undefined) {
                            path = this.path.substr(this.initPath.length)
                        }

                        this.list = path.split("/")

                        if (path.charAt(0) == '/') {
                            this.list[0] = "/"
                        }
                    }

                    if (this.initPath != undefined) {
                        this.list.unshift(this.initPath)
                    }

                    this.directories = result["directories"]

		    if (this.mode != PathMode.DIRECTORY_ONLY) {
			this.files = result["files"]
		    }

		    if (this.mode != PathMode.FILES_ONLY) {
			this.validatePathChange(this.displayPath)
		    }
		    
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

    validatePathChange(path: string) {
	this.propagateChange(path)
	this.change.emit(path)
    }
 }
