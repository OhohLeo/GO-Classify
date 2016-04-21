import {Injectable} from 'angular2/core';

@Injectable()
export class CollectionsService {
	getCollections() {
		return Promise.resolve([]);
	}
}
