import {Component, OnInit, ViewChild} from '@angular/core';
import {ClassifyService, WebSocketStatus} from './classify.service';
import {Collection} from './collections/collection';
import {CollectionsComponent} from './collections/collections.component';
import {ImportDirectoryComponent} from './imports/directory.component';

declare var jQuery:any;

enum AppStatus {
    HOME = 1,
    IMPORT,
    EXPORT,
    CONFIG
}

@Component({
    selector: 'classify',
    templateUrl: 'app/app.component.html',
    providers: [ClassifyService],
    directives: [CollectionsComponent,
                 ImportDirectoryComponent]
})

export class AppComponent implements OnInit {
    @ViewChild(CollectionsComponent) collections: CollectionsComponent

    public appStatus = AppStatus
    public status = AppStatus.HOME
    public websocketStatus: WebSocketStatus

    public title = "Classify"

    public collection: Collection

    public modalTitle: string
    public modalMsg: string

    constructor (private classifyService: ClassifyService) {}

    ngOnInit() {

        // Initialisation de la side bar
        jQuery(".button-collapse").sideNav();

        // Initialisation de la websocket
        this.classifyService.initWebSocket()
            .subscribe((status: WebSocketStatus) => {
                if (this.classifyService.status == WebSocketStatus.ERROR) {
                    this.onError("Websocket", "Impossible to connect the websocket!")
                }
                else if (this.classifyService.status == WebSocketStatus.OPEN) {
                    console.log("websocket ok")
                    this.stopModal()
                    this.onHome()
                }
            })

        this.classifyService.setOnErrors(this.onError)
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

    onNewCollection() {
        if (this.collections) {
            this.collections.onNewCollection()
        }
    }

    onSelectCollection() {
        this.status = AppStatus.HOME
    }

    onCollectionChoosed(collection: Collection) {
        this.status = AppStatus.HOME
        console.log("COLLECTION", collection)
    }

    resetCollectionState() {
        if (this.collections)
            this.collections.resetCollectionState()
    }
}
