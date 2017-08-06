import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { ItemModule } from '../item/item.module';

import { HomeService } from './home.service'

import { HomeComponent } from './home.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
        ItemModule
    ],
    providers: [HomeService],
    declarations: [HomeComponent],
    exports: [HomeComponent]
})

export class HomeModule { }
