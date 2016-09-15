import {Component} from '@angular/core';
import {ClassifyService} from './../classify.service';

@Component({
    selector: 'import-directory',
    templateUrl: 'app/imports/directory.component.html'
})

export class ImportDirectoryComponent {

    public directories: Directory[] = []
    public path : string
    public isRecursive: boolean

    constructor (private classifyService: ClassifyService) {
        //this.getImports()
    }

    // Get whole list of imports
    getImports() {
        this.classifyService.getCollectionImport().subscribe(
            imports => {

                if (imports == undefined
                    || imports["directory"] == undefined) {
                    return
                }

                this.directories = []

                for (let name in imports["directory"]) {
                    let d = imports["directory"][name]
                    this.directories.push(
                        new Directory(name, d.path, d.isRecursive))
                }
            });
    }

    // Analyse directory specified
    onReload(directory: Directory) {
        console.log("launch", directory)
    }

    // Delete existing directory
    onDelete(directory: Directory) {
        console.log("delete", directory)
    }

    // Create new import collection
    onSubmit() {

        // Check if it didn't match an already existing import
        // collection
        // for (let directory of this.directories) {
        //     if (directory.path == this.path) {
        //         //this.classifyService.onError(
        //         console.error("Import path directory already existing!")
        //         return
        //     }
        // }

        console.log("new Directory", this.path, this.isRecursive)

        // this.classifyService.newCollectionImport({
        //     "type": "directory",
        //     "params": {
        //         "path": this.path,
        //         "isRecursive": this.isRecursive
        //     }
        // })
        //     .subscribe(status => {

        //         this.directories.push(
        //             new Directory("test", this.path, this.isRecursive))
        //         console.log("OK", status)
        //     })
    }
}

export class Directory {
    constructor(public name: string,
                public path: string,
                public isRecursive: boolean) {}
}
