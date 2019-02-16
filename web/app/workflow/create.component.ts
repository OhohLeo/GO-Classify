import { Component, NgZone, Input } from '@angular/core'
import { WorkflowType, Workflow } from './workflow'

declare var jQuery: any

@Component({
    selector: 'workflow-create',
    templateUrl: './create.component.html'
})

export class CreateComponent {

    @Input() workflowType: WorkflowType
    @Input() refs: Array<string>
    
    public workflowRef = WorkflowType
    public workflow = new(Workflow)
    
    private action: any

    onNewWorkflow() {
        if (this.action != undefined) {
            return
        }

	console.log("NEW MAPPING!")
	
        this.action = jQuery('div#workflow-' + this.workflowType).modal({
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
