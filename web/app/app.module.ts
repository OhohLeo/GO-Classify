import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { CollectionsModule } from './collections/collections.module';
import { ImportsModule } from './imports/imports.module';

import { ApiService } from './api.service';

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
    providers: [ApiService],
    declarations: [AppComponent],
    bootstrap: [AppComponent]
})

export class AppModule { }
