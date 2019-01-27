import { Component, NgZone } from '@angular/core'
import { ImportsService } from '../imports.service'
import { ApiService } from '../../api.service'
import { ImportCreateComponent } from './../create.component'
import { Directory } from './directory'

@Component({
    selector: 'directory-create',
    templateUrl: './create.component.html'
})

export class DirectoryCreateComponent extends ImportCreateComponent {

    public paths: string[] = []

    constructor(private zone: NgZone,
		private importsService: ImportsService,
		private apiService: ApiService) {

        super(new Directory(""))

        // // Get configuration import
        // importsService.getImportsConfig("directory")
        //     .subscribe(config => {
        //         if (config === undefined)
        //             return

        //         // Get global paths
        //         let paths: string[] = config["*"]
        //         if (paths == undefined)
        //             paths = []

        //         // Add specific collection paths
        //         let collectionName: string = apiService.getCollectionName()
        //         if (collectionName != undefined) {
        //             let collectionPaths: string[] = config[collectionName]
        //             if (collectionPaths != undefined) {
        //                 for (var path of collectionPaths) {
        //                     paths.push(path)
        //                 }
        //             }
        //         }

        //         this.paths = paths
        //     })
    }

    onSuccess(data: Directory) {
        this.zone.run(() => {
            this.data = new Directory("")
        })
    }
}
