import { Component, Input, OnInit } from '@angular/core'
import { Directory } from './directory'

@Component({
    selector: 'directory-display',
    templateUrl: './display.component.html'
})

export class DirectoryDisplayComponent implements OnInit {

	@Input() directory : Directory
	public display: string

    ngOnInit() {

		this.display = this.directory.path

		if (this.directory.isRecursive) {
			this.display += "/**/*"
		}
    }
}
