import { ConfigRef } from './config_ref'

export class ConfigBase {

    private refs: ConfigRef[] = []
    private refsByName: { [name: string]: ConfigRef } = {}

    private cache: boolean
    private onChange: { [name: string]: (newData: any) => void } = {}

    setRefs(refs: any): boolean {

        if (Array.isArray(refs) == false) {
            console.error("Expected config refs array, received", refs)
            return false
        }

        for (let idx in refs) {
            let ref = new ConfigRef(refs[idx])
            this.refs.push(ref)
            this.refsByName[ref.name] = ref
        }

        return true
    }

    getRefs(): ConfigRef[] {
        return this.refs
    }

    getRef(name: string): ConfigRef {
        return this.refsByName[name]
    }

    hasCache(): boolean {
        return this.cache
    }

    enableCache() {
        this.cache = true
    }

    disableCache() {
        this.cache = false
    }

    getData(): any {

        let data = {}

        for (let ref of this.refs) {

            switch (ref.type) {
                case "map":
                    console.log("GET DATA MAP", ref.name, ref)
                case "struct":
                    data[ref.name] = ref.getData()
                    break;

                default:
                    data[ref.name] = ref.data
            }
        }

        return data
    }

    setData(data: any): boolean {

        if (typeof data != "object") {
            console.error("Expected config object, received", data)
            return false
        }

        for (let [key, value] of Object.entries(data)) {

	    if (value == undefined) {
		continue
	    }
	    
            let ref = this.refsByName[key]
            if (ref === undefined) {
                console.error("Config value '" + key + "' not found")
                continue
            }

	    console.log(key, value)
            ref.setData(value)
        }

        return true
    }
}
