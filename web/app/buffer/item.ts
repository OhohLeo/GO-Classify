import { CfgStringList } from '../tools/stringlist.component'

export class BufferItem {

    public id: string
    public name: string
    public type: string

    public image: string

    public cleanedName: string
    public banned: CfgStringList
    public separators: CfgStringList

    public probability: number
    public matchId: string
    public match: BufferItem

    private imports: any[] = []

    public webQuery: string

    private websitesById: { [key: string]: any } = {}
    private websites: any[] = []

    constructor(public data: any) {
        this.id = data.id
        this.name = (data.name != undefined) ? data.name : "<unknown>"
        this.type = data.type

        if (data.cleanedName === undefined)
            return

        this.cleanedName = data.cleanedName
        this.banned = new CfgStringList(data.banned)
        this.separators = new CfgStringList(data.separators)

        this.matchId = data.matchId
        this.match = data.match

        this.probability = data.probability

        if (data.imports != undefined) {
            for (let key in data.imports) {
                data.imports[key].forEach((result: any) => {
                    this.imports.push(result)
                })
            }
        }

        this.webQuery = data.webQuery
        if (data.websites != undefined) {

            for (let websiteId in data.websites) {
                data.websites[websiteId].forEach((result: any) => {

                    if (this.matchId !== undefined
                        && this.matchId === result.id) {
                        this.match = result
                    }

                    if (this.websitesById[result.id] !== undefined)
                        return

                    this.websitesById[result.id] = result
                    this.websites.push(result)
                })
            }
        }
    }

    public getName(): string {
        return (this.cleanedName === "") ? this.name : this.cleanedName;
    }

    public getImports(): any[] {
        return this.imports
    }

    public setMatch(id: string) {
        this.matchId = id
        this.match = this.websitesById[id]
    }

    public getMatch(): any {
        return this.match
    }

    public getWebsite(id: string): any {
        return this.websitesById[id]
    }

    public getWebsites(): any[] {
        return this.websites
    }
}
