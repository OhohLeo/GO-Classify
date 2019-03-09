export class References {

    // Store references with following format :
    // { "type/ref": Reference, ...}
    referenceByTypeRef: Map<string, Reference> = new Map<string, Reference>()

    // Store all types with following format :
    // { "ref": [type, ...], ...}
    refsByType: Map<string, Array<string>> = new Map<string, Array<string>>()
    
    
    constructor(src: any) {
	if (src == undefined) {
	    throw new Error("no references found")
	}

	for (let typeName of Object.keys(src)) {
	    let typeObj = src[typeName]
	    if (typeObj == undefined) {
		throw new Error("no type reference found for " + typeName)
	    }

	    let refsList = Object.keys(typeObj)
	    this.refsByType[typeName] = refsList
	    
	    for (let refName of refsList) {
		this.setReference(typeName, refName, typeObj[refName])
	    }
	}
    }

    setReference(typ: string, ref: string, src: any): Reference {
	let typeRef = typ + "/" + ref
	if (src == undefined) {
	    throw new Error("no reference found for " + typeRef)
	}

	let reference = new Reference(typ, ref, src)
	this.referenceByTypeRef[typeRef] = reference
	return reference
    }

    getReference(typeRef: string): Reference {
	return this.referenceByTypeRef[typeRef]
    }

    getRefs(typeName: string): Array<string> {
	return this.refsByType[typeName]
    }

    getTypeRefs(typeName: string): Array<string> {
	let typeRefs = new Array<string>()
	for (let refName of this.refsByType[typeName]) {
	    typeRefs.push(typeName + "/" + refName)
	}
	return typeRefs
    }
}

export class Reference {

    // Handle reference with following format :
    // { "datas": { "field": "type", ...}, ...}
    public datas: Map<string, DataReference> = new Map<string, DataReference>()

    constructor(public typ: string, public ref: string, src: any) {
	let typeRef = typ + "/" + ref
	if (src == undefined) {
	    throw new Error("no data reference found for " + typeRef)
	}

	let datas = src["datas"]
	if (datas != undefined) {
	    for (let name of Object.keys(datas)) {
		this.datas.set(name, new DataReference(name, datas[name]))
	    }
	}
    }

    getTypeRef(): string {
	return this.typ + "/" + this.ref
    }
    
}

export class DataReference {

    // Handle reference with following format :
    // { "field": "type", ...}
    fields: Map<string, string> = new Map<string, string>()

    constructor(public name: string, src: any) {

	if (src == undefined) {
	    throw new Error("DataReference '" + name + "' found no data")
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
