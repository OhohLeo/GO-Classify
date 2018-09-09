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

import { TweaksComponent } from './tweaks/tweaks.component'
import { TweaksDatasComponent } from './tweaks/datas.component'
import { TweaksFieldsComponent } from './tweaks/fields.component'

import { ConfigsService } from './configs.service'
import { TweaksService } from './tweaks/tweaks.service'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
        ToolsModule,
        ParamsModule
    ],
    declarations: [
		ConfigsComponent,
		ConfigMultiComponent,
		ConfigRefComponent,
		ConfigComponent,
		TweaksComponent,
		TweaksDatasComponent,
		TweaksFieldsComponent
	],
    providers: [ConfigsService, TweaksService],
    exports: [ConfigsComponent, ConfigComponent, TweaksComponent]
})

export class ConfigsModule { }
