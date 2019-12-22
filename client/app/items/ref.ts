import { isNumeric } from 'rxjs/util/isNumeric'
import { isArray } from 'rxjs/util/isArray'
import { isDate } from 'rxjs/util/isDate'
import { Item } from './item'

type RefType = 'bool' | 'string' | 'int' | 'float' | 'date' |
    'struct' | 'map' | 'stringlist' | 'invalid';

export class Ref {

    public ref: RefType

    public lastAttributes: string[] = []
    public byAttributes: { [name: string]: RefType } = {}
    public byRefs: { [ref: string]: string[] } = {}

    constructor(public name: string, data: any) {

        let ref = this.getRef(name, data)
        if (ref == 'invalid') {
            return
        }

        switch (ref) {
            case "struct":
                for (let childName in data) {
                    this.addAttribute(childName,
                        data[childName])
                }
        }
    }

    getRef(name: string, data: any): RefType {

        if (isNumeric(data)) {
            return 'int'
        }

        if ((name == "date") || isDate(data)) {
            return 'date'
        }

        if (isArray(data)) {
            return 'stringlist'
        }

        switch (typeof (data)) {
            case "object":
                return 'struct'
            case "string":
                return 'string'
            default:
                console.error("Unhandled ref type '" + name + "'("
                    + typeof (data) + ")", data)
                return 'invalid'
        }
    }

    hasAttribute(name: string): boolean {
        return (this.byAttributes[name] !== undefined)
    }

    addAttribute(name: string, data: any) {

        let ref = this.getRef(name, data)

        this.byAttributes[name] = ref

        if (this.byRefs[ref] == undefined) {
            this.byRefs[ref] = []
        }

        this.byRefs[ref].push(name)
    }

    getAttributes(filters?: string[]): string[] {

        let attributes: string[] = []

        // check if filtered attributes
        let hasFilter = {}
        if (filters != undefined) {
            for (let name of filters) {
                hasFilter[name] = true
            }
        }

        for (let name in this.byAttributes) {
            let ref = this.byAttributes[name]
            switch (ref) {
                case "struct":
                    continue
            }

            if (filters !== undefined
                && !hasFilter[ref]) {
                continue
            }

            attributes.push(name)
        }

        // Store lastAttributes list
        this.lastAttributes = attributes

        return attributes
    }

    getValues(item: Item, attributes?: string[]): string[] {

        let values: string[] = []

        if (attributes == undefined) {
            attributes = this.lastAttributes
        }

        for (let attribute of attributes) {
            let value: string = ""
            let data = item.data[attribute]
            switch (this.getRef(attribute, data)) {
                case "struct":
                    continue
                case "stringlist":
                    data.join(",")
                    break
                default:
                    value = data
            }

            values.push(value)
        }

        return values
    }
}
