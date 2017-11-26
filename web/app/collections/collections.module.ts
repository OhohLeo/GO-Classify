import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { SimpleModule } from './collection/home/simple/simple.module'

import { CollectionsService } from './collections.service'
import { ApiService } from '../api.service'

import { ListCollectionsComponent } from './list.component'
import { CreateCollectionComponent } from './collection/create.component'
import { ModifyCollectionComponent } from './collection/modify.component'
import { DisplayCollectionComponent } from './collection/display.component'
import { DeleteCollectionComponent } from './collection/delete.component'

import { ListCollectionComponent } from './collection/list/list.component'
import { WorldCollectionComponent } from './collection/world/world.component'
import { TimeLineCollectionComponent } from './collection/timeline/timeline.component'

@NgModule({
    imports: [CommonModule,
        BrowserModule,
        FormsModule,
        SimpleModule],
    providers: [ApiService, CollectionsService],
    declarations: [ListCollectionsComponent,
        CreateCollectionComponent,
        ModifyCollectionComponent,
        DisplayCollectionComponent,
        DeleteCollectionComponent,
        ListCollectionComponent,
        WorldCollectionComponent,
        TimeLineCollectionComponent],
    exports: [ListCollectionsComponent, DisplayCollectionComponent],
})

export class CollectionsModule { }
