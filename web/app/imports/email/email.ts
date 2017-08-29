import { ImportBase } from '../imports.service';

export class Email extends ImportBase {

	public host: string
	public port: number
	public login: string
	public password: string
	public mailbox: string

    constructor(public id: string) {
        super("email", id);
    }

    getParams(): any {
        return {
            "host": this.host,
            "port": this.port,
			"login": this.login,
			"password": this.password,
			"mailbox": this.mailbox,
        }
    }

    compare(i: Email): boolean {
        return super.compare(i)
            && this.host === i.host
			&& this.port == i.port
			&& this.login == i.login
			&& this.mailbox == i.mailbox
    }
}

export function Convert2Email(id: string, params) : ImportBase {

	let email = new Email(id)

	email.host = params['host']
	email.port = params['port']
	email.login = params['login']
	email.password = params['password']
	email.mailbox = params['mailbox']

	return email
}
