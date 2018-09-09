import { Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core'
import { ApiService } from '../../../api.service'
import { Item } from '../../../items/item'

@Component({
    selector: 'simple-icon-item',
    templateUrl: './iconItem.component.html'
})

export class SimpleIconItemComponent implements OnInit, OnDestroy {

    @Input() item: Item
    @Output() open: EventEmitter<Item> = new EventEmitter<Item>()

    public iconUrl: string

    constructor(private apiService: ApiService) {
	console.log("SimpleIconItem")
    }

    ngOnInit() {

	// Get icon url
	this.iconUrl = this.apiService.getIconUrl(this.item)
    }

    ngOnDestroy() { }

    onOpen() {
	this.open.emit(this.item)
    }
}
