import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { ParamsModule } from '../params/params.module'

import { ExportsService } from './exports.service'

import { ExportsComponent } from './exports.component'
import { FileCreateComponent } from './file/create.component'
import { FileDisplayComponent } from './file/display.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
        ParamsModule
    ],
    providers: [ExportsService],
    declarations: [ExportsComponent,
        FileCreateComponent,
        FileDisplayComponent],
    exports: [ExportsComponent]
})

export class ExportsModule { }
