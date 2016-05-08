import {Injectable} from 'angular2/core';
import {Http, Response, RequestOptions, Headers} from 'angular2/http';
import {Observable} from 'rxjs/Rx';
import {Collection} from './collections/collection';

@Injectable()
export class ClassifyService {

    private url = "http://localhost:8080/"
    private references: any
    private collections: Collection[]


    private onChanges: (collection: Collection) => void

    public collectionSelected: Collection

    constructor (private http: Http) {}

    setOnChanges(changes: (collection: Collection) => void) {
        this.onChanges = changes
    }


    selectCollection(collection: Collection) {
        this.collectionSelected = collection
    }

    getOptions() {
        return new RequestOptions({
            headers: new Headers({ 'Content-Type': 'application/json' })
        })
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

                this.onChanges(collection)
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

                // Remove the selected collection
                this.collectionSelected = collection

                this.onChanges(collection)
            })
            .catch(this.handleError);
    }

    // Delete an existing collection
    deleteCollection(name: string) {

        return this.http.delete(this.url + "collections/" + name,
                                this.getOptions())
            .map((res: Response) => {
                if (res.status != 204) {
                    throw new Error('Impossible to modify collection '
                                    + name + ': ' + res.status);
                }

                // Remove the collection from the list
                for (let i = 0; i < this.collections.length; i++) {
                    if (this.collections[i].name === name) {
                        this.collections.splice(i, 1)
                        break
                    }
                }

                // Reset the selected collection
                this.collectionSelected = undefined

                this.onChanges(undefined)
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
        console.error(errMsg);
        return Observable.throw(errMsg);
    }
}
