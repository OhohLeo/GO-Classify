import { ExportBase } from '../exports.service';

export class File extends ExportBase {

    constructor(public id: string) {
        super("file", id);
    }

    getParams(): any {
        return {}
    }

    compare(i: File): boolean {
        return super.compare(i)
    }
}

export function Convert2File(id: string, params): ExportBase {

	if (typeof params != 'object') {
		console.error("Unsupported file parameters!")
		return undefined
	}

	let file = new File(id)

	return file
}
