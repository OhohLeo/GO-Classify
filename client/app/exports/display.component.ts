import { Component, Input } from '@angular/core'

@Component({
    selector: 'exports-display',
    templateUrl: './display.component.html'
})

export class ExportsDisplayComponent {
    @Input() element: string
}
