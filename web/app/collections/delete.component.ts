import {Component} from '@angular/core';
import {ClassifyService} from '../classify.service';
import {Collection} from './collection';

@Component({
    selector: 'collection-delete',
    templateUrl: 'app/collections/delete.component.html',
})

export class DeleteCollectionComponent {

    public title: string

    constructor (private classifySercice: ClassifyService) {

        // Set the name of the collection to delete
        this.title = classifySercice.collectionSelected.name
    }

    onDelete() {

        // Delete the collection
        this.classifySercice.deleteCollection(this.title)
            .subscribe(status => {
            })
    }
}
