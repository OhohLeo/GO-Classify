import { Component, Input } from '@angular/core'

@Component({
    selector: 'imports-create',
    templateUrl: './create.component.html'
})

export class ImportsCreateComponent {

    @Input() currentRef: string
    
    onCreated(importCreated) {
        console.log("CREATED:", importCreated)
    }
}
