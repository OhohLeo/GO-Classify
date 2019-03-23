import { Reference } from '../../references/reference'
import { DataReference, AttributeReference, DataValues } from '../../references/data'

export class Tweaks {

    public refType: string
    public datas: Map<string, Tweak> = new Map<string, Tweak>()

    constructor(public isInput: boolean, public name: string, reference: Reference, src: any) {

	if (reference == undefined) {
	    console.error("[Tweaks] no reference found for " + name)
	    return
	}
	
	this.refType = reference.getTypeRef()

	reference.datas.forEach((value, key) => {
	    this.datas.set(key, new Tweak(
		isInput, key, value, (src != undefined) ? src[key] : undefined))
	})
    }

    getDatas() {

	let datas = {}

	this.datas.forEach((tweak, key) => {
	    datas[key] = tweak.GetAttributes()
	})

	return datas
    }

    getValues() {

	let values = {}

	this.datas.forEach((tweak, key) => {
	    let value = tweak.GetValues()
	    if (value != "") {
		values[key] = value
	    }
	})

	return values
    }
}

export class Tweak {

    public attributes: TweakAttribute[] = []

    constructor(public isInput: boolean, public name: string, data: DataReference, src: any) {
	
	data.getAttributes().forEach((attribute) => {
	    this.attributes.push(new TweakAttribute(
		isInput, attribute.name, attribute.format, undefined))
	})

	this.attributes.sort()
    }

    GetAttributes() {

	let attributes = {}

	this.attributes.forEach((attribute) => {
	    attributes[attribute.name] = {
		"value": attribute.src != "" ? attribute.src : "",
		"format": attribute.typ,
	    }
	})

	return attributes
    }

    GetValues() {

	let values = {}

	this.attributes.forEach((attribute) => {

	    // Ignore empty attribute
	    if (!attribute.src) {
		return
	    }

	    let value

	    if (attribute.isInput) {
		value = { "regexp": attribute.src }
	    } else {
		value = { "value": attribute.src }
	    }

	    values[attribute.name] = value
	})

	return values
    }

    compare(tweak: Tweak): number {
	return this.name.localeCompare(tweak.name)
    }
}

export class TweakAttribute {

    public tag: string

    constructor(public isInput: boolean,
		public name: string,
		public typ: string,
		public src: any) {
	this.tag = this.setTag(isInput, typ)
    }

    setTag(isInput: boolean, typ: string) : string {

	this.isInput = isInput

	if (isInput) {
	    return "input-text"
	}

	switch (typ) {
	case "country":
	case "datetime":
	    // return "selector"
	default:
	    return "input-text"
	}
    }

    compare(attribute: TweakAttribute): number {
	return this.name.localeCompare(attribute.name)
    }
}
