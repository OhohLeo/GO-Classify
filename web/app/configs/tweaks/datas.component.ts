import {
    Component, NgZone, Input, Output, EventEmitter
} from '@angular/core'

import { Tweaks, Tweak } from './tweak'

@Component({
    selector: 'tweaks-datas',
    templateUrl: './datas.component.html',
})

export class TweaksDatasComponent {

	@Output() update = new EventEmitter()

	public name : string
	public refType: string
	public tweaks: Tweaks
	public datas: Tweak[] = []

	constructor(private zone: NgZone) { }

	start(tweaks: Tweaks) {

		this.datas = []
		this.tweaks = tweaks

		tweaks.datas.forEach((value) => {
			this.datas.push(value)
		});

		console.log("[TWEAKS DATA]", tweaks)
		this.zone.run(() => {
			this.name = tweaks.name
			this.refType = tweaks.refType
			this.datas.sort(this.compare);
		})
	}

	compare(a: Tweak, b: Tweak) {
		return a.compare(b)
	}

	onUpdate() {
		this.update.emit()
	}

}
