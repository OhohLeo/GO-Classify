import { ImportBase } from '../imports.service';

export class Imap extends ImportBase {

    public host: string
    public port: number
    public login: string
    public password: string
    public mailbox: string
    public onlyAttached: boolean = false
    public search: Search = new Search()

    constructor(public name: string) {
        super("imap", name);
    }

    getParams(): any {
        return {
            "host": this.host,
            "port": this.port,
            "login": this.login,
            "password": this.password,
            "mailbox": this.mailbox,
            "onlyAttached": this.onlyAttached,
            "search": this.search.getParams(),
        }
    }

    compare(i: Imap): boolean {
        return super.compare(i)
            && this.host === i.host
            && this.port == i.port
            && this.login == i.login
            && this.mailbox == i.mailbox
            && this.onlyAttached == i.onlyAttached
            && this.search.compare(i.search)
    }
}

export class Search {

    public since: string
    public before: string
    public larger: number
    public smaller: number
    public text: string[] = []

    getParams(): any {

        console.log(this.since)

        return {
            "since": this.since,
            "before": this.before,
            "larger": this.larger,
            "smaller": this.smaller,
            "text": this.text,
        }
    }

    compare(i: Search): boolean {
        return this.since === i.since
            && this.before == i.before
            && this.larger == i.larger
            && this.smaller == i.smaller
            && this.text == i.text
    }
}

export function Convert2Imap(name: string, params): ImportBase {

    let imap = new Imap(name)

    imap.host = params['host']
    imap.port = params['port']
    imap.login = params['login']
    imap.password = params['password']
    imap.mailbox = params['mailbox']
    imap.onlyAttached = params['onlyAttached']
    imap.search = Convert2Search(params['search'])

    return imap
}

export function Convert2Search(params): Search {

    let search = new Search()

    search.since = params['since']
    search.before = params['before']
    search.larger = params['larger']
    search.smaller = params['smaller']
    search.text = params['text']

    return search
}
