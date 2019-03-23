import { Injectable } from '@angular/core';
import { Response } from '@angular/http';

import { Observable } from 'rxjs/Rx';

import { ApiService, Event } from '../api.service';
import { BufferService } from '../buffer/buffer.service';
import { ReferencesService } from '../references/references.service'

import { BaseElement } from '../base'

@Injectable()
export class WorkflowService {
    
    private addLinkSlotCallback: Map<string, any> = new Map<string, any>()
    private removeLinkSlotCallback: Map<string, any> = new Map<string, any>()

    constructor(private apiService: ApiService,
		private bufferService: BufferService,
		private referencesService: ReferencesService) { }


    setAddLinkSlotCallback(name: string, callback: any) {
	this.addLinkSlotCallback[name] = callback
    }

    setRemoveLinkSlotCallback(name: string, callback: any) {
	this.removeLinkSlotCallback[name] = callback
    }
    
    addLinkSlot(ref: string, name: string, offsetTop: number) {
	let callback = this.addLinkSlotCallback[ref]
	if (callback ==  undefined) {
	    console.error("addLinkSlotCallback '" + ref + "' not found")
	    return 
	}

	callback(name, offsetTop)
    }

    removeLinkSlot(ref: string, name: string, offsetTop: number) {
	let callback = this.removeLinkSlotCallback[ref]
	if (callback ==  undefined) {
	    console.error("removeLinkSlotCallback '" + ref + "' not found")
	    return 
	}

	callback(name, offsetTop)
    }
}
