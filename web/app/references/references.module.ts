import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ReferencesService } from './references.service'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
    ],
    declarations: [
	],
    providers: [ReferencesService],
    exports: []
})

export class ReferencesModule { }
