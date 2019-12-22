import { BrowserModule } from '@angular/platform-browser'
import { CommonModule } from '@angular/common'
import { ConfigsModule } from '../configs/configs.module'
import { FormsModule } from '@angular/forms'
import { NgModule } from '@angular/core'
import { ParamsModule } from '../params/params.module'
import { ReferencesModule } from '../references/references.module'
import { ToolsModule } from '../tools/tools.module'

import { ImportsService } from './imports.service'

import { ImportsComponent } from './imports.component'
import { ImportsCreateComponent } from './create.component'
import { ImportsDisplayComponent } from './display.component'
import { DirectoryCreateComponent } from './directory/create.component'
import { DirectoryDisplayComponent } from './directory/display.component'
import { ImapCreateComponent } from './imap/create.component'
import { ImapSearchComponent } from './imap/search.component'
import { ImapDisplayComponent } from './imap/display.component'

@NgModule({
    imports: [
        BrowserModule,
        CommonModule,
	ConfigsModule,
        FormsModule,
        ParamsModule,
        ToolsModule,
	ReferencesModule
    ],
    providers: [ImportsService],
    declarations: [
	ImportsComponent,
	ImportsCreateComponent,
	ImportsDisplayComponent,
        DirectoryCreateComponent,
        DirectoryDisplayComponent,
        ImapCreateComponent,
        ImapSearchComponent,
        ImapDisplayComponent
    ],
    exports: [ImportsComponent,
	      ImportsCreateComponent,
	      ImportsDisplayComponent]
})

export class ImportsModule { }
