import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ApiService } from '../api.service'

import { FilterComponent } from './filter.component'

@NgModule({
    imports: [CommonModule,
        BrowserModule,
        FormsModule],
    providers: [ApiService],
    declarations: [FilterComponent],
    exports: [FilterComponent],
})

export class FilterModule { }
