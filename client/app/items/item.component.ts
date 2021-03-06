import { Component, NgZone, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core'
import { ItemService } from './item.service'
import { Item } from './item'

declare var jQuery: any

@Component({
    selector: 'item',
    templateUrl: './item.component.html'
})

export class ItemComponent implements OnInit, OnDestroy {

    @Input() item: Item
    @Output() close = new EventEmitter()

    private action: any

    constructor(private itemService: ItemService) {}

    ngOnInit() {

        this.action = jQuery('div#item').modal({
            complete: () => {
                this.stop()
            }
        })

		this.start()

		console.log(this.item)
    }

    ngOnDestroy() {}

    start() {
        // If has buffer items : open modal
		this.action.modal("open")
    }

    stop() {
		this.close.emit()
    }

    onModify() {
		console.log("Modify")
    }

    onDelete() {
		console.log("Delete")
		this.itemService.deleteItem(this.item)
			.subscribe(status => {})
    }
}
