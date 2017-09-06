import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { FileModule } from '../imports/file/file.module'
import { MoviesModule } from '../collections/movies/movies.module'
import { ConfigModule } from '../config/config.module'
import { ToolsModule } from '../tools/tools.module'

import { BufferService } from './buffer.service'

import { BufferComponent } from './buffer.component'
import { BufferItemComponent } from './item.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
        ConfigModule,
        FileModule,
        MoviesModule,
        ToolsModule
    ],
    providers: [BufferService],
    declarations: [BufferComponent, BufferItemComponent],
    exports: [BufferComponent, BufferItemComponent]
})

export class BufferModule { }
