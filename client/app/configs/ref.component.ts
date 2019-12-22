import { Component, Input, Output, EventEmitter, OnInit, NgZone } from '@angular/core'
import { ConfigsService } from './configs.service'
import { ConfigRef } from './config_ref'
import { CfgStringList } from '../tools/stringlist.component'

@Component({
    selector: 'config-ref',
    templateUrl: './ref.component.html'
})

export class ConfigRefComponent implements OnInit {

    @Input() ref: ConfigRef
    @Output() update = new EventEmitter<ConfigRef>()

    public boolStatus: string
    public inputType: string
    public isStringList: boolean = false

    private pathType: string
    private pathName: string

    constructor(private zone: NgZone) { }

    ngOnInit() {

        switch (this.ref.type) {
            case "bool":
                this.inputType = "checkbox"
                break

            case "int":
                this.inputType = "number"
                break

            case "string":
                if (this.ref.data == undefined)
                    this.ref.data = ""
                break

            case "path":
                this.pathType = "collections/:name"
                this.pathName = "config/" + this.ref.getPathName()
                break
        }

        // Set default input type
        if (this.inputType == undefined)
            this.inputType = this.ref.type
    }

    onChange(event) {

        switch (this.ref.type) {
            case "bool":
                this.ref.data = event.target.checked
                break

            case "int":
                this.ref.data = +event.target.value
                break

            case "stringlist":
                this.ref.data = event.list
                break

            case "path":
                this.ref.data = event
                break

            default:
                this.ref.data = event.target.value
                break
        }

        console.log("Update", this.ref.name, this.ref.data)
        this.update.emit(this.ref)
    }
}
