import { Component, NgZone, OnInit, OnDestroy } from '@angular/core';
import { ExportsService, ExportBase, Directory } from './exports.service';
import { Event } from '../api.service';

declare var jQuery: any;

@Component({
    selector: 'exports',
    templateUrl: './exports.component.html'
})

export class ExportsComponent implements OnInit, OnDestroy {

    public exportTypes: Array<string> = []
    public exports: Map<string, ExportBase[]>

    private events

    constructor(private zone: NgZone,
        private exportsService: ExportsService) {

        // Method called to refresh the export list
        exportsService.setUpdateList(() => {
            this.update()
        })

        this.events = exportsService.subscribeEvents("status")
            .subscribe((e: Event) => {

                let exportBase = this.exportsService.exportsById.get(e.id)
                if (exportBase == undefined) {
                    console.error("Not referenced export with id " + e.id)
                    return
                }

                if (e.event.endsWith("status")) {
                    let item = jQuery("i#" + e.id)
                    if (item == undefined) {
                        console.error("Export with id " + e.id + " not displayed")
                        return
                    }

                    // Set export state
                    exportBase.isRunning = e.data

                    // Status 'TRUE': Rotate refresh logo
                    if (e.data) {
                        item.addClass("rotation")
                    }
                    // Status 'FALSE' : Stop logo rotation
                    else {
                        item.removeClass("rotation")
                    }
                }
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
        this.exportsService.getExports()
            .subscribe((exports: Map<string, ExportBase[]>) => {

                let exportTypes: Array<string> = [];

                exports.forEach((undefined, exportType) => {
                    exportTypes.push(exportType)
                })

                this.zone.run(() => {
                    this.exportTypes = exportTypes
                    this.exports = exports
                })
            })
    }

    onRefresh(item: ExportBase) {
        if (item.isRunning) {
            this.exportsService.stopExport(item)
        } else {
            this.exportsService.startExport(item)
        }
    }

    onDelete(item: ExportBase) {
        this.exportsService.deleteExport(item)
    }
}
