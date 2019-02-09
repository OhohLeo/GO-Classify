import {
    Component, NgZone, OnInit, AfterViewInit, OnDestroy,
    ViewChild, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ApiService, Event } from './../api.service'
import { ConfigModalComponent } from './../configs/modal.component'
import { MappingsService } from './mappings.service'

import { MappingType } from './mapping'

declare var jQuery: any;

@Component({
    selector: 'mappings',
    templateUrl: './mappings.component.html'
})

export class MappingsComponent implements OnInit, OnDestroy {

    public mappingType = MappingType
    public importRefs: Array<string> = []
    public exportRefs: Array<string> = []
    
    constructor(private zone: NgZone,
		private render: Renderer,
		private apiService: ApiService,
		private mappingsService: MappingsService) {

	// Refresh the import/export ref list
        apiService.getReferences()
            .subscribe((references) => {
                this.importRefs = references["imports"]
		this.exportRefs = references["exports"]
            })
    }

    ngOnInit() {
    }

    ngOnDestroy() {
    }
}
