import { Component, OnInit } from '@angular/core';
import { ImportsService } from './imports.service';

@Component({
    selector: 'imports',
    template: '<directory></directory>'
})

export class ImportsComponent implements OnInit {

    constructor(private importsService: ImportsService) { }

    ngOnInit() {
        this.importsService.getImportsConfig()
        this.importsService.getImports()
    }
}
