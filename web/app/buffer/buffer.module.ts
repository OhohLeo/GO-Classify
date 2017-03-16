import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { FileModule } from '../imports/file/file.module'
import { MoviesModule } from '../collections/movies/movies.module'

import { BufferService } from './buffer.service'

import { BufferComponent } from './buffer.component'
import { DetailComponent } from './detail.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
        FileModule,
        MoviesModule
    ],
    providers: [BufferService],
    declarations: [BufferComponent, DetailComponent],
    exports: [BufferComponent]
})

export class BufferModule { }
