import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { CollectionsModule } from './collections/collections.module';
import { ImportsModule } from './imports/imports.module';

import { ClassifyService } from './classify.service';

import { AppComponent } from './app.component';

@NgModule({
    imports: [
        CommonModule,
        HttpModule,
        BrowserModule,
        FormsModule,
        CollectionsModule,
        ImportsModule
    ],
    providers: [ClassifyService],
    declarations: [AppComponent],
    bootstrap: [AppComponent]
})

export class AppModule { }
