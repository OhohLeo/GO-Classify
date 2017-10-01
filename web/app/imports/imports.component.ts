import {
    Component, NgZone, OnInit, AfterViewInit, OnDestroy,
    ViewChildren, QueryList, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ApiService, Event } from './../api.service'
import { ImportsService, ImportBase } from './imports.service'
import { ImportCreateComponent } from './create.component'
import { Convert2Imap } from './imap/imap';
import { Convert2Directory } from './directory/directory';
import { DirectoryCreateComponent } from './directory/create.component';

declare var jQuery: any;

@Component({
    selector: 'imports',
    templateUrl: './imports.component.html'
})

export class ImportsComponent implements OnInit, OnDestroy {

    public refs: Array<string> = []
    public refs2Display: Array<string> = []
    public imports: Map<string, ImportBase[]>
    public currentRef: string = "all"

    public importName: string
    public createComponent: ImportCreateComponent

    private events

    constructor(private zone: NgZone,
        private render: Renderer,
        private apiService: ApiService,
        private importsService: ImportsService) {

        // Refresh the import ref list
        apiService.getReferences()
            .subscribe((references) => {
                this.refs = references["imports"]
            })

        // Refresh the import list
        importsService.setUpdateList(() => {
            this.update()
        })

        // Subscribe to convert received data
        importsService.addConvertToImport("imap", Convert2Imap)
        importsService.addConvertToImport("directory", Convert2Directory)

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
            .subscribe((imports: Map<string, ImportBase[]>) => {
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

    onImportCreated(importCreated) {
        this.createComponent = importCreated
    }

    // Create new import collection
    onSubmit() {

        if (this.createComponent === undefined) {
            console.error("import created component not found", this.currentRef)
            return
        }

        this.importsService.addImport(
            this.importName,
            this.createComponent.data,
            (params) => { return this.createComponent.onParams(params) },
            (newImport) => {

                this.zone.run(() => {
                    this.importName = ""
                })

                return this.createComponent.onSuccess(newImport)
            })
    }

    onRefresh(item: ImportBase) {
        if (item.isRunning) {
            this.importsService.stopImport(item)
        } else {
            this.importsService.startImport(item)
        }
    }

    onConfig(item: ImportBase) {
    }

    onDelete(item: ImportBase) {
        this.importsService.deleteImport(item)
    }
}
