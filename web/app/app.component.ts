import { Component, OnInit, ViewChild } from '@angular/core';

import { ClassifyService, CollectionStatus } from './classify.service';

import { Collection } from './collections/collection';;

import { CollectionsComponent } from './collections/collections.component'

declare var jQuery: any;

enum AppStatus {
    NONE = 0,
    HOME,
    IMPORT,
    EXPORT,
    CONFIG
}

@Component({
    selector: 'classify',
    templateUrl: 'app/app.component.html',
})

export class AppComponent implements OnInit {
    @ViewChild(CollectionsComponent) collections: CollectionsComponent

    public appStatus = AppStatus
    public status = AppStatus.NONE

    public title = "Classify"

    public collection: Collection

    public modalTitle: string
    public modalMsg: string

    constructor(private classifyService: ClassifyService) { }

    ngOnInit() {

        // Initialisation de la side bar
        jQuery(".button-collapse").sideNav();

        // Inscription au flux
        this.classifyService.getStream()
            .subscribe((event) => {
                console.log(event)
            })

        this.classifyService.subscribeCollectionChange(
            (collection: Collection, status: CollectionStatus) => {
                console.log("CHANGE COLLECTION", collection, status)
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
