import { Component } from '@angular/core'
import { ApiService } from '../api.service'
import { CollectionsService } from '../collections/collections.service'
import { Collection } from './collection'
import { References } from '../references/reference'

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
            .subscribe((references : References) => {
		this.websites = references.getRefs("websites")
		this.collections = references.getRefs("collections")
	    });
    }

    onSubmit() {

        // Check that a collection with same name doesn't already
        // exist

        // Create new collection
        this.collectionsService.newCollection(this.collection)
            .subscribe(status => {

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
