import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { SimpleModule } from './home/simple/simple.module'

import { CollectionService } from './collection.service'
import { ApiService } from '../api.service'

import { CollectionsComponent } from './collections.component'
import { CreateCollectionComponent } from './create.component'
import { ModifyCollectionComponent } from './modify.component'
import { DisplayCollectionComponent } from './display.component'
import { DeleteCollectionComponent } from './delete.component'

import { ListCollectionComponent } from './list/list.component'
import { WorldCollectionComponent } from './world/world.component'
import { TimeLineCollectionComponent } from './timeline/timeline.component'

@NgModule({
    imports: [CommonModule,
        BrowserModule,
        FormsModule,
        SimpleModule],
    providers: [ApiService, CollectionService],
    declarations: [CollectionsComponent,
        CreateCollectionComponent,
        ModifyCollectionComponent,
        DisplayCollectionComponent,
        DeleteCollectionComponent,
        ListCollectionComponent,
        WorldCollectionComponent,
        TimeLineCollectionComponent],
    exports: [CollectionsComponent, DisplayCollectionComponent],
})

export class CollectionsModule { }
