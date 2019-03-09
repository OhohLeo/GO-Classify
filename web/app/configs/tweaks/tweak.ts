import { Reference, DataReference } from '../../references/reference'

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
	    datas[key] = tweak.GetFields()
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

    public fields: TweakField[] = []

    constructor(public isInput: boolean, public name: string, data: DataReference, src: any) {
	
	data.fields.forEach((value, key) => {
	    this.fields.push(new TweakField(
		isInput, key, value, undefined))
	})

	this.fields.sort()
    }

    GetFields() {

	let fields = {}

	this.fields.forEach((field, key) => {
	    fields[field.name] = {
		"value": field.src != "" ? field.src : "",
		"format": field.typ,
	    }
	})

	return fields
    }

    GetValues() {

	let values = {}

	this.fields.forEach((field, key) => {

	    // Ignore empty field
	    if (!field.src) {
		return
	    }

	    let value

	    if (field.isInput) {
		value = { "regexp": field.src }
	    } else {
		value = { "value": field.src }
	    }

	    values[field.name] = value
	})

	return values
    }

    compare(tweak: Tweak): number {
	return this.name.localeCompare(tweak.name)
    }
}

export class TweakField {

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

    compare(field: TweakField): number {
	return this.name.localeCompare(field.name)
    }
}
