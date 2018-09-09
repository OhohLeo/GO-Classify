import { Items } from '../items/items'
import { Item } from '../items/item'
import { BaseElement } from '../base'

import { BufferItem } from '../buffer/item'

export type Imports = { [ref: string]: { [name: string]: any }; };

export class Collection extends BaseElement {

    public enableCache: boolean
    public imports: Imports
    public items: Items
    public bufferItems: BufferItem[] = []
    public websites: string[] = []


    constructor(public name: string, ref: string) {
		super("collection", ref, name)
        this.items = new Items(this)
    }

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


    getItems(): Items {
        return this.items
    }

    addItem(item: any) {
        this.items.addItem(item);
    }

    deleteItem(item: Item) {
        if (this.items.removeItem(item)) {
            this.enableCache = false
        }
    }

    updateItem(item: Item): boolean {

        if (this.items.hasItem(item) == false) {
            this.addItem(item)
            return true
        }

        return this.items.updateItem(item)
    }

    toApi(): any {
        return {
            "name": this.name,
            "ref": this.getRef(),
            // "params": {},
            // "config": {},
        }
    }
}
