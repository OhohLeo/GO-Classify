import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { ConfigModule } from '../../config/config.module'

import { DetailFileComponent } from './detail.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
		ConfigModule
    ],
    declarations: [DetailFileComponent],
    exports: [DetailFileComponent]
})

export class FileModule { }
