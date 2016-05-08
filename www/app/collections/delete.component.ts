import {Component} from 'angular2/core';
import {ClassifyService} from '../classify.service';
import {Collection} from './collection';

@Component({
    selector: 'collection-delete',
    templateUrl: 'app/collections/delete.html',
})

export class DeleteCollectionComponent {

    public title: string

    constructor (private classifySercice: ClassifyService) {

        // Set the name of the collection to delete
        this.title = classifySercice.collectionSelected.name
    }

    onDelete() {

        console.log("DELETE", this.title)

        // Delete the collection
        this.classifySercice.deleteCollection(this.title)
            .subscribe(status => {
            })
    }
}
