import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ClassifyService } from './../classify.service';
import { Response } from '@angular/http';

export class ImportBase {

    constructor(private type: string, public id: string) { }

    setId(id: string) {
        this.id = id;
    }

    getId(): string {
        return this.id
    }

    getType(): string {
        if (this.type === undefined)
            throw new Error("attribute 'type' should be defined!")

        return this.type
    }

    getParams(): any {
        throw new Error("method 'getParams' should be defined!")
    }

    display(): string {
        throw new Error("method 'display' should be defined!")
    }

    compare(i: ImportBase): boolean {
        if (this.type === undefined)
            throw new Error("attribute 'type' should be defined!")

        if (this.type != i.getType())
            return false

        return true
    }
}

export class Directory extends ImportBase {

    constructor(public id: string,
        public path: string,
        public isRecursive: boolean) {

        super("directory", id);

        if (isRecursive === undefined) {
            this.isRecursive = false
        }
    }

    getParams(): any {
        return {
            "path": this.path,
            "is_recursive": this.isRecursive ? true : false
        }
    }

    display(): string {
        return this.path.concat(this.isRecursive == true ? "/**" : "")
    }

    compare(i: Directory): boolean {
        return super.compare(i)
            && this.path === i.path
            && this.isRecursive == i.isRecursive
    }
}


@Injectable()
export class ImportsService {

    enableCache: boolean
    imports: Map<string, ImportBase[]> = new Map<string, ImportBase[]>()
    configs: any
    updateList: any

    private convertToImportBase = {
        "directory": function (id: string, params): ImportBase {

            if (typeof params != 'object') {
                console.error("Unsupported directory parameters!")
                return undefined
            }

            let path = params['path']
            if (typeof path != 'string') {
                console.error("Unsupported 'path' directory parameters!")
                return undefined
            }

            let isRecursive = params['is_recursive']
            if (isRecursive !== undefined && typeof isRecursive != 'boolean') {
                console.error("Unsupported 'is_recursive' directory parameters!")
                return undefined
            }

            return new Directory(id, path, isRecursive)
        }
    };


    constructor(private classifyService: ClassifyService) { }

    // Set new update import list function
    setUpdateList(updateList: any) {
        this.updateList = updateList;
    }

    // Call update list function to refresh the import list
    update() {
        if (this.updateList != undefined)
            this.updateList()
    }

    hasImport(search: ImportBase): boolean {

        let imports = this.imports.get(search.getType())
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

    addImport(i: ImportBase) {

        // Disable cache
        this.enableCache = false

        if (this.hasImport(i)) {
            console.error("Already existing " + i.getType())
            return
        }

        return this.classifyService.post(
            "imports", {
                "type": i.getType(),
                "params": i.getParams()
            })
            .subscribe(rsp => {

                if (rsp.status != 200) {
                    throw new Error('Error when adding new import: ' + rsp.status);
                }

                let body = rsp.json()

                if (body === undefined && body.id === undefined) {
                    throw new Error('Id not found when adding new import!');
                }

                i.setId(body.id)

                // Add new import
                if (this.imports.get(i.getType()) === undefined) {
                    this.imports.set(i.getType(), [])
                }

                this.imports.get(i.getType()).push(i)

                this.update()
            })
    }

    deleteImport(i: ImportBase) {

        // Disable cache
        this.enableCache = false

        if (this.hasImport(i) === false) {
            console.error("No existing " + i.getType())
            return
        }

        let urlParams = "?id=" + i.getId()
            + "&collection=" + this.classifyService.getCollectionName();

        return this.classifyService.delete("imports" + urlParams)
            .subscribe(rsp => {

                if (rsp.status != 204) {
                    throw new Error('Error when deleting import: ' + rsp.status)
                }

                // Delete import
                let importList = this.imports.get(i.getType())
                for (let idx in importList) {
                    let importItem = importList[idx]
                    if (importItem.id === i.getId()) {
                        importList.splice(+idx, 1)
                        break;
                    }
                }

                // Remove import types with no imports
                if (importList.length == 0) {
                    this.imports.delete(i.getType())
                }

                this.update()
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
            this.classifyService.get("imports").subscribe(rsp => {

                // Init the import list
                this.imports = new Map<string, ImportBase[]>()

                for (let importType in rsp) {

                    let convert = this.convertToImportBase[importType]
                    if (convert === undefined) {
                        console.error(
                            "Unknown import type '" + importType + "'")
                        continue
                    }

                    for (let importId in rsp[importType]) {
                        let i = convert(importId, rsp[importType][importId])
                        if (i === undefined)
                            continue

                        if (this.imports.get(i.getType()) === undefined) {
                            this.imports.set(i.getType(), [])
                        }

                        this.imports.get(i.getType()).push(i)
                    }
                }

                this.enableCache = true

                observer.next(this.imports)
            })
        })
    }

    // Ask for current import config list
    getImportsConfig(importType: string) {
        return new Observable(observer => {

            // Import config list should not change a lot
            if (this.configs) {
                observer.next(this.configs[importType])
                return
            }

            // Ask for the current import config list
            return this.classifyService.get("imports/config")
                .subscribe(rsp => {

                    // Store as cache the current import config list
                    this.configs = rsp

                    // Return the import config list
                    observer.next(rsp[importType])
                })
        })
    }
}
