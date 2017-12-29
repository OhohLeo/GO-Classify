import {
    Component, NgZone, OnInit, OnDestroy,
    ContentChildren, QueryList, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ExportsService, ExportBase } from './exports.service'
import { ExportCreateComponent } from './create.component'
import { ApiService, Event } from '../api.service'
import { Convert2File } from './file/file';

declare var jQuery: any;

@Component({
    selector: 'exports',
    templateUrl: './exports.component.html'
})

export class ExportsComponent implements OnInit, OnDestroy {

    public refs: Array<string> = []
    public refs2Display: Array<string> = []
    public exports: Map<string, ExportBase[]>
    public currentRef: string = "all"

    public createComponent: ExportCreateComponent

    private events

    constructor(private zone: NgZone,
        private render: Renderer,
        private apiService: ApiService,
        private exportsService: ExportsService) {

        // Refresh the import ref list
        apiService.getReferences()
            .subscribe((references) => {
                this.refs = references["exports"]
            })

        // Method called to refresh the export list
        exportsService.setUpdateList(() => {
            this.update()
        })

        // Subscribe to convert received data
        exportsService.addConvertToExport("file", Convert2File)


        this.events = exportsService.subscribeEvents("status")
            .subscribe((e: Event) => {

                let exportBase = this.exportsService.exportsByName.get(e.name)
                if (exportBase == undefined) {
                    console.error("Not referenced export with name " + e.name)
                    return
                }

                if (e.event.endsWith("status")) {
                    let item = jQuery("i#" + e.name)
                    if (item == undefined) {
                        console.error("Export with name " + e.name + " not displayed")
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
                this.zone.run(() => {
                    this.exports = exports
                })
            })
    }

    onRef(event: any, ref: string) {

        // Set collection-items as active
        event.preventDefault()

        for (let item of event.target.parentElement.children) {
            this.render.setElementClass(item, "active", false)
        }

        this.render.setElementClass(event.target, "active", true)

        this.zone.run(() => {
            this.refs2Display = (ref === "all") ? this.refs : [ref]
            this.currentRef = ref
        })
    }

    onExportCreated(exportCreated) {
        this.createComponent = exportCreated
    }

    // Create new export collection
    onSubmit() {

        if (this.createComponent === undefined) {
            console.error("export created component not found", this.currentRef)
            return
        }

        this.exportsService.addExport(
            this.createComponent.data,
            (params) => {
				return this.createComponent.onParams(params)
			},
            (newExport) => {
                return this.createComponent.onSuccess(newExport)
            })
    }

    onRefresh(item: ExportBase) {
        if (item.isRunning) {
            this.exportsService.stopExport(item)
        } else {
            this.exportsService.forceExport(item)
        }
    }

    onConfig(item: ExportBase) {
    }

    onDelete(item: ExportBase) {
        this.exportsService.deleteExport(item)
    }
}
