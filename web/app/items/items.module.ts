import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ApiService } from '../api.service'

import  { ItemComponent } from './item.component'

@NgModule({
    imports: [CommonModule,
              BrowserModule,
              FormsModule],
    providers: [ApiService],
    declarations: [ItemComponent],
    exports: [ItemComponent],
})

export class ItemsModule { }
