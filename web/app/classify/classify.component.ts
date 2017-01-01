import { Component, OnInit, OnDestroy } from '@angular/core';

import { ApiService, CollectionStatus, Event } from '../api.service';

declare var jQuery: any;

@Component({
    selector: 'classify',
    templateUrl: './classify.component.html',
})

export class ClassifyComponent implements OnInit, OnDestroy {

    constructor(private apiService: ApiService) { }

    ngOnInit() {
    }

    startModal() {
        jQuery('div#classify').openModal()
    }

    ngOnDestroy() {
        jQuery('div#classify').closeModal()
    }
}
