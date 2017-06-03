import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { FileModule } from '../imports/file/file.module'
import { MoviesModule } from '../collections/movies/movies.module'
import { ConfigModule } from '../config/config.module'

import { ItemComponent } from './item.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
        ConfigModule,
        FileModule,
        MoviesModule
    ],
    providers: [],
    declarations: [ItemComponent],
    exports: [ItemComponent]
})

export class ItemModule { }
