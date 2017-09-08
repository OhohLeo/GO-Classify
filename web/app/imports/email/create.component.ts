import { Component, NgZone } from '@angular/core'
import { ImportsService } from './../imports.service'
import { Email } from './email'

@Component({
    selector: 'email-create',
    templateUrl: './create.component.html'
})

export class EmailCreateComponent {

    public email : Email
	public mailboxes : string[] = []

    constructor(private zone: NgZone,
				private importsService: ImportsService) {

		this.email = new Email("")
	}

	onParams(params: any) : boolean {

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

	onSuccess(email: Email) {
		this.zone.run(() => {
			this.email = new Email("")
			this.mailboxes = []
		})
	}

    // Create new import collection
    onSubmit() {
        this.importsService.addImport(
			this.email,
			(params) => { return this.onParams(params) },
			(email) => { return this.onSuccess(email) })
    }
}
