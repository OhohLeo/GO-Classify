import { Component } from '@angular/core';
import { ImportsService, ImportBase } from './imports.service';
import { ApiService } from '../api.service';

export class Directory extends ImportBase {

    constructor(public id: string,
        public path: string,
        public isRecursive: boolean) {

        super("directory", id);

        if (isRecursive === undefined) {
            this.isRecursive = false
        }
    }

    getParams(): any {
        return {
            "path": this.path,
            "is_recursive": this.isRecursive ? true : false
        }
    }

    display(): string {
        return this.path.concat(this.isRecursive == true ? "/**" : "")
    }

    compare(i: Directory): boolean {
        return super.compare(i)
            && this.path === i.path
            && this.isRecursive == i.isRecursive
    }
}

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

		// Subscribe to convert received data
		importsService.addConvertToImport("email", this.onConvert)

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

	onConvert(id: string, params): ImportBase {

		if (typeof params != 'object') {
			console.error("Unsupported directory parameters!")
			return undefined
		}

		let path = params['path']
		if (typeof path != 'string') {
			console.error("Unsupported 'path' directory parameters!")
			return undefined
		}

		let isRecursive = params['is_recursive']
		if (isRecursive !== undefined && typeof isRecursive != 'boolean') {
			console.error("Unsupported 'is_recursive' directory parameters!")
			return undefined
		}

		return new Directory(id, path, isRecursive)
	}
}
