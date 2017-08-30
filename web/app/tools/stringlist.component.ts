import { Component, NgZone, Input, Output, AfterViewChecked, EventEmitter } from '@angular/core';

declare var jQuery: any

export class CfgStringList {

    public stringlist = new Array<string>()

    constructor(values?: string[]) {

        if (values) {
            for (let value of values) {
                this.add(value)
            }
        }
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
        public action: string,
        public list: string[]) { }
}

@Component({
    selector: 'stringlist',
    template: `<div [attr.id]="name" class="chips chips-placeholder"></div>`
})

export class StringListComponent implements AfterViewChecked {

    @Input() name: string
    @Input() value: CfgStringList

    @Output() change = new EventEmitter<StringListEvent>()

    private chip: any = undefined

    constructor(private zone: NgZone) { }

    ngAfterViewChecked() {

        if (this.chip != undefined)
            return

        this.chip = jQuery('div#' + this.name)

		console.log(this.value)

        // Initialisation des valeurs
        this.chip.material_chip({
            data: (this.value instanceof CfgStringList) ?
                this.value.getTags() : [],
            placeholder: 'Enter a tag',
            secondaryPlaceholder: '+Tag',
        })

        this.chip.on('chip.add', (e, chip) => {
            this.change.emit(new StringListEvent(this.name, "add", [chip.tag]))
        })

        this.chip.on('chip.delete', (e, chip) => {
            this.change.emit(new StringListEvent(this.name, "remove", [chip.tag]))
        })
    }
}
