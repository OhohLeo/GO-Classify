import { Injectable } from '@angular/core';
import { Response } from '@angular/http';

import { Observable } from 'rxjs/Rx';

import { ApiService, Event } from '../api.service';
import { BufferService } from '../buffer/buffer.service';
import { ReferencesService } from '../references/references.service'

import { BaseElement } from '../base'

@Injectable()
export class WorkflowService {
    
    constructor(private apiService: ApiService,
		private bufferService: BufferService,
		private referencesService: ReferencesService) { }
}
