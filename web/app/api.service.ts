import { Injectable, NgZone } from '@angular/core';
import { Http, Response, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Rx';

import { ReferencesService } from './references/references.service'

import { Collection, Imports } from './collection/collection';
import { Item } from './items/item'

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

    private onCollectionChange: (collection: Collection,
				 status: CollectionStatus) => void

    private zone = new NgZone({ enableLongStackTrace: false })

    public status = WebSocketStatus.NONE

    public collectionSelected: Collection

    constructor(private http: Http,
		private referencesService: ReferencesService) { }

    headers() {
        return new RequestOptions({
            headers: new Headers({
                'Content-Type': 'application/json'
            })
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
            this.url + this.changePath(path), this.headers())
            .catch(this.handleError);
    }

    putWithData(path: string, data: any) {
        return this.http.put(
            this.url + this.changePath(path),
            JSON.stringify(data), this.headers())
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

    changePath(path: string): string {
        return path.replace("collections/:name",
			    "collections/" + this.collectionSelected.name)
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

    getCollectionReferences() {
	let collectionUrl = this.getCollectionUrl()
        if (collectionUrl == undefined) {
            return
        }

        return new Observable(observer => {
	    let references = this.referencesService.getReference(this.collectionSelected)
	    if (references != undefined) {
		observer.next(references)
		return
	    }

	    let request = this.http.get(collectionUrl + "/references")
                .map(this.extractData)
                .catch(this.handleError);

	    request.subscribe((src) => {
		let references = this.referencesService.setReference(
		    this.collectionSelected, src["datas"])
		console.log("[COLLECTION REFERENCES] OK", references)
		observer.next(references)
	    })
        });
    }

    getReferencesFromSrc(src: string, name: string) {
        // Setup cache on the references
        return new Observable(observer => {
            let request = this.http.get(this.url + "/" + src + "/" + name + "/references")
                .map(this.extractData)
                .catch(this.handleError);

            request.subscribe(references => {
                observer.next(references)
            })
        })
    }

    // Get the classify references
    getReferences() {

        // Setup cache on the references
        return new Observable((observer) => {
	    let references = this.referencesService.getReferences()
            if (references) {
                observer.next(references)
                return
            }

            let request = this.http.get(this.url + "references")
                .map(this.extractData)
                .catch(this.handleError);

            request.subscribe((src) => {
                let references = this.referencesService.setReferences(src)
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

            let imports = this.collectionSelected.imports
            if (imports) {
                observer.next(imports)
                return
            }

            let request = this.http.get(collectionUrl + "/imports")
                .map(this.extractData)
                .catch(this.handleError);

            request.subscribe((imports) => {
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

    getIconUrl(item: Item, size?: string): string {
	return this.getItemUrl(item.id, item.getIconUrl(size))
    }

    getItemUrl(id: string, content?: string): string {

	let url =  this.getCollectionUrl() + "/items/" + id
	if (content == "") {
	    return url
	}

	return url + "?content=" + encodeURIComponent(content)
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
		console.error("CLOSE STREAM!")
                eventSource.close()
            }
        })
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
