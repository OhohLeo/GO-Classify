import { Injectable } from '@angular/core'
import { Response } from '@angular/http'
import { Observable } from 'rxjs/Rx'
import { ApiService, CollectionStatus } from './../api.service'
import { Item } from '../items/item'

@Injectable()
export class ItemsService {

    constructor(private apiService: ApiService) { }

    getPath(item: Item) : string {
	return "collections/" + item.collection.name + "/items/" + item.id
    }
    
    modifyItem(item: Item, data: any) {

	return this.apiService.patch(this.getPath(item), data)
            .map((rsp: Response) => {
                if (rsp.status != 204) {
                    throw new Error('Impossible to modify item '
				    + item.id + '/' + item.name
				    + ' from collection '
				    + item.collection.name + ': ' + rsp.status);
                }

		item.data = data
		console.log("Modify item:", item)
            })
    }

    deleteItem(item: Item) {
	
        return this.apiService.delete(this.getPath(item))
            .map((rsp: Response) => {
                if (rsp.status != 204) {
                    throw new Error('Impossible to delete item '
				    + item.id + '/' + item.name
				    + ' from collection '
				    + item.collection.name + ': ' + rsp.status);
                }

		item.collection.deleteItem(item)
		console.log("Delete item:", item)		
            })
    }    
}
