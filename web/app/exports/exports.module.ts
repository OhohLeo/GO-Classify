import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ExportsService } from './exports.service'

import { ExportsComponent } from './exports.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    providers: [ExportsService],
    declarations: [ExportsComponent],
    exports: [ExportsComponent]
})

export class ExportsModule { }
