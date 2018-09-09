import { Injectable } from '@angular/core';

import { Reference } from './reference'
import { BaseElement } from '../base'

@Injectable()
export class ReferencesService {

	// Store refences as following format :
	// { "ref/typ": { "data_ref": { "field": "type", ...}, ...}, ...}
	references: Map<string, Reference> = new Map<string, Reference>()

	constructor() { }

	setReferences(b: BaseElement, src: any) : Reference {

		if (src == undefined) {
			console.error(
				"References can't handle empty source from '"
					+ b.getTypeRef() + " (" +  b.getName() + ")'")
			return undefined
		}

		try {
			let reference = new Reference(b.getTypeRef(), src["datas"])
			this.references.set(b.getTypeRef(),  reference)
			return reference
		}
		catch (e) {
			console.error(
				"References can't handle source from '" + b.getTypeRef() + "': ", e.message)
			return undefined
		}
	}

	getReferences(b: BaseElement) : Reference {
		return this.references.get(b.getTypeRef())
	}
}
