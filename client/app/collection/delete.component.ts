import { Component, Input } from '@angular/core'
import { ApiService } from '../api.service'
import { CollectionsService } from '../collections/collections.service'
import { Collection } from './collection'

@Component({
    selector: 'collection-delete',
    templateUrl: './delete.component.html',
})

export class DeleteCollectionComponent {

    @Input() collection: Collection

    constructor(private collectionsService: CollectionsService) {
    }

    onDelete() {

        // Delete the collection
        this.collectionsService.deleteCollection(this.collection.name)
            .subscribe(status => {
            })
    }
}
