import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ParamsService } from './params.service'

import { ParamsPathComponent } from './path.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    providers: [ParamsService],
    declarations: [ParamsPathComponent],
    exports: [ParamsPathComponent]
})

export class ParamsModule { }
