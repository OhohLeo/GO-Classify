import {Component} from '@angular/core';
import {ImportsService, ImportInterface} from './imports.service';

@Component({
    selector: 'import-directory',
    templateUrl: 'app/imports/directory.component.html'
})

export class ImportDirectoryComponent {

    public path : string
    public isRecursive: boolean

    constructor (private importsService: ImportsService) {}

    // Create new import collection
    onSubmit() {
        this.importsService.newImport(
            new Directory(this.path, this.isRecursive))
    }
}

export class Directory implements ImportInterface {
    constructor(public path: string,
                public isRecursive: boolean) {

        if (isRecursive === undefined) {
            this.isRecursive = false
        }
    }

    getType() : string {
        return "directory"
    }

    getParams() : any {
        return {
            "path": this.path,
            "isRecursive": this.isRecursive ? true : false
        }
    }

    compare(i : Directory) : boolean {
        return this.path === i.path
            && this.isRecursive == i.isRecursive
    }
}
