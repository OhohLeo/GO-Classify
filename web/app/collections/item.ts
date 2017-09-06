export class Item {

    public id: string
	public ref: string
	public data: any

	constructor(item) {
		this.id = item["id"]
		this.ref = item["ref"]
		this.data = item["data"]
	}

	getRef() {
		return this.ref
	}
}
