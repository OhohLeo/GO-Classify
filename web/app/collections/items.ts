import { Item } from './item'

export class Items {

    public items: Item[] = []

    isEmpty(): boolean {
        return (this.items.length === 0)
    }

    hasItem(idx: number): boolean {

        if (idx >= 0 && idx < this.items.length)
            return true

        return false
    }

    addItem(item: Item) {
        this.items.push(item)
    }

    removeItem(idx: number) {

        if (this.hasItem(idx)) {
            this.items.splice(idx, 1)
            return
        }

        console.error("Invalid item id", idx)
    }

    getItem(idx: number): Item {

        if (this.hasItem(idx)) {
            return this.items[idx]
        }

        return undefined
    }

    getIdx(id: string): number {

        for (let idx in this.items) {
            if (this.items[idx].id === id) {
                return +idx
            }
        }

        console.error("invalid buffer item with id", id)
        return -1
    }

}
