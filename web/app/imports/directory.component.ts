import { Component } from '@angular/core';
import { ImportsService, Directory } from './imports.service';
import { ApiService } from '../api.service';

@Component({
    selector: 'directory',
    templateUrl: './directory.component.html'
})

export class DirectoryComponent {

    public paths: string[] = []
    public path: string
    public isRecursive: boolean

    constructor(private importsService: ImportsService,
        private apiService: ApiService) {

        // Get configuration import
        importsService.getImportsConfig("directory")
            .subscribe(config => {
                if (config === undefined)
                    return

                // Get global paths
                let paths: string[] = config["*"]
                if (paths == undefined)
                    paths = []

                // Add specific collection paths
                let collectionName: string = apiService.getCollectionName()
                if (collectionName != undefined) {
                    let collectionPaths: string[] = config[collectionName]
                    if (collectionPaths != undefined) {
                        for (var path of collectionPaths) {
                            paths.push(path)
                        }
                    }
                }

                this.paths = paths
            })
    }

    // Create new import collection
    onSubmit() {
        this.importsService.addImport(
            new Directory("", this.path, this.isRecursive))
    }
}
