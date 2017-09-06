import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { FileCreateComponent } from './create.component'
import { FileDisplayComponent } from './display.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    declarations: [FileCreateComponent,
				   FileDisplayComponent],
    exports: [FileCreateComponent]
})

export class FileModule { }
