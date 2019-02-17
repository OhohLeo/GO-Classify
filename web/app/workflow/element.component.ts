import { Component, NgZone, Input, OnInit } from '@angular/core'

import { BaseElement } from '../base'

declare var jQuery: any

@Component({
    selector: 'workflow-element',
    templateUrl: './element.component.html',
    styleUrls: [ './element.component.css' ]
})

export class ElementComponent implements OnInit {

    @Input() element: BaseElement

    constructor(private zone: NgZone) {}
    
    ngOnInit() {
	console.log("ELEMENT!", this.element)
    }
}
