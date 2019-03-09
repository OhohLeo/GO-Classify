import { Injectable } from '@angular/core';

import { References, Reference } from './reference'
import { BaseElement } from '../base'

@Injectable()
export class ReferencesService {

    references: References

    constructor() { }

    setReferences(src: any) : References {
	try {
	    let references = new References(src)
	    this.references = references
	    return references
	}
	catch (e) {
	    console.error("All References can't handle source: ", e.message)
	    return undefined
	}
    }

    getReferences(): References {
	return this.references
    }
    
    setReference(b: BaseElement, src: any) : Reference {

	if (src == undefined) {
	    console.error(
		"References can't handle empty source from '"
		    + b.getTypeRef() + " (" +  b.getName() + ")'")
	    return undefined
	}
	
	try {
	    return this.references.setReference(b.getType(), b.getRef(), src)
	}
	catch (e) {
	    console.error(
		"References can't handle source from '" + b.getTypeRef() + "': ",
		e.message)
	    return undefined
	}
    }

    getReference(b: BaseElement) : Reference {
	return this.references.getReference(b.getTypeRef())
    }
}
