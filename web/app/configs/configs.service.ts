import { Injectable } from '@angular/core'
import { Response } from '@angular/http'
import { Observable } from 'rxjs/Rx'
import { ConfigBase } from './config_base'
import { ApiService, Event } from '../api.service'
import { BaseElement } from '../base'

@Injectable()
export class ConfigsService {

    public configs: { [src: string]: { [name: string]: ConfigBase; } } = {}

    constructor(private apiService: ApiService) { }

    public getConfig(src: string, item: BaseElement) {

	let collectionName = this.apiService.getCollectionName()
	if (collectionName == undefined) {
	    console.error("no collection name found")
	    return undefined
	}

	let name = item.name
	
        return new Observable(observer => {

            let needReferences: boolean = true
            let currentCfg = this.configs[src]

            if (currentCfg != undefined) {

                let cfg = currentCfg[name]
                if (cfg != undefined) {
                    needReferences = (cfg.getRefs().length < 1)

                    if (needReferences == false && cfg.hasCache()) {
                        observer.next(cfg)
                        return
                    }
                }
            }

	    let url = src + "/" + name + "/config"
	    let queries = []
	    
	    if (needReferences) {
		queries.push("references")
	    }
	    
	    if (["imports", "exports"].includes(item.getType())) {
		queries.push("collection=" + collectionName) 
	    }

	    if (queries.length > 0) {
		url += "?" + queries.join("&")
	    }
	    
            // Ask for the current configuration
	    this.apiService.get(url)
                .subscribe((res: any) => {

                    // No config
                    if (res == undefined) {
                        console.error("no configuration found for " + src + "/" + name)
                        return
                    }

                    let cfg: ConfigBase

                    // References are expected
                    if (needReferences) {

			console.log(res)
			let references = res["references"]
                        if (references == undefined) {
                            console.error("no references received at " + src + "/" + name)
                            return
                        }

                        cfg = new ConfigBase()

			let forceRefs = []
			let generic = references["generic"]
			if (generic) {
			    forceRefs.push({"name":"generic","type":"struct","childs":generic})
			}
			let specific = references["specific"]
			if (specific) {
			    forceRefs.push({"name":"specific","type":"struct","childs":specific})
			}

			if (forceRefs.length > 0) {
                            cfg.setRefs(forceRefs)          
                        } else {
			    cfg.setRefs(references)
			}
			
                        cfg.setData(res["data"])

                        if (this.configs[src] == undefined) {
                            this.configs[src] = {}
                        }

                        this.configs[src][name] = cfg

                    }
                    // Otherwise needs to update data
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

	let collectionName = this.apiService.getCollectionName()
	if (collectionName == undefined) {
	    console.error("no collection name found")
	    return undefined
	}

	let currentCfg = this.configs[src]
        if (currentCfg == undefined || currentCfg[name] == undefined) {
            console.error("No configuration found for " + src + "/" + name)
            return
        }

        let cfg = currentCfg[name]

        return new Observable(observer => {

            cfg.disableCache();
            console.log(cfg.getData())
            this.apiService.patch(src + "/" + name + "/config?collection=" + collectionName,
				  cfg.getData())
                .subscribe((status) => {
                    console.log(status)

                    observer.next(status)
                })
        })
    }
}
