export type Imports = { [ref: string]: { [name: string]: any }; };

export class Collection {

    public imports: Imports
    public websites: string[]

    constructor(public name: string, public ref: string) { }

    // Retourne vrai lorsqueJSON.parse( l'élément est rajouté à la liste
    addImport(ref: string, name: string, params: any): boolean {

        if (this.imports == undefined) {
            return false
        }

        if (this.imports[ref] == undefined) {
            this.imports[ref] = {}
        }

        if (this.imports[ref][name]) {
            return false
        }

        this.imports[ref][name] = params
        return true
    }

    // Retourne vrai lorsque l'élément est supprimé de la liste
    deleteImport(ref: string, name: string): boolean {

        if (this.imports == undefined
            || this.imports[ref] == undefined
            || this.imports[ref][name] == undefined) {
            return false
        }

        delete this.imports[ref]
        return true
    }
}
