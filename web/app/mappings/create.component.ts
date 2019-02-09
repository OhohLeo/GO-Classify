import { Component, NgZone, Input } from '@angular/core'
import { MappingType, Mapping } from './mapping'

declare var jQuery: any

@Component({
    selector: 'mappings-create',
    templateUrl: './create.component.html'
})

export class CreateComponent {

    @Input() mappingType: MappingType
    @Input() refs: Array<string>
    
    public mappingRef = MappingType
    public mapping = new(Mapping)
    
    private action: any

    onNewMapping() {
        if (this.action != undefined) {
            return
        }

	console.log("NEW MAPPING!")
	
        this.action = jQuery('div#mapping-' + this.mappingType).modal({
            complete: () => {
                this.stop()
            }
        })

        this.action.modal("open")
	
    }

    onSubmit() {
	console.log("ON SUBMIT")
    }
    
    stop() {
        if (this.action == undefined) {
            return
        }

        this.action.modal('close')
        this.action = undefined
    }    
}
