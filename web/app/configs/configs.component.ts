import { Component, Input, OnInit, NgZone, Renderer } from '@angular/core'
import { ConfigsService, ConfigBase, ConfigRef } from './configs.service'

@Component({
    selector: 'configs',
    templateUrl: './configs.component.html'
})

export class ConfigsComponent implements OnInit {

    @Input() src: string
	@Input() name: string

    public mainSelector: string[] = []
	public refs: ConfigRef[] = []

	public validate: boolean = false

    constructor(private zone: NgZone,
				private render: Renderer,
				private configsService: ConfigsService) { }

    ngOnInit() {
        this.update()
    }

    update() {

        this.configsService.getConfigs(this.src, this.name)
            .subscribe((cfg: ConfigBase) => {

				let refs = cfg.getRefs();
				let refsMainList : string[] = []

				for (let idx in refs) {

					let ref = refs[idx]

					switch (ref.type)
					{
					case "struct":
						refsMainList.push(ref.name)
						break
					}

				}

                this.zone.run(() => {
					this.mainSelector = refsMainList
					this.validate = false
                })
            })
    }


    onRef(event: any, refSelected: string) {

        // Set collection-items as active
        event.preventDefault()

        for (let item of event.target.parentElement.children) {
            this.render.setElementClass(item, "active", false)
        }

        this.render.setElementClass(event.target, "active", true)

		this.configsService.getConfigs(this.src, this.name)
            .subscribe((cfg: ConfigBase) => {
				let ref = cfg.getRef(refSelected)
				console.log(cfg.getData())
				this.zone.run(() => {
					this.refs = ref.childs
					this.validate = false
                })
			})
	}

    onChange(event) {
		this.zone.run(() => {
			this.validate = true
		})
    }

	onSubmit(event) {

        event.preventDefault()

		this.configsService.setConfig(this.src, this.name)
			.subscribe((res) => {
				this.zone.run(() => {
					this.validate = false
                })
			})
	}
}
