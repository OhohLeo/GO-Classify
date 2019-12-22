import { Component, NgZone, Input, OnInit, Output, EventEmitter } from '@angular/core'
import { WorkflowType, Workflow } from './workflow'
import { ImportsService } from '../imports/imports.service'
import { ExportsService } from '../exports/exports.service'
import { BaseElement } from '../base'
import { References, Reference } from '../references/reference'

declare var jQuery: any

@Component({
    selector: 'workflow-collection',
    templateUrl: './collection.component.html'
})

export class CollectionComponent implements OnInit {

    @Input() collection: BaseElement
    
    ngOnInit() {
	console.log(this.collection)
    }
}
