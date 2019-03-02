import {
    Component, NgZone, OnInit, AfterViewInit, OnDestroy,
    ViewChild, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ApiService, Event } from './../api.service'
import { ConfigModalComponent } from './../configs/modal.component'
import { ImportsService } from './imports.service'
import { BaseCreateComponent } from '../tools/base_create.component'
import { DirectoryCreateComponent } from './directory/create.component'
import { BaseElement } from '../base'

declare var jQuery: any;

@Component({
    selector: 'imports',
    templateUrl: './imports.component.html'
})

export class ImportsComponent implements OnInit, OnDestroy {

    @ViewChild(ConfigModalComponent) config: ConfigModalComponent

    public refs: Array<string> = []
    public refs2Display: Array<string> = []
    public imports: Map<string, BaseElement[]>
    public currentRef: string = "all"

    public createComponent: BaseCreateComponent

    private events

    constructor(private zone: NgZone,
        private render: Renderer,
        private apiService: ApiService,
        private importsService: ImportsService) {

        // Refresh the import ref list
        apiService.getReferences()
            .subscribe((references) => {
                this.refs = Object.keys(references["imports"])
            })

        // Refresh the import list
        importsService.setUpdateList(() => {
            this.update()
        })

        this.events = importsService.subscribeEvents("status")
            .subscribe((e: Event) => {

                let importBase = this.importsService.importsByName.get(e.name)
                if (importBase == undefined) {
                    console.error("Not referenced import with name " + e.name)
                    return
                }

                if (e.event.endsWith("status")) {
                    let item = jQuery("i#" + e.name)
                    if (item == undefined) {
                        console.error("Import with name " + e.name + " not displayed")
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
        this.update()
    }

    ngOnDestroy() {
        if (this.events != undefined) {
            this.events.unsubscribe()
            this.events = undefined
        }
    }

    update() {
        this.importsService.getImports()
            .subscribe((imports: Map<string, BaseElement[]>) => {
                this.zone.run(() => {
                    this.imports = imports
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

    onCreated(importCreated) {
	this.createComponent = importCreated
    }

    // Create new import collection
    onSubmit() {

        if (this.createComponent === undefined) {
            console.error("import created component not found", this.currentRef)
            return
        }

        this.importsService.addImport(
            this.createComponent.data,
            (params) => {
    		return this.createComponent.onParams(params)
    	    },
            (newImport) => {
                return this.createComponent.onSuccess(newImport)
            })
    }

    onRefresh(item: BaseElement) {
        if (item.isRunning) {
            this.importsService.stopImport(item)
        } else {
            this.importsService.startImport(item)
        }
    }

    onConfig(item: BaseElement) {
	this.config.start("imports", item)
    }

    onDelete(item: BaseElement) {
        this.importsService.deleteImport(item)
    }
}
