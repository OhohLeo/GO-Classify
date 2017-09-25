import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { Response } from '@angular/http';

@Injectable()
export class ParamsService {

    constructor(private apiService: ApiService) { }


    actionParam(type: string, name: string, param: string, data: string) {

        return new Observable(observer => {

            return this.apiService.put(type + "/" + name + "/" + param)
                .subscribe(rsp => {
                    if (rsp.status != 204) {
                        throw new Error('Error when ' + name + '/' + param
                            + ' ' + type + ': '
                            + rsp.status)
                    }

                    observer.next(rsp)
                })
        })
    }

}
