import { BrowserModule } from '@angular/platform-browser'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { NgModule } from '@angular/core'

import { SimpleModule } from './home/simple/simple.module'
import { ItemsModule } from '../items/items.module'

import { CollectionService } from './collection.service'
import { ApiService } from '../api.service'

import { CreateCollectionComponent } from './create.component'
import { DisplayCollectionComponent } from './display.component'
import { ModifyCollectionComponent } from './modify.component'
import { DeleteCollectionComponent } from './delete.component'

import { ListCollectionComponent } from './list/list.component'
import { TimeLineCollectionComponent } from './timeline/timeline.component'
import { WorldCollectionComponent } from './world/world.component'

@NgModule({
    imports: [
		CommonModule,
        BrowserModule,
		FormsModule,
		ItemsModule,
        SimpleModule
	],
    providers: [
		ApiService,
		CollectionService
	],
    declarations: [
		CreateCollectionComponent,
		DisplayCollectionComponent,
		ModifyCollectionComponent,
		DeleteCollectionComponent,

		ListCollectionComponent,
		TimeLineCollectionComponent,
		WorldCollectionComponent
	],
    exports: [
		CreateCollectionComponent,
		DisplayCollectionComponent,
		ModifyCollectionComponent,
		DeleteCollectionComponent,
	],
})

export class CollectionModule { }
