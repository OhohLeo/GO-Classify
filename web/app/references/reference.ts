export class Reference {

	// Store reference as following format :
	// { "data_ref": { "field": "type", ...}, ...}
	datas: Map<string, DataReference> = new Map<string, DataReference>()

	constructor(public name: string, src: any) {

		if (src == undefined) {
			throw new Error("Reference found no data")
		}

		for (let name of Object.keys(src)) {

			let value = src[name]
			if (value == undefined) {
				throw new Error("Reference found no data with '" + name + "'")
			}

			this.datas.set(name, new DataReference(name, value))
		}
	}
}

export class DataReference {

	// Store reference as following format :
	// { "field": "type", ...}
	fields: Map<string, string> = new Map<string, string>()

	constructor(public name: string, src: any) {

		if (src == undefined) {
			throw new Error("DataReference found no data")
		}

		for (let field of Object.keys(src)) {

			let typ = src[field]
			if (typ == undefined) {
				throw new Error("DataReference found no data with '" + typ + "'")
			}

			this.fields.set(field, typ)
		}
	}
}