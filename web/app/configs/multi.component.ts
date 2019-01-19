import {
    Component, Input, Output, EventEmitter, OnInit, NgZone,
    Renderer, ViewChild
} from '@angular/core'
import { ConfigRef } from './config_ref'
import { TweaksComponent } from './tweaks/tweaks.component'
import { BaseElement } from '../base'

@Component({
    selector: 'config-multi',
    templateUrl: './multi.component.html'
})

export class ConfigMultiComponent implements OnInit {

    @Input() item : BaseElement
    @Output() update = new EventEmitter<ConfigRef[]>()

    public tabs: string[] = []
    public structs: string[] = []

    public refsByTab: { [name: string]: ConfigRef[] } = {}
    public refsByStruct: { [name: string]: ConfigRef[] } = {}
    public refs: ConfigRef[] = []

    @ViewChild(TweaksComponent) tweaks
    public hasTweaks: boolean = false

    constructor(private zone: NgZone,
		private render: Renderer) { }

    ngOnInit() { }

    onUpdate(ref: ConfigRef) {

        let tabs: string[] = []

        let childs: ConfigRef[] = []

        switch (ref.type) {
        case "map":
            for (let idx in ref.childs) {
                let refElement = ref.childs[idx]
                tabs.push(refElement.name)
                if (refElement.type === "key") {
                    this.refsByTab[refElement.name] = refElement.childs
                } else {
                    childs.push(refElement)
                }
            }
            break;
        case "struct":
            childs = ref.childs
            break;
	}
    
        this.zone.run(() => {
            this.tabs = tabs
        })

        this.updateRefChilds(childs)
    }

    updateRefChilds(childs: ConfigRef[]) {

        let structs: string[] = []
        let refs: ConfigRef[] = []

        for (let idx in childs) {
            let refChild = childs[idx]
            switch (refChild.type) {
            case "struct":
                structs.push(refChild.name)
                this.refsByStruct[refChild.name] = refChild.childs
                break
	    case "ptr":
		this.hasTweaks = true
		console.log("PTR", refChild.name)
		break
            default:
                refs.push(refChild)
            }
        }

        this.zone.run(() => {
            this.structs = structs
            this.refs = refs
        })
    }

    onChange(ref) {
        this.update.emit(this.refs)
    }

    onRef(event: any, refSelected: string) {

        // Set collection-items as active
        event.preventDefault()

        for (let item of event.target.parentElement.children) {
            this.render.setElementClass(item, "active", false)
        }

        this.render.setElementClass(event.target, "active", true)

        this.updateRefChilds(this.refsByTab[refSelected])
    }
}
