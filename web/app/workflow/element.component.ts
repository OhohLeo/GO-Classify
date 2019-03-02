import { Component, NgZone, Input, OnInit } from '@angular/core'

import { BaseElement } from '../base'

declare var jQuery: any

@Component({
    selector: 'workflow-element',
    template: `
<div [ngSwitch]="ref">
  <imports-display *ngSwitchCase="'import'" [element]="element"></imports-display>
  <exports-display *ngSwitchCase="'export'" [element]="element"></exports-display>
</div>
`,
    styles: [
`div {
    background-color: red;
    color: white;
}`]
})

export class ElementComponent implements OnInit {

    @Input() ref: string
    @Input() element: BaseElement

    constructor(private zone: NgZone) {}
    
    ngOnInit() {
	console.log("ELEMENT!", this.element)
    }
}
