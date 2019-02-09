import {
    Component, NgZone, Input, ViewChild, Renderer
} from '@angular/core'

import { ImportsService } from './imports.service'

@Component({
    selector: 'imports-create',
    templateUrl: './create.component.html'
})

export class ImportsCreateComponent {

    @Input() currentRef: string
    
    constructor(private zone: NgZone,
		private render: Renderer,
		private importsService: ImportsService) {
    }

    onCreated(importCreated) {
        console.log("CREATED:",importCreated)
    }
}
