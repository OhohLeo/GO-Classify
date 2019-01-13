import { Injectable } from '@angular/core'
import { Response } from '@angular/http'
import { Observable } from 'rxjs/Rx'
import { ConfigBase } from './config_base'
import { ApiService, Event } from './../api.service'

@Injectable()
export class ConfigsService {

    public configs: { [src: string]: { [name: string]: ConfigBase; } } = {}

    constructor(private apiService: ApiService) { }

    public getConfig(src: string, name: string) {

	let collectionName = this.apiService.getCollectionName()
	if (collectionName == undefined) {
	    console.error("no collection name found")
	    return undefined
	}
		
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

            // Ask for the current configuration
            this.apiService.get(src + "/" + name + "/config"
		+ "?collection=" + collectionName
		+ (needReferences ? "&references" : ""))
                .subscribe((res: any) => {

                    // No config
                    if (res == undefined) {
                        console.error("No configuration found for "
                            + src + "/" + name)
                        return
                    }

                    let cfg: ConfigBase

                    // Refs are expected
                    if (needReferences) {

			let references = res["references"]
                        if (references == undefined) {
                            console.error("No refs received at " + src + "/" + name)
                            return
                        }

			let generic = res["generic"]
			if (generic == undefined) {
                            console.error("No generic data received at " + src + "/" + name)
                            return
                        }
			
                        cfg = new ConfigBase()
                        cfg.setRefs([
			    { "name": "generic", "type": "struct", "childs": references["generic"] },
			    { "name": "specific", "type": "struct", "childs": references["specific"] },
			])
                        cfg.setData({
			    "generic": res["generic"],
                            "specific": res["specific"],
			})

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
