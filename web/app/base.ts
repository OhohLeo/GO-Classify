import { Md5 } from 'ts-md5/dist/md5';

export class BaseElement {

    public isRunning: boolean

    constructor(private typ: string, private ref: string, public name: string) { }

    getType(): string {
	return this.typ
    }

    getTypeRef(): string {
	return this.typ + "/" + this.ref
    }

    getName(): string {
	return String(Md5.hashStr(JSON.stringify(this.getParams())))
    }

    getRef(): string {
        if (this.ref === undefined)
            throw new Error("attribute 'ref' should be defined!")

        return this.ref
    }

    getParams(): any {
        console.error("[BASE] method 'getParams' should be defined!")
    }

    display(): string {
        throw new Error("method 'display' should be defined!")
    }

    compare(i: BaseElement): boolean {
        if (this.ref === undefined)
            throw new Error("attribute 'ref' should be defined!")

        if (this.ref != i.getRef())
            return false

        return true
    }
}
