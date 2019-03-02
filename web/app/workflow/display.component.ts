import { Component, NgZone, Input, OnInit } from '@angular/core'

import { BaseElement } from '../base'

declare var jQuery: any

@Component({
    selector: 'workflow-display',
    template: `
<div>
  <p>{{elementsRef}}/{{elementsType}}</p>
  <workflow-element *ngFor="let element of elementList"
                    [element]="element"
                    [ref]="elementsRef">
  </workflow-element>
</div>
`,
    styles: [
`div {
  outline-color: red;
  outline-width: 4px;
  outline-style: solid;
}`]
})

export class DisplayComponent  implements OnInit {

    @Input() elementsRef: string
    @Input() elementsType: string
    @Input() elements: Map<string, BaseElement[]>

    public elementList: Array<BaseElement> = []

    constructor(private zone: NgZone) {}
    
    ngOnInit() {
	this.zone.run(() => {
	    this.elementList = this.elements.get(this.elementsType)
	})
    }
}
