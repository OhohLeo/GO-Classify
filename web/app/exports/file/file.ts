import { ExportBase } from '../exports.service';

export class File extends ExportBase {

	public path: string
	public permissions: string

    constructor(public id: string) {
        super("file", id);
    }

    getParams(): any {
        return {
			"path": this.path,
			"permissions": this.permissions
		}
    }

    compare(i: File): boolean {
        return super.compare(i)
			&& this.path === i.path
			&& this.permissions === i.permissions
    }
}

export function Convert2File(id: string, params): ExportBase {

	if (typeof params != 'object') {
		console.error("Unsupported file parameters!")
		return undefined
	}

	let path = params['path']
	if (typeof path != 'string') {
		console.error("Unsupported 'path' file parameters!")
		return undefined
	}

	let permissions = params['permissions']
	if (typeof permissions != 'string') {
		console.error("Unsupported 'permissions' file parameters!")
		return undefined
	}

	let file = new File(id)
	file.path = path
	file.permissions = permissions

	return file
}
