import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { Response } from '@angular/http';

@Injectable()
export class TweaksService {

    constructor(private apiService: ApiService) { }


    getTweak(type: string, name: string) {

        return new Observable(observer => {

            return this.apiService.putWithData(
                type + "/" + name + "/tweaks", data)
                .subscribe(rsp => {
                    if (rsp.status != 200) {
                        throw new Error('Error when ' + name + '/tweaks'
                            + ' ' + type + ': '
                            + rsp.status)
                    }

                    observer.next(rsp.json())
                })
        })
    }

    actionTweak(type: string, name: string, data: any) {

        return new Observable(observer => {

            return this.apiService.putWithData(
                type + "/" + name + "/tweaks", data)
                .subscribe(rsp => {
                    if (rsp.status != 200) {
                        throw new Error('Error when ' + name + '/tweaks'
                            + ' ' + type + ': '
                            + rsp.status)
                    }

                    observer.next(rsp.json())
                })
        })
    }

}
