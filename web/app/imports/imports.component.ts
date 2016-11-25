import { Component, NgZone } from '@angular/core';
import { ImportsService, ImportBase, Directory } from './imports.service';

@Component({
    selector: 'imports',
    templateUrl: './imports.component.html'
})

export class ImportsComponent {

    public importTypes: Array<string> = [];
    public imports: Map<string, ImportBase[]>;

    constructor(private zone: NgZone,
        private importsService: ImportsService) {

        // Method call to refresh the list
        importsService.setUpdateList(() => {
            this.update()
        })
    }

    ngOnInit() {
        this.update();
    }

    update() {
        this.importsService.getImports()
            .subscribe((imports: Map<string, ImportBase[]>) => {

                let importTypes: Array<string> = [];

                imports.forEach(function (undefined, importType) {
                    importTypes.push(importType)
                })

                this.zone.run(() => {
                    this.importTypes = importTypes
                    this.imports = imports
                })
            })
    }

    onRefresh(item: ImportBase) {
        console.log("REFRESH", item)
    }

    onDelete(item: ImportBase) {
        this.importsService.deleteImport(item)
    }
}
