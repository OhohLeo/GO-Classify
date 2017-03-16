import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { Response } from '@angular/http';
import { CfgStringList } from './stringlist.component';

export class ConfigBase {
    private cache: boolean

    public bufferSize: number
    public filters = new CfgStringList()
    public separators = new CfgStringList()
    public banned = new CfgStringList()

    init(config) {

        Object.keys(config).forEach((key) => {
            switch (key) {
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

    hasChanged(name: string, rcv: any): boolean {
        switch (name) {
            case "bufferSize":
                return (this.bufferSize != rcv)
            case "filters":
                return this.filters.hasChanged(rcv)
            case "separators":
                return this.separators.hasChanged(rcv)
            case "banned":
                return this.banned.hasChanged(rcv)
            default:
                console.error("Unhandled configuration '" + name + "'")
        }

        return false;
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
    public configs: Map<string, ConfigBase> = new Map<string, ConfigBase>()

    constructor(private apiService: ApiService) { }

    public hasConfig(collection: string): boolean {
        return this.configs.get(collection) != undefined
    }

    public getConfigs(collection: string) {

        return new Observable(observer => {

            let currentConfs = this.configs.get(collection)
            if (currentConfs != undefined && currentConfs.hasCache()) {
                observer.next(currentConfs)
                return
            }

            // Ask for the current configuration
            this.apiService.get("collections/" + collection + "/config")
                .subscribe((confs) => {

                    // No config
                    if (confs == null) {
                        console.log("No configuration found for collection "
                            + collection)
                        return
                    }

                    let config = new ConfigBase()
                    config.init(confs)
                    config.enableCache()
                    this.configs.set(collection, config)

                    observer.next(config)
                })
        })
    }

    public setConfig(collection: string, name: string, action: any, value: any) {

        let currentConfs = this.configs.get(collection)
        if (currentConfs == undefined) {
            console.log("No configuration found for collection '" + collection + "'")
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
                    observer.next(status)
                })
        })
    }

    public update(collection: string) {
    }
}
