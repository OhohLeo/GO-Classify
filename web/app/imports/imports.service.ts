import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ClassifyService } from './../classify.service';
import { Response } from '@angular/http';

export interface ImportInterface {
    setId(id: string)
	getId(): string
    getType(): string
    getParams(): any
	display(): string
    compare(i: ImportInterface): boolean
}

export class Directory implements ImportInterface {
    constructor(public id: string,
				public path: string,
				public isRecursive: boolean) {

        if (isRecursive === undefined) {
            this.isRecursive = false
        }
    }

	setId(id: string) {
		this.id = id;
	}

	getId(): string {
		return this.id
	}

    getType(): string {
        return "directory"
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
        return this.path === i.path
            && this.isRecursive == i.isRecursive
    }
}


@Injectable()
export class ImportsService {

	enableCache: boolean
    imports: Map<string, ImportInterface[]> = new Map<string, ImportInterface[]>()
    configs

	private convertToImportInterface = {
		"directory": function (id: string, params) : ImportInterface {

			if (typeof params != 'object') {
				console.error("Unsupported directory parameters!")
				return undefined
			}

			let path = params['path']
			if (typeof path !=  'string') {
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

    newImport(i: ImportInterface) {

        if (this.hasImport(i)) {
			console.error("Import already existing")
			return
        }

		// Disable cache
		this.enableCache = false

        return this.classifyService.post(
            "imports", {
                "type": i.getType(),
                "params": i.getParams()
            })
            .subscribe(rsp => {

                if (rsp.status != 200) {
                    throw new Error('Error when creating new import: ' + rsp.status);
                }

				let body = rsp.json()

				if (body === undefined && body.id === undefined) {
					throw new Error('Id not found when creating new import!');
				}

				i.setId(body.id)

				// Add new import
                this.addImport(i)
            })
    }

    hasImport(search: ImportInterface): boolean {

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

    addImport(i: ImportInterface) {

		if (this.hasImport(i)) {
			console.error("Already existing " + i.getType())
			return
		}

		if (this.imports.get(i.getType()) === undefined) {
			this.imports.set(i.getType(), [])
		}

        this.imports.get(i.getType()).push(i)
    }

	deleteImport(i: ImportInterface) {

		// Disable cache
		this.enableCache = false

		if (this.hasImport(i) === false) {
			console.error("No existing " + i.getType())
			return
		}

		let urlParams = "?id=" + i.getId()
			+ "&collection=" + this.classifyService.getCollectionName();

		return this.classifyService.delete("imports" + urlParams)
	}

    getImports() {
        return  new Observable(observer => {

			if (this.imports && this.enableCache === true) {
				observer.next(this.imports)
				return
			}

			this.classifyService.get("imports")
				.subscribe(rsp => {

					// TODO Check if all imports exist or not

					// For all imports type
					for (let importType in rsp) {

						let convert = this.convertToImportInterface[importType]
						if (convert === undefined) {
							console.error(
								"Unknown import type '" + importType + "'")
							continue
						}

						for (let id in rsp[importType])	{
							let importConverted = convert(id, rsp[importType][id])
							if (importConverted === undefined)
								continue

							this.addImport(importConverted)
						}
					}

					this.enableCache = true

					observer.next(this.imports)
				})
		})
    }

    getImportsConfig(importType: string) {
        return new Observable(observer => {
            if (this.configs) {
                observer.next(this.configs[importType])
                return
            }

            return this.classifyService.get("imports/config")
                .subscribe(rsp => {
                    this.configs = rsp
                    observer.next(rsp[importType])
                })
        })
    }
}
