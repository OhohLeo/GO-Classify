export type Imports = { [type: string]: { [name: string]: any }; };

export class Collection {

    public imports: Imports
    public websites: string[]

    constructor(public name: string, public type: string) { }

    // Retourne vrai lorsqueJSON.parse( l'élément est rajouté à la liste
    addImport(type: string, name: string, params: any): boolean {

        if (this.imports == undefined) {
            return false
        }

        if (this.imports[type] == undefined) {
            this.imports[type] = {}
        }

        if (this.imports[type][name]) {
            return false
        }

        this.imports[type][name] = params
        return true
    }

    // Retourne vrai lorsque l'élément est supprimé de la liste
    deleteImport(type: string, name: string): boolean {

        if (this.imports == undefined
            || this.imports[type] == undefined
            || this.imports[type][name] == undefined) {
            return false
        }

        delete this.imports[type]
        return true
    }
}
