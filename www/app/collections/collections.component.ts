import {Component, Input} from 'angular2/core';
import {ClassifyService} from '../classify.service';
import {Collection} from '../collections/collection';
import {CreateCollectionComponent} from './create.component';
import {ModifyCollectionComponent} from './modify.component';
import {DeleteCollectionComponent} from './delete.component';

enum CollectionStatus {
    NONE = 0,
    CREATE,
    CHOOSE,
    MODIFY,
    DELETE
}

@Component({
    selector: 'collections',
    templateUrl: 'app/collections/collections.component.html',
    directives: [CreateCollectionComponent,
                 ModifyCollectionComponent,
                 DeleteCollectionComponent]
})

export class CollectionsComponent {
    @Input() title: string

    public collectionStatus = CollectionStatus
    public collectionState = CollectionStatus.NONE
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

    onNewCollection() {
        this.collectionState = CollectionStatus.CREATE
    }

    onChooseCollection(collection: Collection) {

        // If no collection exists : ask to create new one
        if (this.collections.length === 0) {
            this.collectionState = CollectionStatus.CREATE
            return
        }

        if (collection !== undefined) {
            this.onSelectCollection(collection)
        }

        this.collectionState = CollectionStatus.CHOOSE
    }

    onSelectCollection(collection: Collection) {
        this.classifyService.selectCollection(collection)
        this.collectionState = CollectionStatus.NONE
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
