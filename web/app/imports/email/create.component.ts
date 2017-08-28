import { Component } from '@angular/core';
import { ImportsService, ImportBase } from './../imports.service';
import { ApiService } from '../../api.service';
import { Email } from './email';

@Component({
    selector: 'email-create',
    templateUrl: './create.component.html'
})

export class EmailCreateComponent {

    public email : Email

    constructor(private importsService: ImportsService,
				private apiService: ApiService) {

		this.email = new Email("")

		// Subscribe to convert received data
		importsService.addConvertToImport("email", this.onConvert)
	}

    // Create new import collection
    onSubmit() {
        this.importsService.addImport(this.email)
    }

	onConvert(id: string, params) : ImportBase {

		let email = new Email(id)

		email.host = params['host']
		email.port = params['port']
		email.login = params['login']
		email.password = params['password']

		return email
	}
}
