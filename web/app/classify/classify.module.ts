import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { FileModule } from '../imports/file/file.module'
import { MoviesModule } from '../collections/movies/movies.module'

import { ClassifyService } from './classify.service'

import { ClassifyComponent } from './classify.component'
import { DetailComponent } from './detail.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
		FileModule,
		MoviesModule
    ],
    providers: [ClassifyService],
    declarations: [ClassifyComponent, DetailComponent],
    exports: [ClassifyComponent]
})

export class ClassifyModule { }
