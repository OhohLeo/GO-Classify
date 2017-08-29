import { Component, NgZone } from '@angular/core'
import { ImportsService, ImportBase } from '../imports.service'
import { ApiService } from '../../api.service'
import { Directory } from './directory'

@Component({
    selector: 'directory-create',
    templateUrl: './create.component.html'
})

export class DirectoryCreateComponent {

	public paths: string[] = []
    public directory : Directory


    constructor(private zone: NgZone,
				private importsService: ImportsService,
				private apiService: ApiService) {

		this.directory = new Directory("")

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

	onSuccess(directory: Directory) {
		this.zone.run(() => {
			this.directory = new Directory("")
		})
	}

    // Create new import collection
    onSubmit() {
        this.importsService.addImport(
			this.directory,
			undefined,
			(directory) => { this.onSuccess(directory) })
    }

}
