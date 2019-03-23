import { Component, NgZone, Input, OnInit } from '@angular/core'
import { DataReference, AttributeReference, DataValues } from '../references/data'
import { BaseElement } from '../base'

@Component({
    selector: 'workflow-data',
    template: `
<p class="title">{{data.name}}
  <i class="material-icons"
     (click)="addDataValues()">plus_one</i>
</p>
<div *ngFor="let attributeValues of values; let index = index">
  <workflow-attribute
    *ngFor="let attributeReference of attributes" 
    [element]="element"
    [ref]="attributeReference"
    [index]="index"
    [value]="attributeValues.getAttribute(attributeReference.name)">
  </workflow-attribute>
  <i class="material-icons" (click)="deleteDataValues(index)">delete_forever</i>
</div>`,
    styles: [
`.title {
    background-color: red;
    color: white;
}`],
})

export class DataComponent implements OnInit {

    @Input() data: DataReference
    @Input() element: BaseElement
    
    public attributes: Array<AttributeReference> = new Array<AttributeReference>()
    public values: Array<DataValues> = new Array<DataValues>()
        
    constructor(private zone: NgZone) {}
    
    ngOnInit() {	
	this.zone.run(() => {
	    this.attributes = this.data.getAttributes()
	    this.addDataValues()
	})
    }

    addDataValues() {
	this.values.push(new DataValues(this.data, undefined))
    }

    deleteDataValues(index: number) {
	this.values.splice(index, 1)
    }
}
