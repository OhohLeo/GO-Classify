import {Component} from 'angular2/core';
import {CreateCollectionComponent} from '../collections/create.component';
import {ModifyCollectionComponent} from '../collections/modify.component';
import {DeleteCollectionComponent} from '../collections/delete.component';

@Component({
    selector: 'config',
    templateUrl: 'app/config/config.html',
	directives: [CreateCollectionComponent,
				 ModifyCollectionComponent,
				 DeleteCollectionComponent]
})

export class ConfigComponent {

}
