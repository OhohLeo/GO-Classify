import { Component } from '@angular/core';
import { ImportsService, ImportInterface } from './imports.service';

@Component({
    selector: 'imports',
    template: `<list></list>
<directory></directory>`
})

export class ImportsComponent {

    constructor(private importsService: ImportsService) {}
}
