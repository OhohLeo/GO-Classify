import { Injectable, NgZone } from '@angular/core';
import { Http, Response, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Rx';
import { Collection, Imports } from './collections/collection';

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
    DELETED,
}

export interface Event {
    event: string
    status: string
    name: string
    data: any
}

declare var EventSource: any

@Injectable()
export class ApiService {

    private url = "http://localhost:1234/"
    private references: any
    private collections: Collection[]

    private onCollectionChange: (collection: Collection,
        status: CollectionStatus) => void

    private zone = new NgZone({ enableLongStackTrace: false })

    public status = WebSocketStatus.NONE

    public collectionSelected: Collection

    constructor(private http: Http) {}

    headers() {
        return new RequestOptions({
            headers: new Headers({ 'Content-Type': 'application/json' })
        })
    }

    get(path: string) {
        return this.http.get(
            this.url + path, this.headers())
            .map(this.extractData)
            .catch(this.handleError);
    }

    post(path: string, data: any) {
        return this.http.post(
            this.url + path, JSON.stringify(data), this.headers())
            .catch(this.handleError);
    }

    put(path: string) {
        return this.http.put(
            this.url + path, this.headers())
            .catch(this.handleError);
    }

    patch(path: string, value: any) {
        return this.http.patch(
            this.url + path, JSON.stringify(value), this.headers())
            .catch(this.handleError);
    }

    delete(path: string) {
        return this.http.delete(
            this.url + path, this.headers())
            .catch(this.handleError);
    }

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

    getCollectionName(): string {

        if (this.collectionSelected == undefined) {
            this.handleError("Select a collection first!")
            return undefined
        }

        return this.collectionSelected.name
    }

    getCollectionUrl(): string {
        return this.url + "collections/" + this.getCollectionName()
    }

    // Create a new collection
    newCollection(collection: Collection) {

        return this.http.post(this.url + "collections",
            JSON.stringify(collection),
            this.headers())
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
            this.headers())
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
            this.headers())
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

    // GET /stream
    getStream() {
        return Observable.create(observer => {

            let eventSource = new EventSource(this.url + "stream")

            eventSource.onmessage =
                event => observer.next(JSON.parse(
                    event.data))

            eventSource.onerror =
                error => {
                    console.error("EVENT SOURCE", error)
                }

            return () => {
                eventSource.close()
            }
        })
    }

    // Get the collections list
    getCollections() {

        return new Observable<Collection[]>(observer => {

			console.log("COLLECTIONS BEFORE !!", this.collections)

            if (this.collections) {
                observer.next(this.collections)
                return
            }

            let request = this.http.get(this.url + "collections",
                this.headers())
                .map(this.extractData)
                .catch(this.handleError);

            request.subscribe(collections => {

				console.log("COLLECTIONS!!", collections)

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
            this.headers())
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

    private handleError(error: any) {
        let errMsg = error.message || 'Server error';
        console.error(errMsg)
        return Observable.throw(errMsg);
    }
}
