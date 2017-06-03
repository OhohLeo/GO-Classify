import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core';
import { FileInfo } from './file'

@Component({
    selector: 'detail-file',
    templateUrl: './detail.component.html'
})

export class DetailFileComponent implements OnInit, OnDestroy {

    @Input() file: FileInfo

    private infos: string[] = []
    private needDetails: boolean = true

    constructor(private zone: NgZone) {
    }

    ngOnInit() {
        this.infos.push("fullpath")
    }

    ngOnDestroy() {
    }

    getDetails() {
        this.zone.run(() => {
            this.needDetails = !this.needDetails;
        })
    }


}
