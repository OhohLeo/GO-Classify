import { Reference, DataReference } from '../../references/reference'

export class Tweaks {

	public refType: string
	public datas: Map<string, Tweak> = new Map<string, Tweak>()

	constructor(public isInput: boolean, public name: string, reference: Reference, src: any) {

		this.refType = reference.name

		reference.datas.forEach((value, key) => {
			console.log(value, key)
			this.datas.set(key, new Tweak(
				isInput, key, value, (src != undefined) ? src[key] : undefined))
		})
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

	compare(tweak: Tweak): number {
		return this.name.localeCompare(tweak.name)
	}
}

export class TweakField {

	public tag: string

	constructor(isInput: boolean,
				public name: string,
				public typ: string,
				src: any) {
		this.tag = this.setTag(isInput, typ)
		console.log(name, typ, this.tag)
	}

	setTag(isInput: boolean, typ: string) : string {

		if (isInput) {
			return "input-text"
		}

		switch (typ) {
		case "country":
		case "datetime":
			return "selector"
		default:
			return "input-text"
		}
	}

	compare(field: TweakField): number {
		return this.name.localeCompare(field.name)
	}
}
