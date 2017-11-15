import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { Response } from '@angular/http';
import { CfgStringList, StringListEvent } from '../tools/stringlist.component';

export class ConfigRef {

    public name: string
    public type: string
    public comments: string
    public data: any

    public parent: ConfigRef
    public childs: ConfigRef[] = []
    private childsByName: { [name: string]: ConfigRef } = {}

    public callbacks: { [name: string]: (newValue: any) => void } = {}

    constructor(ref: any, parent?: ConfigRef) {

        this.name = ref["name"]
        this.type = ref["type"]
        this.comments = ref["comments"]

        this.parent = parent

        let childs

        // Handle arrays
        if (Array.isArray(ref)) {
            childs = ref
        }
        else {
            childs = ref["childs"]
        }

        if (childs != undefined && Array.isArray(childs)) {
            for (let idx in childs) {
                let ref = new ConfigRef(childs[idx], this)
                this.childs.push(ref)
                this.childsByName[ref.name] = ref
            }
        }
        // Handle map type
        else if (this.type == "map") {

            let map = ref["map"]
            for (let keyIdx in map) {

                let keyRef = new ConfigRef({
                    name: keyIdx,
                    type: "key",
                    childs: map[keyIdx]
                }, this)

                this.childs.push(keyRef)
                this.childsByName[keyRef.name] = keyRef
            }
        }
    }

    subscribeDataChange(name: string, cb: (newValue: any) => void) {
        this.callbacks[name] = cb
    }

    unsubscribeDataChange(name: string) {
        delete this.callbacks[name]
    }

    getFullName(): string {

        if (this.parent != undefined)
            return this.parent.getFullName() + "-" + this.name

        return this.name
    }

    getPathName(): string {

        if (this.parent != undefined)
            return this.parent.getFullName()

        return ""
    }

    getData(): any {

        let data = {}

        for (let ref of this.childs) {

            console.log(ref.type, ref.name)

            switch (ref.type) {
                case "stringlist":
                    if (ref.data instanceof CfgStringList) {
                        data[ref.name] = ref.data.stringlist
                    } else if (Array.isArray(ref.data)) {
                        data[ref.name] = ref.data
                    } else {
                        data[ref.name] = []
                    }
                    break;
                case "map":
                case "key":
                case "struct":
                    data[ref.name] = ref.getData()
                    break;

                default:
                    data[ref.name] = ref.data
            }
        }
        console.log(data)
        return data
    }

    setData(data: any) {

        switch (this.type) {
            case "stringlist":
                this.data = Array.isArray(data) ? data : []
                return;

            case "key":
            case "map":
            case "struct":

                if (typeof data != "object") {
                    console.error("struct expect object", this, data)
                    return;
                }

                for (let [key, value] of Object.entries(data)) {

                    let ref = this.childsByName[key]
                    if (ref === undefined) {
                        console.error("Config child value '" + key
                            + "' not found in", this.name)
                        continue
                    }

                    ref.setData(value)
                }

                return
        }

        // When values are different :
        if (data !== this.data) {
            for (let name in this.callbacks) {
                this.callbacks[name](data)
            }
        }

        this.data = data
    }
}

export class ConfigBase {

    private refs: ConfigRef[] = []
    private refsByName: { [name: string]: ConfigRef } = {}

    private cache: boolean
    private onChange: { [name: string]: (newData: any) => void } = {}

    setRefs(refs: any): boolean {

        if (Array.isArray(refs) == false) {
            console.error("Expected config refs array, received", refs)
            return false
        }

        for (let idx in refs) {
            let ref = new ConfigRef(refs[idx])

            console.log("SET REFS", ref)

            this.refs.push(ref)
            this.refsByName[ref.name] = ref
        }

        return true
    }

    getRefs(): ConfigRef[] {
        return this.refs
    }

    getRef(name: string): ConfigRef {
        return this.refsByName[name]
    }

    hasCache(): boolean {
        return this.cache
    }

    enableCache() {
        this.cache = true
    }

    disableCache() {
        this.cache = false
    }

    getData(): any {

        let data = {}

        for (let ref of this.refs) {

            switch (ref.type) {
                case "map":
                    console.log("GET DATA MAP", ref.name, ref)
                case "struct":
                    data[ref.name] = ref.getData()
                    break;

                default:
                    data[ref.name] = ref.data
            }
        }

        return data
    }

    setData(data: any): boolean {

        if (typeof data != "object") {
            console.error("Expected config object, received", data)
            return false
        }

        for (let [key, value] of Object.entries(data)) {

            let ref = this.refsByName[key]
            if (ref === undefined) {
                console.error("Config value '" + key + "' not found")
                continue
            }

            ref.setData(value)
        }

        return true
    }
}

@Injectable()
export class ConfigsService {

    public configs: { [src: string]: { [name: string]: ConfigBase; } } = {}

    constructor(private apiService: ApiService) { }

    public getConfigs(src: string, name: string) {

        return new Observable(observer => {

            let needRefs: boolean = true
            let currentCfg = this.configs[src]

            if (currentCfg != undefined) {

                let cfg = currentCfg[name]
                if (cfg != undefined) {
                    needRefs = (cfg.getRefs().length < 1)

                    if (needRefs == false && cfg.hasCache()) {
                        observer.next(cfg)
                        return
                    }
                }
            }

            // Ask for the current configuration
            this.apiService.get(src + "/" + name + "/config"
                + (needRefs ? "?refs" : ""))
                .subscribe((res: any) => {

                    // No config
                    if (res == undefined) {
                        console.error("No configuration found for "
                            + src + "/" + name)
                        return
                    }

                    let cfg: ConfigBase

                    // Refs are expected
                    if (needRefs) {

                        if (res["refs"] == undefined) {
                            console.error("No refs received at "
                                + src + "/" + name)
                            return
                        }

                        if (res["data"] == undefined) {
                            console.error("No data received at "
                                + src + "/" + name)
                            return
                        }

                        cfg = new ConfigBase()
                        cfg.setRefs(res["refs"])
                        cfg.setData(res["data"])

                        if (this.configs[src] == undefined) {
                            this.configs[src] = {}
                        }

                        this.configs[src][name] = cfg

                    }
                    // Needs to update data
                    else {
                        cfg = currentCfg[name]
                        cfg.setData(res)
                    }

                    cfg.enableCache()

                    observer.next(cfg)
                })
        })
    }

    public setConfig(src: string, name: string) {

        let currentCfg = this.configs[src]
        if (currentCfg == undefined || currentCfg[name] == undefined) {
            console.error("No configuration found for " + src + "/" + name)
            return
        }

        let cfg = currentCfg[name]

        return new Observable(observer => {

            cfg.disableCache();
            console.log(cfg.getData())
            this.apiService.patch(src + "/" + name + "/config", cfg.getData())
                .subscribe((status) => {
                    console.log(status)

                    observer.next(status)
                })
        })
    }
}
