import { Component, Input } from '@angular/core'
import { ApiService } from '../api.service'
import { CollectionsService } from '../collections/collections.service'
import { Collection } from './collection'

@Component({
    selector: 'collection-modify',
    templateUrl: './modify.component.html',
})

export class ModifyCollectionComponent {

    @Input() collection: Collection
    private websites: string[]

    constructor(private apiService: ApiService,
				private collectionsService: CollectionsService) {

        apiService.getReferences()
            .subscribe(
				references => {
					this.websites = references["websites"]
				});
    }

    onSubmit() {

        // Check that the parameters of the collection differ

        // Modify the collection
        this.collectionsService.modifyCollection(this.collection.name, this.collection)
            .subscribe(status => {

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
