import { BrowserModule } from '@angular/platform-browser'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { NgModule } from '@angular/core'

import { ReferencesModule } from '../references/references.module'
import { ImportsModule } from '../imports/imports.module'
import { ExportsModule } from '../exports/exports.module'
import { ToolsModule } from '../tools/tools.module'
import { ConfigsModule } from '../configs/configs.module'

import { WorkflowService } from './workflow.service'

import { WorkflowComponent } from './workflow.component'
import { SelectComponent } from './select.component'
import { InstanceComponent } from './instance.component'
import { CollectionComponent } from './collection.component'


@NgModule({
    imports: [
        BrowserModule,
        CommonModule,
        FormsModule,
	ReferencesModule,
	ImportsModule,
	ExportsModule,
	ToolsModule,
	ConfigsModule
    ],
    providers: [WorkflowService],
    declarations: [
	WorkflowComponent,
	SelectComponent,
	InstanceComponent,
	CollectionComponent,
    ],
    exports: [WorkflowComponent]
})

export class WorkflowModule {}
