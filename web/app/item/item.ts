import { CfgStringList } from '../config/stringlist.component'

export class Item {

    public id: string
    public probability: number
    public name: string
    public image: string

    public cleanedName: string
    public banned: CfgStringList
    public separators: CfgStringList

    public bestMatchId: string
    public bestMatch: Item
    private imports: any[] = []

    public webQuery: string

    private websitesById: { [key: string]: any } = {}
    private websites: any[] = []

    constructor(public data: any) {
        this.id = data.id
        this.name = (data.name != undefined) ? data.name : "<unknown>"

        this.cleanedName = data.cleanedName
        this.banned = new CfgStringList(data.banned)
        this.separators = new CfgStringList(data.separators)

        this.bestMatchId = data.bestMatchId
        this.bestMatch = data.bestMatch

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

                    if (this.bestMatchId !== undefined
                        && this.bestMatchId === result.id) {
                        this.bestMatch = result
                    }

                    if (this.websitesById[result.id] !== undefined)
                        return

                    this.websitesById[result.id] = result
                    this.websites.push(result)
                })
            }
        }

        console.log(this.bestMatch)
    }

    public getName(): string {
        return (this.cleanedName === "") ? this.name : this.cleanedName;
    }

    public getImports(): any[] {
        return this.imports
    }

    public setBestMatch(id: string) {
        this.bestMatchId = id
        this.bestMatch = this.websitesById[id]
    }

    public getBestMatch(): any {
        return this.bestMatch
    }

    public getWebsite(id: string): any {
        return this.websitesById[id]
    }

    public getWebsites(): any[] {
        return this.websites
    }
}
