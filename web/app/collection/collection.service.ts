import { Injectable } from '@angular/core'
import { Response } from '@angular/http'
import { Observable } from 'rxjs/Rx'

import { ApiService, CollectionStatus } from '../api.service'

import { Collection } from './collection'
import { Items } from '../items/items'
import { Item, ItemObserver, ItemEvent } from '../items/item'
import { Event } from '../api.service'

@Injectable()
export class CollectionService {
    constructor(private apiService: ApiService) { }
}
