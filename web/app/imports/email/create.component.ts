import { Component, NgZone } from '@angular/core'
import { ImportsService, ImportBase } from './../imports.service'
import { ApiService } from '../../api.service'
import { Email } from './email'

@Component({
    selector: 'email-create',
    templateUrl: './create.component.html'
})

export class EmailCreateComponent {

    public email : Email
	public mailboxes : string[] = []

    constructor(private zone: NgZone,
				private importsService: ImportsService,
				private apiService: ApiService) {

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
