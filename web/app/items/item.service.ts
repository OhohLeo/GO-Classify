import { Injectable } from '@angular/core'
import { Response } from '@angular/http'
import { Observable } from 'rxjs/Rx'
import { ApiService } from './../api.service'
import { Item } from '../items/item'

@Injectable()
export class ItemService {

    constructor(private apiService: ApiService) { }

	getReferences(item: Item) {

		// return this.apiService.get(this.getCollectionPath(item), data)
        //     .map((rsp: Response) => {
        //        	console.log("Modify item:", item)
        //     })
	}

    modifyItem(item: Item, data: any) {

		return this.apiService.patch(item.getPath(), data)
            .map((rsp: Response) => {
                if (rsp.status != 204) {
                    throw new Error('Impossible to modify item '
									+ item.getPath()+ '/' + item.id +
									': ' + rsp.status);
                }

				item.data = data
				console.log("Modify item:", item)
            })
    }

    deleteItem(item: Item) {

        return this.apiService.delete(item.getPath())
            .map((rsp: Response) => {
                if (rsp.status != 204) {
                    throw new Error('Impossible to delete item '
									+ item.getPath()+ '/' + item.id +
									': ' + rsp.status);
                }

				item.delete()
				console.log("Delete item:", item)
            })
    }
}
