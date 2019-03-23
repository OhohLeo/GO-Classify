import { ElementRef, Component, NgZone, Input, OnInit, AfterViewInit } from '@angular/core'
import { AttributeReference, AttributeValue } from '../references/data'
import { WorkflowService } from './workflow.service'
import { BaseElement } from '../base'

@Component({
    selector: 'workflow-attribute',
    template: `
<p>{{ref.name}} - {{ref.format}}</p>
`,
    styles: [`
p {
  text-align: right;
}
`]
})

export class AttributeComponent implements OnInit, AfterViewInit {

    @Input() element: BaseElement
    @Input() ref: AttributeReference
    @Input() index: number
    @Input() value: AttributeValue
        
    constructor(private eltRef: ElementRef,
		private zone: NgZone,
		private workflowService: WorkflowService) {}
    
    ngOnInit() {

    }

    ngAfterViewInit() {
	this.workflowService.addLinkSlot(
	    this.ref.getRefName(),
	    this.ref.getLinkName(this.element.name, this.index),
	    this.eltRef.nativeElement.offsetTop)
    }
}
