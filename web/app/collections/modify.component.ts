import {Component} from '@angular/core';
import {ClassifyService} from '../classify.service';
import {Collection} from './collection';

@Component({
    selector: 'collection-modify',
    templateUrl: 'app/collections/modify.component.html',
})

export class ModifyCollectionComponent {

    public title: string
    public collection: Collection
    private websites: string[]

    constructor (private classifySercice: ClassifyService) {

        this.collection = classifySercice.collectionSelected
        this.title = this.collection.name

        classifySercice.getReferences()
            .subscribe(
                references => {
                    this.websites = references["websites"]
                });
    }

    onSubmit() {

        // Check that the parameters of the collection differ

        // Modify the collection
        this.classifySercice.modifyCollection(this.title, this.collection)
            .subscribe(status => {

                // Reset the collection
                this.collection = new Collection('', '')
            })
    }

}
