import { Component, NgZone, Input, Output, OnInit, EventEmitter } from '@angular/core'
import { WorkflowType, Workflow } from './workflow'
import { ImportsService } from '../imports/imports.service'
import { ExportsService } from '../exports/exports.service'
import { BaseElement } from '../base'
import { References, Reference } from '../references/reference'

declare var jQuery: any

@Component({
    selector: 'workflow-select',
    templateUrl: './select.component.html'
})

export class SelectComponent implements OnInit {

    @Input() workflowType: WorkflowType
    @Output() eventInstance = new EventEmitter<BaseElement>()

    public workflowRef = WorkflowType
    public workflow = new(Workflow)

    public instances: BaseElement[] = []
    public instanceByIDs: Map<string, BaseElement> = new Map<string, BaseElement>()
    public currentInstanceID: string
    public selectedInstance: number

    public service: any	
    private action: any

    constructor(private zone: NgZone,
		private importsService: ImportsService,
		private exportsService: ExportsService) {
    }

    ngOnInit() {
	let service, updateFunc
	switch (this.workflowType) {
	case WorkflowType.IMPORT:
	    service = this.importsService
	    updateFunc = () => {
		this.importsService.getImports()
		    .subscribe((imports: Map<string, BaseElement[]>) => {
			this.updateInstances(imports)
		    })
	    }
	    break
	case WorkflowType.EXPORT:
	    service = this.exportsService
	    updateFunc = () => {
		this.exportsService.getExports()
		    .subscribe((exports: Map<string, BaseElement[]>) => {
			this.updateInstances(exports)
		    })
	    }
	    break
	}

	this.service = service
	service.setUpdateList = updateFunc
	updateFunc()
    }    
    
    updateInstances(inputs: Map<string, BaseElement[]>) {
	let instances = []
	let instanceByIDs = new Map<string, BaseElement>()
	for (let instanceList of inputs.values()) {
	    for (let instance of instanceList) {
		instances.push(instance)
		instanceByIDs.set(instance.getID(), instance)
	    }
	}

	this.zone.run(() => {
	    this.instances = instances
	    this.instanceByIDs = instanceByIDs
	})

	if (instances.length > 0) {
	    this.instanceToSelect(instances[0])
	}
    }

    instanceToSelect(searched: BaseElement) {
	let selectedInstance = -1;
	for (let instance of this.instances) {
	    selectedInstance += 1
	    if (instance == searched) {
		break
	    }
	}

	if (selectedInstance < 0) {
	    console.error("instance '" + searched.getID() + "' position not found")
	    return
	}

	this.zone.run(() => {
	    this.selectedInstance = selectedInstance
	})

	this.eventInstance.emit(searched)
    }
    
    onInstanceCreated() {
        if (this.action != undefined) {
            return
        }
	
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

    onInstanceChanged() {
	let currentInstance = this.instanceByIDs.get(this.currentInstanceID)
	if (currentInstance == undefined) {
	    console.error("instance with ID '" + this.currentInstanceID + "' not found")
	    return
	}
	console.log(currentInstance)
	this.eventInstance.emit(currentInstance)
    }
}
