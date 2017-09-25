import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { BufferService } from './../buffer/buffer.service';
import { Response } from '@angular/http';

export class ImportBase {

    public isRunning: boolean

    constructor(private ref: string, public name: string) { }

    getRef(): string {
        if (this.ref === undefined)
            throw new Error("attribute 'ref' should be defined!")

        return this.ref
    }

    getParams(): any {
        throw new Error("method 'getParams' should be defined!")
    }

    display(): string {
        throw new Error("method 'display' should be defined!")
    }

    compare(i: ImportBase): boolean {
        if (this.ref === undefined)
            throw new Error("attribute 'ref' should be defined!")

        if (this.ref != i.getRef())
            return false

        return true
    }
}

@Injectable()
export class ImportsService {

    enableCache: boolean
    imports: Map<string, ImportBase[]> = new Map<string, ImportBase[]>()
    importsByName: Map<string, ImportBase> = new Map<string, ImportBase>()
    configs: any
    updateList: any

    private eventObservers = {}
    private convertToImport = {}

    constructor(private apiService: ApiService,
        private bufferService: BufferService) { }

    // Set update import list function
    setUpdateList(updateList: any) {
        this.updateList = updateList;
    }

    addConvertToImport(name: string, callback) {
        this.convertToImport[name] = callback
    }

    // Refresh the import list
    private update() {
        if (this.updateList != undefined)
            this.updateList()
    }

    // Check if import does exist
    hasImport(search: ImportBase): boolean {
        return this.importsByName.get(search.name) != undefined
    }

    // Check if import does exist
    hasSameImport(search: ImportBase): boolean {
        let imports = this.imports.get(search.getRef())
        if (imports === undefined) {
            return false
        }

        for (let i of imports) {
            if (i.compare(search)) {
                return true
            }
        }

        return false
    }

    private add(i: ImportBase) {

        // Store imports by name
        this.importsByName.set(i.name, i)

        // Store imports by ref
        if (this.imports.get(i.getRef()) === undefined) {
            this.imports.set(i.getRef(), [])
        }

        this.imports.get(i.getRef()).push(i)
    }

    addImport(name: string, i: ImportBase, onParams: any, onSuccess: any) {

        // Disable cache
        this.enableCache = false

        if (this.hasSameImport(i)) {
            console.error("Already existing " + i.getRef())
            return
        }

        return this.apiService.post(
            "imports", {
                "name": name,
                "ref": i.getRef(),
                "params": i.getParams(),
                "collections": [this.apiService.getCollectionName()],
            })
            .subscribe(rsp => {

                if (rsp.status != 200) {
                    throw new Error('Error when adding new import: ' + rsp.status);
                }

                let body = rsp.json()

                if (body === undefined || body.name === undefined) {

                    if (onParams !== undefined && onParams(body))
                        return

                    throw new Error('Name not found when adding new import!')
                }

                this.add(i)

                this.update()

                if (onSuccess !== undefined) {
                    onSuccess(i)
                }
            })
    }

    private delete(i: ImportBase) {

        // Delete import by name
        this.importsByName.delete(i.name)

        // Delete import by ref
        let importList = this.imports.get(i.getRef())
        for (let idx in importList) {
            let importItem = importList[idx]
            if (importItem.name === i.name) {
                importList.splice(+idx, 1)
                break;
            }
        }

        // Remove import refs with no imports
        if (importList.length == 0) {
            this.imports.delete(i.getRef())
        }
    }

    deleteImport(i: ImportBase) {

        // Disable cache
        this.enableCache = false

        if (this.hasImport(i) === false) {
            console.error("No existing " + i.getRef())
            return
        }

        let urlParams = "?name=" + i.name
            + "&collection=" + this.apiService.getCollectionName();

        return this.apiService.delete("imports" + urlParams)
            .subscribe(rsp => {

                if (rsp.status != 204) {
                    throw new Error('Error when deleting import: ' + rsp.status)
                }

                // Delete import
                this.delete(i)

                this.update()
            })
    }

    startImport(i: ImportBase) {
        return this.actionImport(true, i)
    }

    stopImport(i: ImportBase) {
        return this.actionImport(false, i)
    }

    actionImport(isStart: boolean, i: ImportBase) {

        if (this.hasImport(i) === false) {
            console.error("No existing " + i.getRef())
            return
        }

        let action = isStart ? "start" : "stop"
        let urlParams = "?name=" + i.name
            + "&collection=" + this.apiService.getCollectionName();

        return this.apiService.put("imports/" + action + urlParams)
            .subscribe(rsp => {
                if (rsp.status != 204) {
                    throw new Error('Error when ' + action + ' import: ' + rsp.status)
                }

                if (isStart)
                    this.bufferService.disableCache();
            })
    }

    // Ask for current imports list
    getImports() {
        return new Observable(observer => {

            // Returns the cache if the list should not have changed
            if (this.imports && this.enableCache === true) {
                observer.next(this.imports)
                return
            }

            // Ask for the current list
            this.apiService.get("imports").subscribe(rsp => {

                // Init the import lists
                this.imports = new Map<string, ImportBase[]>()
                this.importsByName = new Map<string, ImportBase>()

                for (let importRef in rsp) {

                    let convert = this.convertToImport[importRef]
                    if (convert === undefined) {
                        console.error(
                            "Unknown import ref '" + importRef + "'")
                        continue
                    }

                    for (let importName in rsp[importRef]) {
                        let i = convert(importName, rsp[importRef][importName])
                        if (i === undefined)
                            continue

                        this.add(i)
                    }
                }

                this.enableCache = true

                observer.next(this.imports)
            })
        })
    }

    // Ask for current import config list
    getImportsConfig(importRef: string) {
        return new Observable(observer => {

            // Import config list should not change a lot
            if (this.configs) {
                observer.next(this.configs[importRef])
                return
            }

            // Ask for the current import config list
            return this.apiService.get("imports/config")
                .subscribe(rsp => {

                    // Store as cache the current import config list
                    this.configs = rsp

                    // Return the import config list
                    observer.next(rsp[importRef])
                })
        })
    }

    subscribeEvents(name: string): Observable<Event> {

        if (this.eventObservers[name] != undefined) {
            console.error("Already existing observer", name)
            return;
        }

        return Observable.create(observer => {

            // Initialisation de l'observer
            this.eventObservers[name] = observer

            return () => delete this.eventObservers[name]
        })
    }

    addEvent(event: Event) {
        for (let name in this.eventObservers) {
            this.eventObservers[name].next(event)
        }
    }
}
