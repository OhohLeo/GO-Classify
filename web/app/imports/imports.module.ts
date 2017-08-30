import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { ToolsModule } from '../tools/tools.module'

import { ImportsService } from './imports.service'

import { ImportsComponent } from './imports.component'
import { DirectoryCreateComponent } from './directory/create.component'
import { DirectoryDisplayComponent } from './directory/display.component'
import { EmailCreateComponent } from './email/create.component'
import { EmailSearchComponent } from './email/search.component'
import { EmailDisplayComponent } from './email/display.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
		ToolsModule
    ],
    providers: [ImportsService],
    declarations: [ImportsComponent,
				   DirectoryCreateComponent,
				   DirectoryDisplayComponent,
				   EmailCreateComponent,
				   EmailSearchComponent,
				   EmailDisplayComponent],
    exports: [ImportsComponent]
})

export class ImportsModule { }
