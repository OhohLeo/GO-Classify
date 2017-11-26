import { Component } from '@angular/core';
import { ApiService } from '../../api.service';
import { Collection } from '../collection';

@Component({
    selector: 'collection-modify',
    templateUrl: './modify.component.html',
})

export class ModifyCollectionComponent {

    public title: string
    public collection: Collection
    private websites: string[]

    constructor(private apiSercice: ApiService) {

        this.collection = apiSercice.collectionSelected
        this.title = this.collection.name

        apiSercice.getReferences()
            .subscribe(
            references => {
                this.websites = references["websites"]
            });
    }

    onSubmit() {

        // Check that the parameters of the collection differ

        // Modify the collection
        this.apiSercice.modifyCollection(this.title, this.collection)
            .subscribe(status => {

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
