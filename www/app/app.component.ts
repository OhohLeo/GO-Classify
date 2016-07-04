import {Component, OnInit, ViewChild} from 'angular2/core';
import {ClassifyService} from './classify.service';
import {Collection} from './collections/collection';
import {CollectionsComponent} from './collections/collections.component';
import {ImportDirectoryComponent} from './imports/directory.component';

declare var jQuery:any;

enum AppStatus {
    HOME = 1,
    IMPORT,
    EXPORT,
    CONFIG,
}

@Component({
    selector: 'classify',
    templateUrl: 'app/app.html',
    providers: [ClassifyService],
    directives: [CollectionsComponent,
                 ImportDirectoryComponent]
})

export class AppComponent implements OnInit {
    @ViewChild(CollectionsComponent) collections: CollectionsComponent

    public status = AppStatus
    public display = AppStatus.HOME
    public title = "Classify"

    public collection: Collection

    constructor (private classifyService: ClassifyService) {
    }

    ngOnInit() {
        jQuery(".button-collapse").sideNav();
    }

    onHome() {
        this.resetCollectionState()
        this.display = AppStatus.HOME
    }

    onImport() {
        this.resetCollectionState()
        this.display = AppStatus.IMPORT
    }

    onExport() {
        this.resetCollectionState()
        this.display = AppStatus.EXPORT
    }

    onConfig() {
        this.resetCollectionState()
        this.display = AppStatus.CONFIG
    }

    onNewCollection() {
        this.collections.onNewCollection()
    }

    onSelectCollection() {
        this.display = AppStatus.HOME
    }

    onCollectionChoosed(collection: Collection) {
        this.display = AppStatus.HOME
        console.log("COLLECTION", collection)
    }

    resetCollectionState() {
        this.collections.resetCollectionState()
    }
}
