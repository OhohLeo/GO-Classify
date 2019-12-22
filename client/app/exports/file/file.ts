import { BaseElement } from '../../base'

export class File extends BaseElement {

    public path: string = ""
    public permissions: string

    constructor(public name: string) {
        super("exports", "file", name);
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

export function Convert2File(name: string, params): File {

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

    let file = new File(name)
    file.path = path
    file.permissions = permissions

    return file
}
