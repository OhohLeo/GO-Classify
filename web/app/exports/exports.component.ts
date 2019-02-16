import {
    Component, NgZone, OnInit, OnDestroy,
    ContentChildren, QueryList, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ExportsService } from './exports.service'
import { BaseCreateComponent } from '../tools/base_create.component'
import { ApiService, Event } from '../api.service'

import { BaseElement } from '../base'

declare var jQuery: any;

@Component({
    selector: 'exports',
    templateUrl: './exports.component.html'
})

export class ExportsComponent implements OnInit, OnDestroy {

    public refs: Array<string> = []
    public refs2Display: Array<string> = []
    public exports: Map<string, BaseElement[]>
    public currentRef: string = "all"

    public createComponent: BaseCreateComponent

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
            .subscribe((exports: Map<string, BaseElement[]>) => {
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

    onCreated(exportCreated) {
    //     this.createComponent = exportCreated
    }

    // // Create new export collection
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

    onRefresh(item: BaseElement) {
        if (item.isRunning) {
            this.exportsService.stopExport(item)
        } else {
            this.exportsService.forceExport(item)
        }
    }

    onConfig(item: BaseElement) {
    }

    onDelete(item: BaseElement) {
        this.exportsService.deleteExport(item)
    }
}
