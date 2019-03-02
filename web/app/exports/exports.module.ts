import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { ParamsModule } from '../params/params.module'
import { ToolsModule } from '../tools/tools.module'

import { ExportsService } from './exports.service'

import { ExportsComponent } from './exports.component'
import { ExportsCreateComponent } from './create.component'
import { ExportsDisplayComponent } from './display.component'
import { FileCreateComponent } from './file/create.component'
import { FileDisplayComponent } from './file/display.component'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
	ToolsModule,
        ParamsModule
    ],
    providers: [ExportsService],
    declarations: [ExportsComponent,
		   ExportsCreateComponent,
		   ExportsDisplayComponent,
		   FileCreateComponent,
		   FileDisplayComponent],
    exports: [ExportsComponent, ExportsCreateComponent, ExportsDisplayComponent]
})

export class ExportsModule { }
