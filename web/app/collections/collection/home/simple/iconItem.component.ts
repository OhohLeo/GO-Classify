import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'
import { ApiService } from '../../../../api.service'
import { Item } from '../../../../items/item'

@Component({
    selector: 'simple-icon-item',
    templateUrl: './iconItem.component.html'
})

export class SimpleIconItemComponent implements OnInit, OnDestroy {

    @Input() item: Item
    public iconUrl: string

    constructor(private apiService: ApiService) {

    }
    
    ngOnInit() {

	// Get icon url
	this.iconUrl = this.apiService.getIconUrl(this.item)
    }
    
    ngOnDestroy() { }
}
