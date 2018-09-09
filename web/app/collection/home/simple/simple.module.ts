import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { PdfViewerModule } from 'ng2-pdf-viewer'

import { ApiService } from '../../../api.service'

import { SimpleCollectionComponent } from './collection.component'
import { SimpleItemsComponent } from './items.component'
import { SimpleIconItemComponent } from './iconItem.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
		PdfViewerModule
    ],
    providers: [ApiService],
    declarations: [
        SimpleCollectionComponent,
        SimpleItemsComponent,
		SimpleIconItemComponent
    ],
    exports: [
        SimpleCollectionComponent
    ]
})

export class SimpleModule { }
