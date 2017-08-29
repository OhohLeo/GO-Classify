import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ImportsService } from './imports.service'

import { ImportsComponent } from './imports.component'
import { DirectoryCreateComponent } from './directory/create.component'
import { DirectoryDisplayComponent } from './directory/display.component'
import { EmailCreateComponent } from './email/create.component'
import { EmailDisplayComponent } from './email/display.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    providers: [ImportsService],
    declarations: [ImportsComponent,
				   DirectoryCreateComponent,
				   DirectoryDisplayComponent,
				   EmailCreateComponent,
				   EmailDisplayComponent],
    exports: [ImportsComponent]
})

export class ImportsModule { }
