import {
    Component, NgZone, OnInit, AfterViewInit, OnDestroy,
    ViewChild, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ApiService, Event } from './../api.service'
import { WorkflowService } from './workflow.service'

import { CanvasComponent } from './canvas.component'

import { WorkflowType } from './workflow'

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
    public exportRefs: Array<string> = []
    
    constructor(private zone: NgZone,
		private render: Renderer,
		private apiService: ApiService,
		private workflowService: WorkflowService) {

	// Refresh the import/export ref list
        apiService.getReferences()
            .subscribe((references) => {
                this.importRefs = references["imports"]
		this.exportRefs = references["exports"]
            })
    }

    ngOnInit() {
    }

    ngAfterViewInit() {
	console.log(this.importLinks, this.exportLinks)
    }

    ngOnDestroy() {
    }
}
