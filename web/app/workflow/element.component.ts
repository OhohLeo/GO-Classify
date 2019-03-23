import { Component, NgZone, Input, OnInit } from '@angular/core'
import { Reference } from '../references/reference'
import { DataReference } from '../references/data'
import { BaseElement } from '../base'

@Component({
    selector: 'workflow-element',
    template: `
<div [ngSwitch]="reference.typ">
  <imports-display *ngSwitchCase="'imports'" [element]="element"></imports-display>
  <exports-display *ngSwitchCase="'exports'" [element]="element"></exports-display>
</div>
<workflow-data *ngFor="let data of datas"
               [data]="data"
               [element]="element">
</workflow-data>
`,
    styles: [
`div {
    background-color: red;
    color: white;
}`]
})

export class ElementComponent implements OnInit {

    @Input() reference: Reference
    @Input() element: BaseElement

    public datas: Array<DataReference> = []
    
    constructor(private zone: NgZone) {}
    
    ngOnInit() {
	this.datas = this.reference.getDataReferences()
	console.log("ELEMENT!", this.element, this.datas)
    }
}
