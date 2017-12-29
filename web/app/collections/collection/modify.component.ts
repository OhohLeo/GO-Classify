import { Component } from '@angular/core'
import { ApiService } from '../../api.service'
import { CollectionsService } from '../collections.service'
import { Collection } from '../collection'

@Component({
    selector: 'collection-modify',
    templateUrl: './modify.component.html',
})

export class ModifyCollectionComponent {

    public title: string
    public collection: Collection
    private websites: string[]

    constructor(private apiService: ApiService,
        private collectionsService: CollectionsService) {

        this.collection = apiService.collectionSelected
        this.title = this.collection.name

        apiService.getReferences()
            .subscribe(
            references => {
                this.websites = references["websites"]
            });
    }

    onSubmit() {

        // Check that the parameters of the collection differ

        // Modify the collection
        this.collectionsService.modifyCollection(this.title, this.collection)
            .subscribe(status => {

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
