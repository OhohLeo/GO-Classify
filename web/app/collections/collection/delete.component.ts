import { Component } from '@angular/core'
import { ApiService } from '../../api.service'
import { CollectionsService } from '../collections.service'
import { Collection } from '../collection'

@Component({
    selector: 'collection-delete',
    templateUrl: './delete.component.html',
})

export class DeleteCollectionComponent {

    public title: string

    constructor(private apiService: ApiService,
        private collectionsService: CollectionsService) {

        // Set the name of the collection to delete
        this.title = apiService.collectionSelected.name
    }

    onDelete() {

        // Delete the collection
        this.collectionsService.deleteCollection(this.title)
            .subscribe(status => {
            })
    }
}
