import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ClassifyMovieComponent } from './classify.component'
import { DetailMovieComponent } from './detail.component'
import { WebsiteMovieComponent } from './website.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    declarations: [ClassifyMovieComponent, DetailMovieComponent, WebsiteMovieComponent],
    exports: [ClassifyMovieComponent, DetailMovieComponent, WebsiteMovieComponent]
})

export class MoviesModule { }
