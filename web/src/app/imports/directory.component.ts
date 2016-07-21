import {Component} from '@angular/core';
import {ClassifyService} from './../classify.service';

@Component({
    selector: 'import-directory',
    templateUrl: 'app/imports/directory.component.html'
})

export class ImportDirectoryComponent {

    public newDirectory = new Directory()

    constructor (private classifyService: ClassifyService) {

    }

    onSubmit() {
        console.log(this.newDirectory)
    }
}

class Directory {
    path: string
    isRecursive: boolean
}
