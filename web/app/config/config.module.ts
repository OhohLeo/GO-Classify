import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ConfigComponent } from './config.component'
import { StringListComponent } from './stringlist.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    declarations: [ConfigComponent, StringListComponent],
    exports: [ConfigComponent, StringListComponent]
})

export class ConfigModule { }
