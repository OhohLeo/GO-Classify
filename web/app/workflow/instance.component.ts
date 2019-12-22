import { Component, NgZone, Input, OnInit, Output, EventEmitter } from '@angular/core'
import { WorkflowType, Workflow } from './workflow'
import { ImportsService } from '../imports/imports.service'
import { ExportsService } from '../exports/exports.service'
import { BaseElement } from '../base'
import { References, Reference } from '../references/reference'

declare var jQuery: any

@Component({
    selector: 'workflow-instance',
    templateUrl: './instance.component.html'
})

export class InstanceComponent implements OnInit {
    @Input() instance: BaseElement
    @Input() collection: BaseElement
    
    ngOnInit() {
	console.log(this.instance)
    }
}
