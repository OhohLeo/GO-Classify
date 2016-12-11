import { Component } from '@angular/core';
import { ClassifyService } from '../classify.service';
import { Collection } from './collection';

@Component({
    selector: 'collection-create',
    templateUrl: './create.component.html',
})

export class CreateCollectionComponent {

    public collection: Collection = new Collection('', '')
    private types: string[]
    private websites: string[]

    constructor(private classifySercice: ClassifyService) {

        classifySercice.getReferences()
            .subscribe(
            references => {
                this.websites = references["websites"]
                this.types = references["types"]
            });
    }

    onSubmit() {

        // Check that a collection with same name doesn't already
        // exist

        this.collection.websites = this.websites;

        // Create new collection
        this.classifySercice.newCollection(this.collection)
            .subscribe(status => {

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
