import { ImportBase } from '../imports.service';

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
