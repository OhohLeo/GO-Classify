import {
    Component, NgZone, OnInit, AfterViewInit, OnDestroy,
    ViewChild, Renderer
} from '@angular/core'
import { NgSwitch } from '@angular/common'
import { ApiService, Event } from './../api.service'
import { MappingsService } from './mappings.service'

import { CanvasComponent } from './canvas.component'

import { MappingType } from './mapping'

declare var jQuery: any;

@Component({
    selector: 'mappings',
    templateUrl: './mappings.component.html'
})

export class MappingsComponent implements OnInit, AfterViewInit, OnDestroy {

    @ViewChild("importsLinks") importLinks: CanvasComponent
    @ViewChild("exportsLinks") exportLinks: CanvasComponent
    
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

    ngAfterViewInit() {
	console.log(this.importLinks, this.exportLinks)
    }

    ngOnDestroy() {
    }
}
