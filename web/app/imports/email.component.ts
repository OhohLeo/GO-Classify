import { Component } from '@angular/core';
import { ImportsService, ImportBase } from './imports.service';
import { ApiService } from '../api.service';

export class Email extends ImportBase {

	public host: string
	public port: number
	public login: string
	public password: string

    constructor(public id: string) {
        super("email", id);
    }

    getParams(): any {
        return {
            "host": this.host,
            "port": this.port,
			"login": this.login,
			"password": this.password,
        }
    }

	display(): string {
        return this.host + ':' + this.port + '(' + this.login + ')'
    }

    compare(i: Email): boolean {
        return super.compare(i)
            && this.host === i.host
			&& this.port == i.port
			&& this.login == i.login
    }
}

@Component({
    selector: 'email',
    templateUrl: './email.component.html'
})

export class EmailComponent {

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
