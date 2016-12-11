import { Component, NgZone, OnInit, OnDestroy } from '@angular/core';
import { ImportsService, ImportBase, Directory } from './imports.service';
import { Event } from '../classify.service';

declare var jQuery: any;

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

        this.events = importsService.subscribeEvents("status")
            .subscribe((e: Event) => {

				let importBase = this.importsService.importsById.get(e.id)
				if (importBase == undefined)
				{
					console.error("Not referenced import with id "+ e.id)
					return
				}

				if (e.event.endsWith("status"))
				{
					let item = jQuery("i#" + e.id)
					if (item == undefined) {
						console.error("Import with id " + e.id + " not displayed")
						return
					}

					// Set import state
					importBase.isRunning = e.data

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
        this.importsService.getImports()
            .subscribe((imports: Map<string, ImportBase[]>) => {

                let importTypes: Array<string> = [];

                imports.forEach((undefined, importType) => {
					importTypes.push(importType)
                })

                this.zone.run(() => {
                    this.importTypes = importTypes
                    this.imports = imports
                })
            })
    }

    onRefresh(item: ImportBase) {
		if (item.isRunning) {
			this.importsService.stopImport(item)
		} else {
			this.importsService.startImport(item)
		}
    }

    onDelete(item: ImportBase) {
        this.importsService.deleteImport(item)
    }
}
