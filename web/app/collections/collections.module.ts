import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { ApiService } from '../api.service';

import { CollectionsComponent } from './collections.component';
import { CreateCollectionComponent } from './create.component';
import { ModifyCollectionComponent } from './modify.component';
import { DeleteCollectionComponent } from './delete.component';

@NgModule({
    imports: [CommonModule, BrowserModule, FormsModule],
    providers: [ApiService],
    declarations: [CollectionsComponent,
        CreateCollectionComponent,
        ModifyCollectionComponent,
        DeleteCollectionComponent],
    exports: [CollectionsComponent],
})

export class CollectionsModule { }
