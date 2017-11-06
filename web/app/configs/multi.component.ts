import { Component, Output, EventEmitter, OnInit, NgZone , Renderer} from '@angular/core'
import { ConfigRef } from './configs.service'

@Component({
    selector: 'config-multi',
    templateUrl: './multi.component.html'
})

export class ConfigMultiComponent implements OnInit {

	@Output() update = new EventEmitter<ConfigRef[]>()

	public collections: string[] = []

	public multiRef: { [name: string]: ConfigRef } = {}
	public refs: ConfigRef[] = []

	constructor(private zone: NgZone,
			   	private render: Renderer) { }

    ngOnInit() {

	}

	onUpdate(ref : ConfigRef) {
		console.log("MULTI", ref)

		let collections : string[] = []

		let childs: ConfigRef[] = []

		switch (ref.type)
		{
		case "map":
			for (let idx in ref.childs) {
				let refChild = ref.childs[idx]
				collections.push(refChild.name)
			}
			break;
		case "struct":
			childs = ref.childs
			break;
		}

		this.zone.run(() => {
			this.collections = collections
			this.refs = childs
		})
	}

	onChange(event) {
		this.update.emit(this.refs)
	}


    onRef(event: any, refSelected: string) {

        // Set collection-items as active
        event.preventDefault()

        for (let item of event.target.parentElement.children) {
            this.render.setElementClass(item, "active", false)
        }

        this.render.setElementClass(event.target, "active", true)

		console.log(refSelected)
	}
}
