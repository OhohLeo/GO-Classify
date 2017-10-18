import { Component, Input, Output, EventEmitter, OnInit, NgZone } from '@angular/core'
import { ConfigsService, ConfigRef } from './configs.service'
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

	constructor(private zone: NgZone) { }

    ngOnInit() {

		switch (this.ref.type)
		{
		case "bool":
			this.inputType = "checkbox"
			break

		case "int":
			this.inputType = "number"
			break

		default:
			this.inputType = this.ref.type
		}
    }

	onChange(event) {

		switch (this.ref.type)
		{
		case "bool":
			this.ref.data = event.target.checked
			break

		case "int":
			this.ref.data = +event.target.value
			break

		case "stringlist":
			this.ref.data = event.list
			break

		default:
			this.ref.data = event.target.value
			break
		}


		console.log("STRINGLIST", this.ref.name, this.ref.data)
		this.update.emit(this.ref)
	}
}
