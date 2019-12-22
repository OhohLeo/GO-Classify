import { ElementRef, Component, NgZone, OnInit, Input } from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ApiService, Event } from './../api.service'
import { WorkflowType } from './workflow'
import { BaseElement } from '../base'

@Component({
    selector: 'workflow',
    templateUrl: './workflow.component.html',
    styleUrls: ['./workflow.component.css']
})

export class WorkflowComponent implements OnInit {

    @Input() collection: BaseElement
        
    public workflowType = WorkflowType
    public importInstance: BaseElement
    public exportInstance: BaseElement
    
    constructor(private zone: NgZone) {
    }

    ngOnInit() {
    }

    setImportInstance(importInstance: BaseElement) {
	this.zone.run(() => {
	    this.importInstance = importInstance
        })
    }

    setExportInstance(exportInstance: BaseElement) {
	this.zone.run(() => {
	    this.exportInstance = exportInstance
        })
    }
}
