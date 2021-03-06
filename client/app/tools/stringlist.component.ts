import { Component, NgZone, Input, Output, AfterViewChecked, EventEmitter } from '@angular/core';

declare var jQuery: any

export class CfgStringList {

    public stringlist = new Array<string>()

    constructor(values?: string[]) {
		this.init(values)
    }

    init(values: any) {

        if (values == undefined) {
            this.stringlist = new Array<string>()
            return
        }

        // S'agit-il d'un tableau
        if ((values instanceof Array) == false) {
            console.error("CfgStringList expect 'Array' has init")
            return
        }

        values.forEach((item: string) => {
            this.stringlist.push(item)
        })
    }

    add(name: string): boolean {

        // Check that the name doesn't already exist
        if (this.hasName(name)) {
            console.log("String '" + name + "'already existing")
            return false
        }

        this.stringlist.push(name)
        return true
    }

    remove(name: string): boolean {

        // Check if the name exist
        let index = this.stringlist.indexOf(name)
        if (index < 0) {
            console.log("String '" + name + "'not existing")
            return false
        }

        this.stringlist.splice(index, 1)
        return true
    }

    hasName(search: string): boolean {

        this.stringlist.forEach((name) => {
            if (name === search)
                return true;
        })

        return false
    }

    hasChanged(values: string[]): boolean {

        // S'agit-il d'un tableau
        if ((values instanceof Array) == false) {
            console.error("hasChanged expect 'Array' has parameter")
            return false;
        }

        values.forEach((name) => {
            if (this.hasName(name) == false)
                return true;
        })

        return false;
    }

    getTags(): string[] {
        let result = []

        this.stringlist.forEach((item) => {
            result.push({ 'tag': item })
        })

        return result
    }
}

export class StringListEvent {

    constructor(public name: string,
				public list: string[]) { }
}

@Component({
    selector: 'stringlist',
    template: `<div [attr.id]="name" class="chips chips-placeholder"></div>`
})

export class StringListComponent implements AfterViewChecked {

    @Input() name: string
    @Input() values

    @Output() change = new EventEmitter<StringListEvent>()

	private cfg = new CfgStringList()
    private chip: any = undefined

    constructor(private zone: NgZone) { }

    ngAfterViewChecked() {

        if (this.chip != undefined)
            return

		this.cfg = (this.values instanceof CfgStringList)
			? this.values : new CfgStringList(this.values)
        this.chip = jQuery('div#' + this.name)

        // Initialisation des valeurs
        this.chip.material_chip({
            data: this.cfg.getTags(),
            placeholder: 'Enter a tag',
            secondaryPlaceholder: '+Tag',
        })

        this.chip.on('chip.add', (e, chip) => {

			if (this.cfg.add(chip['tag'])) {
				this.change.emit(new StringListEvent(
					this.name, this.cfg.stringlist))
			}
        })

        this.chip.on('chip.delete', (e, chip) => {

			if (this.cfg.remove(chip['tag']))	{
				this.change.emit(new StringListEvent(
					this.name, this.cfg.stringlist))
			}
        })
    }
}
