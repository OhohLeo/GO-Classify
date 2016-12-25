import { Component } from '@angular/core';
import { ApiService } from '../api.service';
import { Collection } from './collection';

@Component({
    selector: 'collection-delete',
    templateUrl: './delete.component.html',
})

export class DeleteCollectionComponent {

    public title: string

    constructor(private apiSercice: ApiService) {

        // Set the name of the collection to delete
        this.title = apiSercice.collectionSelected.name
    }

    onDelete() {

        // Delete the collection
        this.apiSercice.deleteCollection(this.title)
            .subscribe(status => {
            })
    }
}
