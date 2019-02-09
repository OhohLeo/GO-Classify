import { BrowserModule } from '@angular/platform-browser'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { NgModule } from '@angular/core'

import { ReferencesModule } from '../references/references.module'
import { ImportsModule } from '../imports/imports.module'
import { ExportsModule } from '../exports/exports.module'
import { ToolsModule } from '../tools/tools.module'

import { MappingsService } from './mappings.service'

import { MappingsComponent } from './mappings.component'
import { CreateComponent } from './create.component'


@NgModule({
    imports: [
        BrowserModule,
        CommonModule,
        FormsModule,
	ReferencesModule,
	ImportsModule,
	ExportsModule,
	ToolsModule
    ],
    providers: [MappingsService],
    declarations: [
	MappingsComponent,
	CreateComponent,	
    ],
    exports: [MappingsComponent]
})

export class MappingsModule {}
