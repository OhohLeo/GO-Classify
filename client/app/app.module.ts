import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { HttpModule } from '@angular/http'

import { BufferModule } from './buffer/buffer.module'
import { CollectionModule } from './collection/collection.module'
import { CollectionsModule } from './collections/collections.module'
import { ConfigsModule } from './configs/configs.module'
import { WorkflowModule } from './workflow/workflow.module'
import { ImportsModule } from './imports/imports.module'
import { ExportsModule } from './exports/exports.module'
import { FilterModule } from './filter/filter.module'
import { ReferencesModule } from './references/references.module'
import { ToolsModule } from './tools/tools.module'

import { ApiService } from './api.service'

import { AppComponent } from './app.component'

@NgModule({
    imports: [
     	CommonModule,
        BrowserModule,
        BufferModule,
        CollectionModule,
        CollectionsModule,
        ConfigsModule,
        FilterModule,
        FormsModule,
        HttpModule,
	WorkflowModule,
        ImportsModule,
	ExportsModule,
        ToolsModule,
	ReferencesModule,
    ],
    providers: [ApiService],
    declarations: [AppComponent],
    bootstrap: [AppComponent]
})

export class AppModule { }
