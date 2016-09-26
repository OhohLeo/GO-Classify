import {Component, OnInit} from '@angular/core';
import {ImportDirectoryComponent} from './directory.component';
import {ImportsService} from './imports.service';

@Component({
    selector: 'imports',
    template: '<import-directory></import-directory>',
    providers: [ImportsService],
    directives: [ImportDirectoryComponent]
})

export class ImportsComponent implements OnInit {

    constructor (private importsService: ImportsService) {}

    ngOnInit() {
        this.importsService.getImports()
    }
}
