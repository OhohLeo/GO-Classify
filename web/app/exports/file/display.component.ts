import { Component, Input, NgZone } from '@angular/core'
import { File } from './file'

@Component({
    selector: 'file-display',
    templateUrl: './display.component.html'
})

export class FileDisplayComponent {
	@Input() file : File
}
