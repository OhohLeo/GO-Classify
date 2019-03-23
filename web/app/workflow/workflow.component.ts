import {
    ElementRef, Component, NgZone, OnInit, AfterViewInit, OnDestroy,
    ViewChild, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ApiService, Event } from './../api.service'
import { ImportsService } from '../imports/imports.service'
import { ExportsService } from '../exports/exports.service'
import { WorkflowService } from './workflow.service'

import { CanvasComponent } from './canvas.component'

import { WorkflowType } from './workflow'
import { References, Reference } from '../references/reference'
import { BaseElement } from '../base'

@Component({
    selector: 'workflow',
    templateUrl: './workflow.component.html',
    styleUrls: ['./workflow.component.css']
})

export class WorkflowComponent implements OnInit, AfterViewInit, OnDestroy {

    @ViewChild("importsLinks") importLinks: CanvasComponent
    @ViewChild("exportsLinks") exportLinks: CanvasComponent
    
    public workflowType = WorkflowType

    public references: References
    
    public importTypeRefNames: Array<string> = []
    public importTypes: Array<string> = []
    public imports: Map<string, BaseElement[]>
	
    public exportTypeRefNames: Array<string> = []
    public exportTypes: Array<string> = []
    public exports: Map<string, BaseElement[]>
    
    constructor(private zone: NgZone,
		private render: Renderer,
		private eltRef: ElementRef,
		private apiService: ApiService,
		private importsService: ImportsService,
		private exportsService: ExportsService,
		private workflowService: WorkflowService) {

	// Refresh the import/export ref list
        apiService.getReferences()
            .subscribe((references: References) => {
		this.references = references
		this.importTypeRefNames = references.getTypeRefs('imports')
		this.exportTypeRefNames = references.getTypeRefs('exports')
            })
	
	// Refresh the import/export lists
        importsService.setUpdateList(() => {
            this.updateImports()
        })

	exportsService.setUpdateList(() => {
            this.updateExports()
        })
    }

    ngOnInit() {
 	this.updateAll()
    }

    ngAfterViewInit() {
	console.log("WORKFLOW ELT REF", this.eltRef.nativeElement.height)
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
		    this.importTypes = Array.from(imports.keys())
                    this.imports = imports
                })
            })
    }

    updateExports() {
	this.exportsService.getExports()
            .subscribe((exports: Map<string, BaseElement[]>) => {
                this.zone.run(() => {
		    this.exportTypes = Array.from(exports.keys())
                    this.exports = exports
                })
            })
    }

}
