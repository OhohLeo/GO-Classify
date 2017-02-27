import { Component, NgZone, Input, OnInit, OnDestroy} from '@angular/core';

@Component({
	selector: 'config',
	templateUrl: './config.component.html'
})

export class ConfigComponent {

    constructor(private zone: NgZone) {}
}
