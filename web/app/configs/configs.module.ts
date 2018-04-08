import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'
import { ToolsModule } from '../tools/tools.module'
import { ParamsModule } from '../params/params.module'

import { ConfigsComponent } from './configs.component'
import { ConfigComponent } from './config.component'
import { ConfigMultiComponent } from './multi.component'
import { ConfigRefComponent } from './ref.component'
import { ConfigsService } from './configs.service'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
        ToolsModule,
        ParamsModule
    ],
    declarations: [ConfigsComponent, ConfigMultiComponent, ConfigRefComponent, ConfigComponent] ,
    providers: [ConfigsService],
    exports: [ConfigsComponent, ConfigComponent]
})

export class ConfigsModule { }
