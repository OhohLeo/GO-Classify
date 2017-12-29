import { Component } from '@angular/core'
import { ApiService } from '../../api.service'
import { CollectionsService } from '../collections.service'
import { Collection } from '../collection'

@Component({
    selector: 'collection-create',
    templateUrl: './create.component.html',
})

export class CreateCollectionComponent {

    public collection: Collection = new Collection('', '')
    private collections: string[]
    private websites: string[]

    constructor(private apiService: ApiService,
        private collectionsService: CollectionsService) {

        apiService.getReferences()
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
        this.collectionsService.newCollection(this.collection)
            .subscribe(status => {

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
