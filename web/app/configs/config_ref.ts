import { CfgStringList, StringListEvent } from '../tools/stringlist.component'

export class ConfigRef {

    public name: string
    public type: string
    public comments: string
    public data: any

    public parent: ConfigRef
    public childs: ConfigRef[] = []
    private childsByName: { [name: string]: ConfigRef } = {}

    public callbacks: { [name: string]: (newValue: any) => void } = {}

    constructor(ref: any, parent?: ConfigRef) {

        this.name = ref["name"]
        this.type = ref["type"]
        this.comments = ref["comments"]

        this.parent = parent

        let childs

        // Handle arrays
        if (Array.isArray(ref)) {
            childs = ref
        }
        else {
            childs = ref["childs"]
        }

        if (childs != undefined && Array.isArray(childs)) {
            for (let idx in childs) {
                let ref = new ConfigRef(childs[idx], this)
                this.childs.push(ref)
                this.childsByName[ref.name] = ref
            }
        }
        // Handle map type
        else if (this.type == "map") {

            let map = ref["map"]
            for (let keyIdx in map) {

                let keyRef = new ConfigRef({
                    name: keyIdx,
                    type: "key",
                    childs: map[keyIdx]
                }, this)

                this.childs.push(keyRef)
                this.childsByName[keyRef.name] = keyRef
            }
        }
    }

    subscribeDataChange(name: string, cb: (newValue: any) => void) {
        this.callbacks[name] = cb
    }

    unsubscribeDataChange(name: string) {
        delete this.callbacks[name]
    }

    getFullName(): string {

        if (this.parent != undefined)
            return this.parent.getFullName() + "-" + this.name

        return this.name
    }

    getPathName(): string {

        if (this.parent != undefined)
            return this.parent.getFullName()

        return ""
    }

    getData(): any {

        let data = {}

        for (let ref of this.childs) {

            // console.log(ref.type, ref.name)
	    
            switch (ref.type) {
                case "stringlist":
                    if (ref.data instanceof CfgStringList) {
                        data[ref.name] = ref.data.stringlist
                    } else if (Array.isArray(ref.data)) {
                        data[ref.name] = ref.data
                    } else {
                        data[ref.name] = []
                    }
                    break;
                case "map":
                case "key":
                case "struct":
                    data[ref.name] = ref.getData()
                    break;

                default:
                    data[ref.name] = ref.data
            }
        }
        // console.log(data)
        return data
    }

    setData(data: any) {

        switch (this.type) {
            case "stringlist":
                this.data = Array.isArray(data) ? data : []
                return;

            case "key":
            case "map":
            case "struct":

                if (typeof data != "object") {
                    console.error("struct expect object", this, data)
                    return;
                }

                for (let [key, value] of Object.entries(data)) {

                    let ref = this.childsByName[key]
                    if (ref === undefined) {
                        console.error("Config child value '" + key
                            + "' not found in", this.name)
                        continue
                    }

                    ref.setData(value)
                }

                return
        }

        // When values are different :
        if (data !== this.data) {
            for (let name in this.callbacks) {
                this.callbacks[name](data)
            }
        }

        this.data = data
    }
}
