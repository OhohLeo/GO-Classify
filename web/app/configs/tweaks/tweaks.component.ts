import {
    Component, NgZone, Input, Output, EventEmitter, OnInit
} from '@angular/core'

import { Observable } from 'rxjs/Rx';

import { TweaksService } from './tweaks.service'
import { ApiService, Event } from '../../api.service'

import { Tweaks } from './tweak'
import { ConfigRef } from '../config_ref'
import { BaseElement } from '../../base'

@Component({
    selector: 'tweaks',
    templateUrl: './tweaks.component.html',
})

export class TweaksComponent implements OnInit {

    @Input() item : BaseElement
    @Input() ref: ConfigRef
    @Output() update = new EventEmitter()

    public needHelp: boolean
    public input : Tweaks
    public output : Tweaks
    
    constructor(private zone: NgZone,
		private apiService: ApiService,
		private tweakService: TweaksService) {}

    ngOnInit() {
	this.start({"input": null, "output": null})
    }

    start(tweak) {
	let item = this.item
	this.tweakService.getReferences(item).subscribe((references) => {
	    
	    if (tweak == undefined) {
		console.error("[TWEAKS] tweak response is invalid")
		return
	    }

	    if (references == undefined) {
		console.error("[TWEAKS] references response is invalid")
		return
	    }

	    // console.log("[TWEAK] INPUT", tweak["input"], references[0])
	    // console.log("[TWEAK] OUTPUT", tweak["output"], references[1])

	    let input: string
	    let output: string

	    switch (item.getType()) {
	    case "imports":
		input = item.name
		output = this.apiService.getCollectionName()
		break
	    case "exports":
		input = this.apiService.getCollectionName()
		output = item.name
		break
	    default:
		console.error("[TweaksComponent] item not possible on '" + item.getType() + "'")
		return
	    }

	    this.zone.run(() => {
		this.input = new Tweaks(true, input, references[0], tweak["input"])
		this.output = new Tweaks(false, output, references[1], tweak["output"])
	    })
	})
    }

    onHelp() {
	this.zone.run(() => {
	    this.needHelp = !this.needHelp
	})
    }

    onUpdate() {
	this.ref.data = {
	    "input": this.input.getValues(),
	    "output": this.output.getValues()
	}
	this.update.emit(this.ref)
    }
}
