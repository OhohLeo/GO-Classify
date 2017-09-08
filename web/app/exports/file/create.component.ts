import { Component, NgZone } from '@angular/core'
import { ExportsService } from './../exports.service'
import { File } from './file'

@Component({
    selector: 'file-create',
    templateUrl: './create.component.html'
})

export class FileCreateComponent {
	public file : File

    constructor(private zone: NgZone,
				private exportsService: ExportsService) {

		this.file = new File("")
	}

	onSuccess(file: File) {
		this.zone.run(() => {
			this.file = new File("")
		})
	}

	// Create new export collection
    onSubmit() {
        this.exportsService.addExport(
			this.file,
			undefined,
			(file) => { return this.onSuccess(file) })
    }
}
