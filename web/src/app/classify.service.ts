import {Injectable} from '@angular/core';
import {Http, Response, RequestOptions, Headers} from '@angular/http';
import {Observable} from 'rxjs/Rx';
import {Collection, Imports} from './collections/collection';

export enum WebSocketStatus {
    NONE = 0,
    CONNECTING,
    OPEN,
    CLOSING,
    CLOSE,
    ERROR,
}

export enum CollectionStatus {
    NONE = 0,
    CREATED,
    SELECTED,
    MODIFIED,
    DELETED
}


@Injectable()
export class ClassifyService {

    private url = "http://localhost:3333/"
    private references: any
    private collections: Collection[]
    private websocket: any
    private websocketStatus = {
        'open': WebSocketStatus.OPEN,
        'error': WebSocketStatus.ERROR,
        'close': WebSocketStatus.CLOSE,
    }
    private websocketTimer

    private onCollectionChange: (collection: Collection,
                                 status: CollectionStatus) => void

    public status = WebSocketStatus.NONE

    public collectionSelected: Collection

    constructor (private http: Http) {}

    subscribeCollectionChange(cb: (collection: Collection,
                                   status: CollectionStatus) => void) {
        this.onCollectionChange = cb
    }

    setCollection(collection: Collection, status: CollectionStatus) {

        if (status == CollectionStatus.SELECTED
            && this.collectionSelected == collection) {
            return
        }

        this.collectionSelected = collection

        if (this.onCollectionChange) {
            this.onCollectionChange(collection, status)
        }
    }

    initWebSocket(): Observable<WebSocketStatus>{
        return Observable.create(
            observer => this.connectWebSocket(observer))
    }

    connectWebSocket(observer) {

        console.log("websocket connecting...")

        // Etablissement de la connexion avec la websocket
        this.websocket = new WebSocket('ws://localhost:3333/ws')

        observer.next(WebSocketStatus.CONNECTING)

        let handleWebSocketStatus = (expected) => {
            return (evt) => {
                let status = this.getWebSocketStatus(evt, expected)

                // Vérification de l'état du status
                if (status == undefined) {
                    return
                }

                // Vérification que le status a bien changé
                if (this.status == expected) {
                    return
                }

                // Attribution du nouveau status
                this.status = expected

                // En cas de status d'erreur ou de fermeture
                // inattendue, lorsque le timer n'est pas défini : on
                // relance périodiquement la tentative de connexion
                if (this.websocketTimer === undefined
                    && (this.status === WebSocketStatus.ERROR
                        || this.status === WebSocketStatus.CLOSE))  {
                    console.log(
                        "websocket ",
                        this.status === WebSocketStatus.CLOSE ? "close" : "error")
                    this.websocketTimer = setTimeout(
                        () => {
                            console.log("websocket retry ...")
                            this.websocketTimer = undefined
                            this.connectWebSocket(observer)
                        }, 5000)
                }

                observer.next(expected)
            }
        }

        this.websocket.onopen = handleWebSocketStatus(WebSocketStatus.OPEN)
        this.websocket.onerror = handleWebSocketStatus(WebSocketStatus.ERROR)
        this.websocket.onclose = handleWebSocketStatus(WebSocketStatus.CLOSE)
        this.websocket.onmessage = (evt) => {
            console.log("RECEIVED: " + evt.data)
        }
    }

    getWebSocketStatus(evt, expected: WebSocketStatus): WebSocketStatus{

        let status = this.websocketStatus[evt.type]
        if (status == undefined) {
            console.error("Unknown received websocket status type: " + evt.type)
            return undefined
        }

        if (status != expected) {
            console.error("Websocket status error: expected " + expected
                          + ", received " + status)
            return undefined
        }

        return status
    }

    getWebSocket(): Observable<any>{
        return Observable.fromEvent(this.websocket,'message')
    }

    getOptions() {
        return new RequestOptions({
            headers: new Headers({ 'Content-Type': 'application/json' })
        })
    }

    getCollectionUrl() {

        if (this.collectionSelected == undefined) {
            this.handleError("Select a collection first!")
            return undefined
        }

        return this.url + "collections/" + this.collectionSelected.name
    }

    // Create a new collection
    newCollection(collection: Collection) {

        return this.http.post(this.url + "collections",
                              JSON.stringify(collection),
                              this.getOptions())
            .map((res: Response) => {
                if (res.status != 204) {
                    throw new Error('Impossible to create new collection: ' + res.status);
                }

                // Ajoute la collection nouvellement créée
                this.collections.push(collection)

                // Choisit automatiquement la nouvelle collection
                this.setCollection(collection, CollectionStatus.CREATED)
            })
            .catch(this.handleError);
    }

    // Modify an existing collection
    modifyCollection(name: string, collection: Collection) {

        return this.http.patch(this.url + "collections/" + name,
                               JSON.stringify(collection),
                               this.getOptions())
            .map((res: Response) => {
                if (res.status != 204) {
                    throw new Error('Impossible to modify collection '
                                    + name + ': ' + res.status);
                }

                // Replace the collection from the list
                for (let i in this.collections) {
                    if (this.collections[i].name === name) {
                        this.collections[i] = collection
                        break
                    }
                }

                // Choisit automatiquement la nouvelle collection
                this.setCollection(collection, CollectionStatus.MODIFIED)
            })
            .catch(this.handleError);
    }

    // Delete an existing collection
    deleteCollection(name: string) {

        return this.http.delete(this.url + "collections/" + name,
                                this.getOptions())
            .map((res: Response) => {
                if (res.status != 204) {
                    throw new Error('Impossible to delete collection '
                                    + name + ': ' + res.status);
                }

                // Remove the collection from the list
                for (let i = 0; i < this.collections.length; i++) {
                    if (this.collections[i].name === name) {
                        this.collections.splice(i, 1)
                        break
                    }
                }

                // Warn about deleted collection
                this.setCollection(undefined, CollectionStatus.DELETED)
            })
            .catch(this.handleError);
    }

    // Get the collections list
	getAll() {

        return new Observable<Collection[]>(observer => {
            if (this.collections) {
                observer.next(this.collections)
                return
            }

            let request =  this.http.get(this.url + "collections",
                                         this.getOptions())
                .map(this.extractData)
                .catch(this.handleError);

            request.subscribe(collections => {

                if (collections) {
                    this.collections = collections
                    observer.next(collections)
                }
            })
        });
    }


    // Get the collections references
	getReferences() {

        // Setup cache on the references
        return new Observable(observer => {
            if (this.references) {
                observer.next(this.references)
                return
            }

            let request = this.http.get(this.url + "references")
                .map(this.extractData)
                .catch(this.handleError);

            request.subscribe(references => {
                this.references = references
                observer.next(references)
            })
        })
    }

    newCollectionImport(newImport: any) {

        let collectionUrl = this.getCollectionUrl()
        if (collectionUrl == undefined) {
            return
        }

        return this.http.post(collectionUrl + "/imports",
                              JSON.stringify(newImport),
                              this.getOptions())
            .map((res: Response) => {
                if (res.status != 204) {
                    throw new Error('Impossible to create new import: ' + res.status);
                }

            })
            .catch(this.handleError);
    }

    getCollectionImport() {

        let collectionUrl = this.getCollectionUrl()
        if (collectionUrl == undefined) {
            return
        }

        return new Observable<Imports>(observer => {

            if (this.collectionSelected.imports) {
                observer.next(this.collectionSelected.imports)
                return
            }

            let request = this.http.get(collectionUrl + "/imports")
                .map(this.extractData)
                .catch(this.handleError);

            request.subscribe(imports => {
                console.log(imports)
                this.collectionSelected.imports = imports
                observer.next(imports)
            })
        });
    }

    // Delete an existing collection import
    deleteCollectionImport(name: string) {

        let collectionUrl = this.getCollectionUrl()
        if (collectionUrl == undefined) {
            return
        }

        return this.http.delete(collectionUrl + "/imports/" + name,
                                this.getOptions())
            .map((res: Response) => {
                if (res.status != 204) {
                    throw new Error('Impossible to delete import collection '
                                    + name + ': ' + res.status);
                }
            })
            .catch(this.handleError);
    }


    private extractData(res: Response) {

        if (res.status < 200 || res.status >= 300) {
            throw new Error('Bad response status: ' + res.status);
        }

        // No content to return
        if (res.status === 204) {
            return true
        }

        return res.json();
    }

    private handleError (error: any) {
        let errMsg = error.message || 'Server error';
        console.error(errMsg)
        return Observable.throw(errMsg);
    }
}
