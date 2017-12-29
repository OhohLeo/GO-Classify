import { Component, OnInit, Output, EventEmitter } from '@angular/core'

@Component({
	selector: 'export-create'
})

export class ExportCreateComponent implements OnInit {

	@Output() onCreated = new EventEmitter();

	constructor(public data : any) {
	}

	ngOnInit() {
		this.onCreated.emit(this)
	}

	onParams(data: any) {}
	onSuccess(data: any) {}
}
