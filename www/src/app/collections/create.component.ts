import {Component} from '@angular/core';
import {ClassifyService} from '../classify.service';
import {Collection} from './collection';

@Component({
    selector: 'collection-create',
    templateUrl: 'app/collections/create.component.html',
})

export class CreateCollectionComponent {

    public collection: Collection = new Collection('', '')
    private types: string[]
    private websites: string[]

    constructor (private classifySercice: ClassifyService) {

        classifySercice.getReferences()
            .subscribe(
                references => {
                    this.websites = references["websites"]
                    this.types = references["types"]
                    console.log(this.websites, this.types)
                });
    }

    onSubmit() {

        // Check that a collection with same name doesn't already
        // exist

        console.log("CREATE", this.collection)

        // Create new collection
        this.classifySercice.newCollection(this.collection)
            .subscribe(status => {
                console.log(status)

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
