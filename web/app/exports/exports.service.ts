import { Injectable } from '@angular/core'
import { Observable } from 'rxjs/Rx'
import { ApiService, Event } from './../api.service'
import { BufferService } from './../buffer/buffer.service'
import { Response } from '@angular/http'

import { Convert2File } from './file/file'

import { BaseElement } from '../base'

@Injectable()
export class ExportsService {

    enableCache: boolean
    exports: Map<string, BaseElement[]> = new Map<string, BaseElement[]>()
    exportsByName: Map<string, BaseElement> = new Map<string, BaseElement>()
    configs: any
    updateList: any

    private eventObservers = {}
    public convertToExport: { [index:string]: (string, any) => BaseElement } = {
        "file": Convert2File,
    }

    constructor(private apiService: ApiService,
		private bufferService: BufferService) {}

    // Set update export list function
    setUpdateList(updateList: any) {
        this.updateList = updateList;
    }

    // Refresh the export list
    private update() {
        if (this.updateList != undefined)
            this.updateList()
    }

    // Check if export does exist
    hasExport(search: BaseElement): boolean {
        return this.hasSameExportName(search.name)
    }

    // Check if export does exist
    hasSameExportName(name: string): boolean {
        return this.exportsByName.get(name) != undefined
    }

    // Check if export does exist
    hasSameExport(search: BaseElement): boolean {
        let exports = this.exports.get(search.getRef())
        if (exports === undefined) {
            return false
        }

        for (let i of exports) {
            if (i.compare(search)) {
                return true
            }
        }

        return false
    }

    private add(i: BaseElement) {

        // Store exports by name
        this.exportsByName.set(i.name, i)

        // Store exports by ref
        if (this.exports.get(i.getRef()) === undefined) {
            this.exports.set(i.getRef(), [])
        }

        this.exports.get(i.getRef()).push(i)
    }

    addExport(i: BaseElement, onParams: any, onSuccess: any) {

        // Disable cache
        this.enableCache = false

        if (this.hasSameExport(i)) {
            console.error("Already existing " + i.getRef())
            return
        }

	let name = i.getName()
	if (this.hasSameExportName(name)) {
	    console.error("Already existing name " + name)
            return
	}

        return this.apiService.post(
            "exports", {
		"name": name,
                "ref": i.getRef(),
                "params": i.getParams(),
                "collections": [this.apiService.getCollectionName()],
            })
            .subscribe(rsp => {

                if (rsp.status != 200) {
                    throw new Error('Error when adding new export: ' + rsp.status);
                }

                let body = rsp.json()

                if (body === undefined && body.name === undefined) {

		    if (onParams !== undefined && onParams(body))
			return

                    throw new Error('Name not found when adding new export!');
                }

                this.add(i)

                this.update()

		if (onSuccess !== undefined) {
		    onSuccess(i)
		}
            })
    }

    private delete(i: BaseElement) {

        // Delete export by name
        this.exportsByName.delete(i.name)

        // Delete export by ref
        let exportList = this.exports.get(i.getRef())
        for (let idx in exportList) {
            let exportItem = exportList[idx]
            if (exportItem.name === i.name) {
                exportList.splice(+idx, 1)
                break;
            }
        }

        // Remove export refs with no exports
        if (exportList.length == 0) {
            this.exports.delete(i.getRef())
        }
    }

    deleteExport(i: BaseElement) {

        // Disable cache
        this.enableCache = false

        if (this.hasExport(i) === false) {
            console.error("No existing " + i.getRef())
            return
        }

        let urlParams = "?name=" + i.name
            + "&collection=" + this.apiService.getCollectionName();

        return this.apiService.delete("exports" + urlParams)
            .subscribe(rsp => {

                if (rsp.status != 204) {
                    throw new Error('Error when deleting export: ' + rsp.status)
                }

                // Delete export
                this.delete(i)

                this.update()
            })
    }

    forceExport(e: BaseElement) {
        return this.actionExport(true, e)
    }

    stopExport(e: BaseElement) {
        return this.actionExport(false, e)
    }

    actionExport(isForce: boolean, e: BaseElement) {

        if (this.hasExport(e) === false) {
            console.error("No existing " + e.getRef())
            return
        }

        let action = isForce ? "force" : "stop"
        let urlParams = "?name=" + e.name
            + "&collection=" + this.apiService.getCollectionName();

        return this.apiService.put("exports/" + action + urlParams)
            .subscribe(rsp => {
                if (rsp.status != 204) {
                    throw new Error(' Error when '
				    + action + ' export: ' + rsp.status)
                }

                if (isForce)
                    this.bufferService.disableCache();
            })
    }

    // Ask for current exports list
    getExports() {
        return new Observable(observer => {

            // Returns the cache if the list should not have changed
            if (this.exports && this.enableCache === true) {
                observer.next(this.exports)
                return
            }

            // Ask for the current list
            this.apiService.get("exports").subscribe(rsp => {

                // Init the export lists
                this.exports = new Map<string, BaseElement[]>()
                this.exportsByName = new Map<string, BaseElement>()

                for (let exportRef in rsp) {

                    let convert = this.convertToExport[exportRef]
                    if (convert === undefined) {
                        console.error(
                            "Unknown export ref '" + exportRef + "'")
                        continue
                    }

                    for (let exportName in rsp[exportRef]) {
                        let i = convert(exportName, rsp[exportRef][exportName])
                        if (i === undefined)
                            continue

                        this.add(i)
                    }
                }

                this.enableCache = true

                observer.next(this.exports)
            })
        })
    }

    // Ask for current export config list
    getExportsConfig(exportRef: string) {
        return new Observable(observer => {

            // Export config list should not change a lot
            if (this.configs) {
                observer.next(this.configs[exportRef])
                return
            }

            // Ask for the current export config list
            return this.apiService.get("exports/config")
                .subscribe(rsp => {

                    // Store as cache the current export config list
                    this.configs = rsp

                    // Return the export config list
                    observer.next(rsp[exportRef])
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
