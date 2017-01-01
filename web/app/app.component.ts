import { Component, OnInit, ViewChild } from '@angular/core';

import { ApiService, CollectionStatus, Event } from './api.service';
import { ImportsService } from './imports/imports.service';

import { Collection } from './collections/collection';

import { CollectionsComponent } from './collections/collections.component'
import { ClassifyComponent } from './classify/classify.component'

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
    @ViewChild(ClassifyComponent) classify: ClassifyComponent

    public appStatus = AppStatus
    public status = AppStatus.NONE

    public title = "Classify"

    public collection: Collection

    public modalTitle: string
    public modalMsg: string

    constructor(private apiService: ApiService,
        private importsService: ImportsService) { }

    ngOnInit() {

        // Initialisation de la side bar
        jQuery(".button-collapse").sideNav();

        // Logo d'importation
        let importsLoop = jQuery("i#imports-loop")
        let importsRunningNb = 0;

        // Inscription au flux
        this.apiService.getStream()
            .subscribe((e: Event) => {
                console.log("EVENT!", e)

                if (e.event === "start") {
                    // restart application
                    window.location.replace("/");
                    return;
                }

                // Send data to the import service
                if (e.event.startsWith("import")) {

                    // Send notifications to the imports list
                    this.importsService.addEvent(e);

                    // Display imports status
                    if (e.event.endsWith("status")) {

                        // Status 'TRUE': rotate refresh logo
                        if (e.data) {
                            importsLoop.addClass("rotation")
                            importsRunningNb++
                        }
                        // Status 'FALSE'
                        else if (importsRunningNb > 0) {
                            importsRunningNb--
                        }

                        // No more imports : stop logo rotation
                        if (importsRunningNb < 1) {
                            importsLoop.removeClass("rotation")
                        }
                    }

                    return;
                }

                // Items reception
                if (e.event.startsWith("item")) {
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

    onClassify() {
        this.classify.startModal()
    }

    onNewState(nextStatus: AppStatus) {
        this.resetCollectionState()
        this.status = nextStatus
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
