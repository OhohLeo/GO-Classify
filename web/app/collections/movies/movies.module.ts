import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { BufferMovieComponent } from './buffer.component'
import { DetailMovieComponent } from './detail.component'
import { WebsiteMovieComponent } from './website.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    declarations: [BufferMovieComponent, DetailMovieComponent, WebsiteMovieComponent],
    exports: [BufferMovieComponent, DetailMovieComponent, WebsiteMovieComponent]
})

export class MoviesModule { }
