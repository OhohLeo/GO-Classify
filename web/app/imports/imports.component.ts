import { Component, NgZone, OnInit, OnDestroy } from '@angular/core';
import { ImportsService, ImportBase, Directory } from './imports.service';
import { Event } from '../classify.service';

@Component({
    selector: 'imports',
    templateUrl: './imports.component.html'
})

export class ImportsComponent implements OnInit, OnDestroy {

    public importTypes: Array<string> = [];
    public imports: Map<string, ImportBase[]>;

    private events

    constructor(private zone: NgZone,
        private importsService: ImportsService) {

        // Method called to refresh the import list
        importsService.setUpdateList(() => {
            this.update()
        })

        // Method called to refresh import status
        importsService.setImportStatus((item: ImportBase, isStart: boolean) => {
            this.onImportStatus(item, isStart)
        })

        this.events = importsService.subscribeEvents("status")
            .subscribe((e: Event) => {
                console.log("IMPORT EVENT!", e)
            })

    }

    ngOnInit() {
        this.update();
    }

    ngOnDestroy() {
        if (this.events != undefined) {
            this.events.unsubscribe()
            this.events = undefined
        }
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
        this.importsService.startImport(item)
    }

    onImportStatus(item: ImportBase, status: boolean) {
        console.log(item, status)
    }

    onDelete(item: ImportBase) {
        this.importsService.deleteImport(item)
    }
}
