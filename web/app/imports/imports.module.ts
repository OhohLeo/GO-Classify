import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { ImportsService } from './imports.service';

import { ImportsComponent } from './imports.component';
import { ListComponent } from './list.component';
import { DirectoryComponent } from './directory.component';

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    providers: [ImportsService],
    declarations: [ImportsComponent, ListComponent, DirectoryComponent],
    exports: [ImportsComponent]
})

export class ImportsModule {}
