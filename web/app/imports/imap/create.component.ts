import { Component, NgZone } from '@angular/core'
import { ImportsService } from './../imports.service'
import { Imap } from './imap'

@Component({
    selector: 'imap-create',
    templateUrl: './create.component.html'
})

export class ImapCreateComponent {

    public imap: Imap
    public mailboxes: string[] = []

    constructor(private zone: NgZone,
        private importsService: ImportsService) {

        this.imap = new Imap("")
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

    onSuccess(imap: Imap) {
        this.zone.run(() => {
            this.imap = new Imap("")
            this.mailboxes = []
        })
    }

    // Create new import collection
    onSubmit() {
        this.importsService.addImport(
            this.imap,
            (params) => { return this.onParams(params) },
            (imap) => { return this.onSuccess(imap) })
    }
}
