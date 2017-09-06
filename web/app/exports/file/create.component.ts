import { Component, NgZone } from '@angular/core'
import { ApiService } from '../../api.service'
import { File } from './file'

@Component({
    selector: 'file-create',
    templateUrl: './create.component.html'
})

export class FileCreateComponent {
	public file : File
}
