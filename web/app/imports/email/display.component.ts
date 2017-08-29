import { Component, Input } from '@angular/core'
import { Email } from './email'

@Component({
    selector: 'email-display',
    templateUrl: './display.component.html'
})

export class EmailDisplayComponent {

    @Input() email : Email
}
