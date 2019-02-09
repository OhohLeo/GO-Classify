import { Component, NgZone } from '@angular/core'
import { ExportsService } from './../exports.service'
import { BaseCreateComponent } from '../../tools/base_create.component'
import { File } from './file'

@Component({
    selector: 'file-create',
    templateUrl: './create.component.html'
})

export class FileCreateComponent extends BaseCreateComponent {

    constructor(private zone: NgZone,
		private exportsService: ExportsService) {

	super(new File(""))
    }

    onSuccess(file: File) {
	this.zone.run(() => {
	    this.data = new File("")
	})
    }
}
