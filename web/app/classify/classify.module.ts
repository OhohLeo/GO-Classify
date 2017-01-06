import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms'

import { ClassifyService } from './classify.service'

import { ClassifyComponent } from './classify.component';

@NgModule({
    imports: [
        CommonModule,
        BrowserModule,
        FormsModule
    ],
    providers: [ClassifyService],
    declarations: [ClassifyComponent],
    exports: [ClassifyComponent]
})

export class ClassifyModule { }
