export class DataReference {

    // Handle reference with following format :
    // { "attribute": "format", ...}
    public attributes: Map<string, AttributeReference> = new Map<string, AttributeReference>()

    constructor(public typeRef: string, public name: string, src: any) {

	if (src == undefined) {
	    throw new Error("DataReference '" + name + "' found no data")
	}

	for (let attributeName of Object.keys(src)) {

	    let format = src[attributeName]
	    if (format == undefined) {
		throw new Error("DataReference found no data with '" + format + "'")
	    }

	    this.attributes.set(attributeName, new AttributeReference(typeRef, name, attributeName, format))
	}
    }

    getAttributeNames(): Array<string> {
	return Array.from(this.attributes.keys()).sort()
    }

    getAttributes(): Array<AttributeReference> {
	return Array.from(this.attributes.values())
    }
}

export class DataValues {

    public defaultValueByFormat: Map<string, any> = new Map<string, any>([
	["string", ""],
	["number", 0],
	["date", Date.now()],
    ])
    
    public values: Map<string, AttributeValue> = new Map<string, AttributeValue>()
    constructor(public ref: DataReference, src: any) {
	for (let attributeReference of ref.getAttributes()) {
	    let value

	    // No source defined : set default values
	    if (src == undefined) {
		value = this.defaultValueByFormat[attributeReference.format]
	    }

	    this.values.set(
		attributeReference.name,
		new AttributeValue(attributeReference, value))
	}
    }

    getAttribute(name: string): AttributeValue {
	return undefined
    }
}

export class AttributeReference {
    constructor(public typeRef: string, public dataName: string, public name: string, public format: string) {}

    getRefName() : string {
	return this.typeRef.split("/")[0]
    }

    getTypeName() : string {
	return this.typeRef.split("/")[1]
    }
    
    getLinkName(elementName: string, index: number) {
	return [this.typeRef, elementName, this.dataName, index, this.name].join("/")
    }
}

export class AttributeValue {
    constructor(public ref: AttributeReference, public value: any) {}
}
