import {
    Component, NgZone, Input, ViewChild, Renderer
} from '@angular/core'

import { ExportsService } from './exports.service'

@Component({
    selector: 'exports-create',
    templateUrl: './create.component.html'
})

export class ExportsCreateComponent {

    @Input() currentRef: string
    
    constructor(private zone: NgZone,
		private render: Renderer,
		private exportsService: ExportsService) {
    }

    onCreated(exportCreated) {
        console.log("CREATED:",exportCreated)
    }
}
