import { ImportBase } from '../imports.service';

export class Directory extends ImportBase {

	public path: string
	public isRecursive: boolean = false

    constructor(public name: string) {
        super("directory", name);
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

export function Convert2Directory(name: string, params): ImportBase {

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

	let directory = new Directory(name)

	directory.path = path
	directory.isRecursive = isRecursive ? true : false

	return directory
}
