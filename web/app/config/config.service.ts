import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { ApiService, Event } from './../api.service';
import { Response } from '@angular/http';

export class ConfigBase {

	public numberData: Map<string, number> = new Map<string, number>()
	public stringListData: Map<string, string[]> = new Map<string, string[]>()

	setNumber(name: string, item: number)
	{
		this.numberData.set(name, item);
	}

	initStringList(name: string, items: string[])
	{
		this.stringListData.set(name, items);
	}

	addToStringList(name: string, item: string)
	{
		// Check if the item is already in the list

		this.stringListData.set(name, items);
	}

	removeFromStringList(name: string, item: string)
	{
		// Check if the item is in the list

		this.stringListData.set(name, items);
	}
}

@Injectable()
export class ConfigService {

	// Collections configurations
	public configs: Map<string, ConfigBase> = new Map<string, ConfigBase>()

    constructor(private apiService: ApiService) { }

	public refresh(collection: string)
	{
	}

	public update(collection: string)
	{
	}
}
