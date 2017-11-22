import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { SimpleCollectionComponent } from './collection.component'
import { SimpleItemsComponent } from './items.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    declarations: [
        SimpleCollectionComponent,
		SimpleItemsComponent
    ],
    exports: [
        SimpleCollectionComponent
    ]
})

export class SimpleModule {}
