import {
    Component, NgZone, Input, Output, EventEmitter, OnInit, ViewChildren, QueryList
} from '@angular/core'

import { Observable } from 'rxjs/Rx';

import { TweaksService } from './tweaks.service'
import { ApiService, Event } from '../../api.service'

import { TweaksDatasComponent } from './datas.component'

import { Tweaks } from './tweak'
import { BaseElement } from '../../base'

@Component({
    selector: 'tweaks',
    templateUrl: './tweaks.component.html',
})

export class TweaksComponent implements OnInit {

	public needHelp: boolean
	public canValidate: boolean

	private item: BaseElement

	@ViewChildren(TweaksDatasComponent) datas: QueryList<TweaksDatasComponent>;

	constructor(private zone: NgZone,
				private apiService: ApiService,
				private tweakService: TweaksService) {}

    ngOnInit() {}

	start(item: BaseElement) {
		this.item = item

		Observable.combineLatest(
			this.tweakService.getTweak(item.getType(), item.getName()),
			this.tweakService.getReferences(item)
		).subscribe(([ tweak, references ]) => {

			if (tweak == undefined) {
				console.error("[Tweaks] tweak response is invalid")
			}

			if (references == undefined) {
				console.error("[Tweaks] references response is invalid")
			}

			console.log("[TWEAK] INPUT", tweak["input"], references[0])
			console.log("[TWEAK] OUTPUT", tweak["output"], references[1])

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
				this.datas.first.start(
					new Tweaks(true, input, references[0], tweak["input"]))
				this.datas.last.start(
			 		new Tweaks(false, output, references[1], tweak["output"]))
			})
		})
	}

	onHelp() {
		this.zone.run(() => {
			this.needHelp = !this.needHelp
		})
	}

	onUpdate() {
		this.zone.run(() => {
			this.canValidate = true
		})
	}

	onValidate() {

		let input = this.datas.first
		let output = this.datas.last

		let values = {
			"input": input.tweaks.getValues(),
			"output": output.tweaks.getValues()
		}

		console.log("[TWEAK] VALIDATE", values)
		this.tweakService.setTweak(this.item.getType(), this.item.getName(), values)
			.subscribe((rsp) => {
				console.log("[TWEAK SET]", rsp)
			})
	}
}
