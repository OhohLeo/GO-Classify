import {Component, OnInit} from 'angular2/core';
import {ClassifyService} from './classify.service';
import {Collection} from './collections/collection';
import {CreateCollectionComponent} from './collections/create.component';
import {ModifyCollectionComponent} from './collections/modify.component';
import {DeleteCollectionComponent} from './collections/delete.component';

declare var jQuery:any;

enum AppStatus {
    HOME = 1,
    IMPORT,
    EXPORT,
    CONFIG,
}

enum CollectionStatus {
    NONE = 0,
    CREATE,
    CHOOSE,
    MODIFY,
    DELETE
}

@Component({
    selector: 'classify',
    templateUrl: 'app/app.html',
    providers: [ClassifyService],
    directives: [CreateCollectionComponent,
                 ModifyCollectionComponent,
                 DeleteCollectionComponent]
})

export class AppComponent implements OnInit {
    public status = AppStatus
    public display = AppStatus.HOME
    public collectionStatus = CollectionStatus
    public collectionState = CollectionStatus.NONE

    public title = "Classify"
    public collections: Collection[] = []

    constructor (private classifyService: ClassifyService) {

        classifyService.setOnChanges((collection: Collection) => {
            this.onChooseCollection(collection)
        })

        classifyService.getAll().subscribe(
            list => {
                console.log(list)
                this.collections = list
                this.onChooseCollection(undefined)
            });
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
        this.collectionState = CollectionStatus.CREATE
    }

    onChooseCollection(collection: Collection) {

        // If no collection exists : ask to create new one
        if (this.collections.length === 0) {
            this.title = "Classify"
            this.collectionState = CollectionStatus.CREATE
            return
        }

        if (collection !== undefined) {
            this.onSelectCollection(collection)
        } else {
            this.title = "Classify"
        }


        this.collectionState = CollectionStatus.CHOOSE
    }

    onSelectCollection(collection: Collection) {
        this.classifyService.selectCollection(collection)
        this.title = collection.name
        this.collectionState = CollectionStatus.NONE
        this.display = AppStatus.HOME
    }

    onModifyCollection(collection: Collection) {
        this.onSelectCollection(collection)
        this.collectionState = CollectionStatus.MODIFY
    }

    onDeleteCollection(collection: Collection) {
        this.onSelectCollection(collection)
        this.collectionState = CollectionStatus.DELETE
    }

    resetCollectionState() {

        // If no collection exists : ask to create new one
        if (this.collections.length === 0) {
            this.collectionState = CollectionStatus.CREATE
            return
        }

        if (this.classifyService.collectionSelected == undefined) {
            this.collectionState = CollectionStatus.CHOOSE
            return
        }

        this.collectionState = CollectionStatus.NONE
    }
}
