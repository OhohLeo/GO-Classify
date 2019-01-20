import {
    Component, NgZone, Input, Output, EventEmitter
} from '@angular/core'

import { Tweaks, Tweak } from './tweak'

@Component({
    selector: 'tweaks-datas',
    templateUrl: './datas.component.html',
})

export class TweaksDatasComponent {

    @Input() tweaks: Tweaks
    @Output() update = new EventEmitter()

    public name : string
    public refType: string
    public datas: Tweak[] = []

    constructor(private zone: NgZone) { }

    ngOnInit() { 	
	let datas = []
	this.tweaks.datas.forEach((value) => {
	    datas.push(value)
	});
	datas.sort(this.compare);

	// console.log("[TWEAKS DATA]", this.tweaks)
	this.name = this.tweaks.name
	this.refType = this.tweaks.refType
	this.datas = datas
    }

    compare(a: Tweak, b: Tweak) {
	return a.compare(b)
    }

    onUpdate() {
	this.update.emit()
    }

}
