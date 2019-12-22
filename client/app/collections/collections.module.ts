import { BrowserModule } from '@angular/platform-browser'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { NgModule } from '@angular/core'

import { CollectionModule } from '../collection/collection.module'

import { CollectionsService } from './collections.service'
import { ApiService } from '../api.service'

import { ListCollectionsComponent } from './list.component'

@NgModule({
    imports: [
		CollectionModule,
		CommonModule,
        BrowserModule,
        FormsModule
	],
    providers: [
		ApiService,
		CollectionsService
	],
    declarations: [ListCollectionsComponent],
    exports: [ListCollectionsComponent],
})

export class CollectionsModule { }
