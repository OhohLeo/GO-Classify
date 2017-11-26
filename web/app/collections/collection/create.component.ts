import { Component } from '@angular/core';
import { ApiService } from '../../api.service';
import { Collection } from '../collection';

@Component({
    selector: 'collection-create',
    templateUrl: './create.component.html',
})

export class CreateCollectionComponent {

    public collection: Collection = new Collection('', '')
    private collections: string[]
    private websites: string[]

    constructor(private apiSercice: ApiService) {

        apiSercice.getReferences()
            .subscribe(
            references => {
                this.websites = references["websites"]
                this.collections = references["collections"]
            });
    }

    onSubmit(website: string) {

        // Check that a collection with same name doesn't already
        // exist

        this.collection.websites = [website];

        // Create new collection
        this.apiSercice.newCollection(this.collection)
            .subscribe(status => {

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
