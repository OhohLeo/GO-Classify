import { Component, NgZone, Input, AfterViewChecked} from '@angular/core';

declare var jQuery: any;

@Component({
	selector: 'stringlist',
	template: `<div [attr.id]="name" class="chips chips-placeholder"></div>`
})

export class StringListComponent implements AfterViewChecked {

	@Input() name: string

	private chip: any = undefined

    constructor(private zone: NgZone) {}

	ngAfterViewChecked() {

		if (this.chip == undefined)
		{
			this.chip = jQuery('div#' + this.name)

			// Initialisation des valeurs
			this.chip.material_chip({
				data: [],
				placeholder: 'Enter a tag',
				secondaryPlaceholder: '+Tag',
			})

			this.chip.on('chip.add', function(e, chip){
				console.log("ADD", e, chip);
			})

			this.chip.on('chip.delete', function(e, chip){
				console.log("DELETE", e, chip);
			})
		}
	}
}
