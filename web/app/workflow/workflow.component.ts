import {
    Component, NgZone, OnInit, AfterViewInit, OnDestroy,
    ViewChild, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ApiService, Event } from './../api.service'
import { ImportsService } from '../imports/imports.service'
import { ExportsService } from '../exports/exports.service'
import { WorkflowService } from './workflow.service'

import { CanvasComponent } from './canvas.component'

import { WorkflowType } from './workflow'
import { BaseElement } from '../base'

declare var jQuery: any;

@Component({
    selector: 'workflow',
    templateUrl: './workflow.component.html'
})

export class WorkflowComponent implements OnInit, AfterViewInit, OnDestroy {

    @ViewChild("importsLinks") importLinks: CanvasComponent
    @ViewChild("exportsLinks") exportLinks: CanvasComponent
    
    public workflowType = WorkflowType
    public importRefs: Array<string> = []
    public imports: Map<string, BaseElement[]>
    public exportRefs: Array<string> = []
    public exports: Map<string, BaseElement[]>
    
    constructor(private zone: NgZone,
		private render: Renderer,
		private apiService: ApiService,
		private importsService: ImportsService,
		private exportsService: ExportsService,
		private workflowService: WorkflowService) {

	// Refresh the import/export ref list
        apiService.getReferences()
            .subscribe((references) => {
                this.importRefs = references["imports"]
		this.exportRefs = references["exports"]
            })

	// // Refresh the import/export lists
        // importsService.setUpdateList(() => {
        //     this.updateImports()
        // })

	// exportsService.setUpdateList(() => {
        //     this.updateExports()
        // })
    }

    ngOnInit() {
    }

    ngAfterViewInit() {
	console.log(this.importLinks, this.exportLinks)
    }

    ngOnDestroy() {
    }


    updateAll() {
	this.updateImports()
	this.updateExports()
    }
    
    updateImports() {
        this.importsService.getImports()
            .subscribe((imports: Map<string, BaseElement[]>) => {
                this.zone.run(() => {
                    this.imports = imports
                })
            })
    }

    updateExports() {
	this.exportsService.getExports()
            .subscribe((exports: Map<string, BaseElement[]>) => {
                this.zone.run(() => {
                    this.exports = exports
                })
            })
    }

}
