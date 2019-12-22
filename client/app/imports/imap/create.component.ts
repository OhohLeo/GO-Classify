import { Component, NgZone } from '@angular/core'
import { ImportsService } from '../imports.service'
import { BaseCreateComponent } from '../../tools/base_create.component'
import { Imap } from './imap'

@Component({
    selector: 'imap-create',
    templateUrl: './create.component.html'
})

export class ImapCreateComponent extends BaseCreateComponent {

    public mailboxes: string[] = []

    constructor(private zone: NgZone,
		private importsService: ImportsService) {
	super(new Imap(""))
    }

    onParams(params: any): boolean {

        if (params instanceof Object) {

            let mailboxes = params["mailboxes"]

            if (Array.isArray(mailboxes) && mailboxes.length > 0) {

                this.zone.run(() => {
                    this.mailboxes = mailboxes
                })

                return true
            }
        }

        return false
    }

    onSuccess(data: Imap) {
        this.zone.run(() => {
            this.data = new Imap("")
            this.mailboxes = []
        })
    }
}
