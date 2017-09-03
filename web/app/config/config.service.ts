import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { BufferService } from './../buffer/buffer.service';
import { Response } from '@angular/http';
import { CfgStringList, StringListEvent } from '../tools/stringlist.component';

export class ConfigBase {
    private cache: boolean

    public enableStore: boolean
    public enableBuffer: boolean
    public bufferSize: number
    public filters = new CfgStringList()
    public separators = new CfgStringList()
    public banned = new CfgStringList()

    init(config) {
        Object.keys(config).forEach((key) => {
            switch (key) {
                case "enableStore":
                    this.enableStore = config[key] ? true : false
                    break
                case "enableBuffer":
                    this.enableBuffer = config[key] ? true : false
                    break
                case "bufferSize":
                    this.bufferSize = config[key]
                    break
                case "filters":
                    this.filters.init(config[key])
                    break
                case "separators":
                    this.separators.init(config[key])
                    break
                case "banned":
                    this.banned.init(config[key])
                    break
                default:
                    console.error("Unhandled configuration '" + key + "'")
            }
        })
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
}

@Injectable()
export class ConfigService {

    // Collections configurations
    public configs: { [key: string]: ConfigBase; } = {}

    constructor(private apiService: ApiService,
        private bufferService: BufferService) { }

    public hasConfig(collection: string): boolean {
        return this.configs[collection] != undefined
    }

    public getConfigs(collection: string) {

        return new Observable(observer => {

            let currentConfs = this.configs[collection]
            if (currentConfs != undefined && currentConfs.hasCache()) {
                observer.next(currentConfs)
                return
            }

            // Ask for the current configuration
            this.apiService.get("collections/" + collection + "/config")
                .subscribe((confs) => {

                    // No config
                    if (confs == null) {
                        console.error("No configuration found for collection "
                            + collection)
                        return
                    }

                    let config = new ConfigBase()
                    config.init(confs)
                    config.enableCache()
                    this.configs[collection] = config

                    observer.next(config)
                })
        })
    }

    public setConfig(collection: string, name: string, action: any, value: any) {

        let currentConfs = this.configs[collection]
        if (currentConfs == undefined) {
            console.error("No configuration found for collection '"
                + collection + "'")
            return
        }

        let body = {
            'name': name,
        }

        body['action'] = action

        switch (action) {
            case "add":
            case "remove":
                body['list'] = value
                break;
            default:
                body['value'] = value
        }

        return new Observable(observer => {

            currentConfs.disableCache();

            this.apiService.patch("collections/" + collection + "/config", body)
                .subscribe((status) => {
                    console.log(status)

                    // On change on buffer size : update buffer list
                    if (name === "bufferSize")
                        this.bufferService.disableCache();

                    observer.next(status)
                })
        })
    }

    public update(collection: string) {
    }

    public onChange(collection: string, event) {

        console.log("onChange", event)

        let name, action, value

        if (event instanceof StringListEvent) {
            name = event.name
            action = event.action
            value = event.list
        } else {
            name = event.target.name
            switch (event.target.type) {
                case "number":
                    value = Number(event.target.value)
                default:
                    value = event.target.value
            }
        }

        console.log("CHANGE", name, action, value);

        let observable = this.setConfig(
            collection, name, action, value)
        if (observable != undefined) {
            observable.subscribe((status) => {
            })
        }
    }
}
