import { Component, OnInit, ViewChild } from '@angular/core';

import { ApiService, CollectionStatus, Event } from './api.service';
import { ImportsService } from './imports/imports.service';
import { BufferService } from './buffer/buffer.service';

import { Collection } from './collections/collection';

import { CollectionsComponent } from './collections/collections.component'
import { BufferComponent } from './buffer/buffer.component'

declare var jQuery: any;

enum AppStatus {
    NONE = 0,
    HOME,
    IMPORT,
    EXPORT,
    CONFIG
}

@Component({
    selector: 'app',
    templateUrl: './app.component.html',
})

export class AppComponent implements OnInit {
    @ViewChild(CollectionsComponent) collections: CollectionsComponent
    @ViewChild(BufferComponent) buffer: BufferComponent

    public appStatus = AppStatus
    public status = AppStatus.NONE

    public title = "Buffer"

    public collection: Collection

    public modalTitle: string
    public modalMsg: string

    private importsLoop: any
    private importsRunningNb: number

    constructor(private apiService: ApiService,
        private importsService: ImportsService,
        private bufferService: BufferService) { }

    ngOnInit() {

        // Initialisation de la side bar
        jQuery(".button-collapse").sideNav();

        // Logo d'importation
        this.importsLoop = jQuery("i#imports-loop")
        this.importsRunningNb = 0;

        // Inscription au flux
        this.apiService.getStream()
            .subscribe((e: Event) => {

                //console.log("EVENT!", e)

                if (e.event === "start") {
                    // restart application
                    window.location.replace("/");
                    return;
                }

                // Send data to the import service
                if (e.event.startsWith("import")) {
                    this.handleImport(e);
                    return;
                }

                // Items reception
                if (e.event.startsWith("item")) {
                    this.handleItem(e);
                    return;
                }
            })

        this.apiService.subscribeCollectionChange(
            (collection: Collection, status: CollectionStatus) => {

                if (collection === undefined) {
                    this.onChangeCollection()
                    return
                }

                this.title = collection.name
                this.collection = collection

                switch (status) {
                    case CollectionStatus.CREATED:
                    case CollectionStatus.MODIFIED:
                    case CollectionStatus.SELECTED:
                        this.onHome()
                        break;
                    case CollectionStatus.DELETED:
                        this.status = AppStatus.NONE
                        break
                }
            })
    }

    onHome() {
        this.onNewState(AppStatus.HOME)
    }

    onImport() {
        this.onNewState(AppStatus.IMPORT)
    }

    onExport() {
        this.onNewState(AppStatus.EXPORT)
    }

    onConfig() {
        this.onNewState(AppStatus.CONFIG)
    }

    onBuffer() {
        this.buffer.start()
    }

    onNewState(nextStatus: AppStatus) {
        this.resetCollectionState()
        this.status = nextStatus
    }

    // Gestion des nouveaux imports
    handleImport(e: Event) {
        // Send notifications to the imports list
        this.importsService.addEvent(e);

        // Display imports status
        if (e.event.endsWith("status")) {

            // Status 'TRUE': rotate refresh logo
            if (e.data) {
                this.importsLoop.addClass("rotation")
                this.importsRunningNb++
            }
            // Status 'FALSE'
            else if (this.importsRunningNb > 0) {
                this.importsRunningNb--
            }

            // No more imports : stop logo rotation
            if (this.importsRunningNb < 1) {
                this.importsLoop.removeClass("rotation")
            }
        }
    }

    // Gestion des éléments à classer
    handleItem(e: Event) {
        // Replace type
        e.event = e.event.substring(e.event.indexOf('/') + 1)

        // Send notifications to the imports list
        this.bufferService.addEvent(e)
    }

    onError(title: string, msg: string) {

        console.error(title + "error :" + msg)

        this.modalTitle = title + " error!"
        this.modalMsg = msg
        jQuery('#modal').openModal()
    }

    stopModal() {
        jQuery('#modal').closeModal()
    }

    // Affiche la création d'une nouvelle collection
    onNewCollection() {
        if (this.collections) {
            this.collections.onNewCollection()
        }
    }

    // Affiche la liste des collections à sélectionner
    onChangeCollection() {
        if (this.collections
            && this.collections.onChooseCollection(undefined)) {
            this.status = AppStatus.NONE
        }
    }

    // Affiche la liste des collections à sélectionner si aucune
    // collection n'est actuellement sélectionnée
    resetCollectionState() {
        if (this.collections)
            this.collections.resetCollectionState()
    }
}
