import { Injectable } from '@angular/core';
import { ClassifyService } from './../classify.service';
import { Response } from '@angular/http';
import { Directory } from './directory.component';

export interface ImportInterface {
    getType(): string
    getParams(): any
    compare(i: ImportInterface): boolean
}

@Injectable()
export class ImportsService {

    imports: Array<ImportInterface> = [];
    freshness: number = 300

    constructor(private classifyService: ClassifyService) { }

    newImport(i: ImportInterface) {

        if (this.hasImport(i)) {
            console.error("Import already existing")
            return
        }

        return this.classifyService.post(
            "imports", {
                "type": i.getType(),
                "params": i.getParams()
            })
            .subscribe(rsp => {

                if (rsp.status != 204) {
                    throw new Error('Impossible to create new import: ' + rsp.status);
                }

                this.addImport(i)
            })
    }

    getImports() {
        return this.classifyService.get("imports")
            .subscribe(rsp => {
                console.log(rsp)
                for (let i of rsp) {
                    this.addImport(rsp[i])
                }
            })
    }

    getImportsConfig() {
        return this.classifyService.get("imports/config")
            .subscribe(rsp => {
                console.log(rsp)
            })
    }

    addImport(i: ImportInterface) {
        this.imports.push(i)
    }

    hasImport(search: ImportInterface): boolean {
        for (let i of this.imports) {

            if (i.getType() === search.getType()
                && i.compare(search)) {
                return true
            }
        }

        return false
    }
}
