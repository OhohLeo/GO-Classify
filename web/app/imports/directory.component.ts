import { Component, NgZone } from '@angular/core';
import { ImportsService, ImportInterface } from './imports.service';
import { ClassifyService } from '../classify.service';

@Component({
    selector: 'directory',
    templateUrl: './directory.component.html'
})

export class DirectoryComponent {

    public paths: string[] = []
    public path: string
    public isRecursive: boolean

    constructor(private zone: NgZone,
        private importsService: ImportsService,
        private classifyService: ClassifyService) {

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
                let collectionName: string = classifyService.getCollectionName()
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

    getType(): string {
        return "directory"
    }

    getParams(): any {
        return {
            "path": this.path,
            "isRecursive": this.isRecursive ? true : false
        }
    }

    compare(i: Directory): boolean {
        return this.path === i.path
            && this.isRecursive == i.isRecursive
    }
}
