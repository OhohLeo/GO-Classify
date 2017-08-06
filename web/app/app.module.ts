import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { CollectionsModule } from './collections/collections.module';
import { ConfigModule } from './config/config.module';
import { ImportsModule } from './imports/imports.module';
import { ExportsModule } from './exports/exports.module';
import { BufferModule } from './buffer/buffer.module';
import { HomeModule } from './home/home.module';
import { ItemModule } from './item/item.module';

import { ApiService } from './api.service';

import { AppComponent } from './app.component';

@NgModule({
    imports: [
        CommonModule,
        HttpModule,
        BrowserModule,
        FormsModule,
        CollectionsModule,
        ImportsModule,
        ExportsModule,
        HomeModule,
        BufferModule,
        ConfigModule,
        ItemModule
    ],
    providers: [ApiService],
    declarations: [AppComponent],
    bootstrap: [AppComponent]
})

export class AppModule { }
