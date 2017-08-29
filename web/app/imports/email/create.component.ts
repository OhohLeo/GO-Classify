import { Component, NgZone } from '@angular/core';
import { ImportsService, ImportBase } from './../imports.service';
import { ApiService } from '../../api.service';
import { Email } from './email';

declare var jQuery: any;

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

		if (typeof(params) === "object") {
			let mailboxes = params["mailboxes"]

			if (Array.isArray(mailboxes) && mailboxes.length > 0) {

				this.zone.run(() => {
					this.mailboxes = mailboxes
					console.log(this.mailboxes)
				})

				this.zone.run(() => {
					jQuery('select').material_select()
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
			(params) => { this.onParams(params) },
			(email) => { this.onSuccess(email) })
    }
}
