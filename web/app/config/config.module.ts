import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { BrowserModule } from '@angular/platform-browser'
import { FormsModule } from '@angular/forms'

import { ToolsModule } from '../tools/tools.module'
import { ConfigComponent } from './config.component'
import { ConfigService } from './config.service'

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule,
		ToolsModule
    ],
    declarations: [ConfigComponent],
    providers: [ConfigService],
    exports: [ConfigComponent]
})

export class ConfigModule {}
