import { BrowserModule } from '@angular/platform-browser'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { NgModule } from '@angular/core'

import { ReferencesModule } from '../references/references.module'
import { ImportsModule } from '../imports/imports.module'
import { ExportsModule } from '../exports/exports.module'
import { ToolsModule } from '../tools/tools.module'

import { WorkflowService } from './workflow.service'

import { WorkflowComponent } from './workflow.component'
import { CreateComponent } from './create.component'
import { DisplayComponent } from './display.component'
import { ElementComponent } from './element.component'
import { CanvasComponent } from './canvas.component'


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
    providers: [WorkflowService],
    declarations: [
	WorkflowComponent,
	CreateComponent,
	DisplayComponent,
	ElementComponent,
	CanvasComponent,
    ],
    exports: [WorkflowComponent]
})

export class WorkflowModule {}
