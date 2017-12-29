import { Component, Input } from '@angular/core'
import { Imap } from './imap'

@Component({
    selector: 'imap-display',
    templateUrl: './display.component.html'
})

export class ImapDisplayComponent {

    @Input() imap: Imap
}
