import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { Response } from '@angular/http';

@Injectable()
export class ParamsService {

    constructor(private apiService: ApiService) { }


    actionParam(type: string, name: string, param: string, data: any) {

        return new Observable(observer => {

            return this.apiService.putWithData(
                type + "/" + name + "/params/" + param, data)
                .subscribe(rsp => {
                    if (rsp.status != 200) {
                        throw new Error('Error when ' + name + '/params/' + param
                            + ' ' + type + ': '
                            + rsp.status)
                    }

                    observer.next(rsp.json())
                })
        })
    }

}
