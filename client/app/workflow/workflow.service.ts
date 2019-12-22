import { Injectable } from '@angular/core'

import { ApiService, Event } from '../api.service'
import { ReferencesService } from '../references/references.service'

import { BaseElement } from '../base'

@Injectable()
export class WorkflowService {
       
    constructor(private apiService: ApiService,
		private referencesService: ReferencesService) { }

    addInstance(instance: BaseElement) {}
    removeInstance(instance: BaseElement) {}

    getInstance(instance: BaseElement) {}

    getInstanceConfig(instance: BaseElement) {}
    setInstanceConfig(instance: BaseElement) {}
}
