import { Component, Input } from '@angular/core'

@Component({
    selector: 'imports-display',
    templateUrl: './display.component.html'
})

export class ImportsDisplayComponent {
    @Input() element: string
}
