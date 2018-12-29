import { Component, Input, OnInit, ViewChild, NgZone, Renderer } from '@angular/core'
import { ConfigsService, ConfigBase, ConfigRef } from './configs.service'
import { ConfigMultiComponent } from './multi.component'

@Component({
    selector: 'configs',
    templateUrl: './configs.component.html'
})

export class ConfigsComponent implements OnInit {

    @Input() src: string
    @Input() name: string
    @Input() init: boolean
    
    @ViewChild(ConfigMultiComponent) multi;

    public collections: string[] = []

    public refMulti: ConfigRef

    public validate: boolean = false

    constructor(private zone: NgZone,
        private render: Renderer,
        private configsService: ConfigsService) { }

    ngOnInit() {
	if (this.init) {
            this.update()
	}
    }

    forceInit(src: string, name: string) {
	this.src = src
	this.name = name
	this.update()
    }
    
    update() {

        this.configsService.getConfigs(this.src, this.name)
            .subscribe((cfg: ConfigBase) => {

                let refs = cfg.getRefs();
                let refsMainList: string[] = []

                for (let idx in refs) {

                    let ref = refs[idx]

                    switch (ref.type) {
                        case "struct":
                        case "map":
                            refsMainList.push(ref.name)
                            break
                    }

                }

                this.zone.run(() => {
                    this.collections = refsMainList
                    this.validate = false
                })
            })
    }


    onMulti(event: any, refSelected: string) {

        // Set collection-items as active
        event.preventDefault()

        for (let item of event.target.parentElement.children) {
            this.render.setElementClass(item, "active", false)
        }

        this.render.setElementClass(event.target, "active", true)

        this.configsService.getConfigs(this.src, this.name)
            .subscribe((cfg: ConfigBase) => {

                let ref = cfg.getRef(refSelected)

                this.zone.run(() => {
                    this.multi.onUpdate(ref)
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

        console.log(this.src, this.name)

        this.configsService.setConfig(this.src, this.name)
            .subscribe((res) => {
                this.zone.run(() => {
                    this.validate = false
                })
            })
    }
}
