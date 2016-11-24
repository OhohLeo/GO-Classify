import { Component, NgZone } from '@angular/core';
import { ImportsService, ImportInterface, Directory } from './imports.service';
import { ClassifyService } from '../classify.service';

@Component({
    selector: 'list',
    templateUrl: './list.component.html'
})

export class ListComponent {

	public importTypes: Array<string> = [];
	public imports: Map<string, ImportInterface[]>;

	constructor(private zone: NgZone,
				private importsService: ImportsService) {}

    ngOnInit() {
		this.update();
    }

	update() {
		this.importsService.getImports()
			.subscribe((imports: Map<string, ImportInterface[]>) => {

				let importTypes: Array<string> = [];

				imports.forEach(function(undefined, importType) {
					importTypes.push(importType)
				})

				this.zone.run(() => {
					this.importTypes = importTypes
					this.imports = imports
				})
			})
	}

	onRefresh(item: ImportInterface) {
		console.log("REFRESH", item)
	}

	onDelete(item: ImportInterface) {
		console.log("DELETE", item)
		this.importsService.deleteImport(item)
			.subscribe(rsp => {
				let imports = this.imports.get(item.getType());
				console.log(rsp)

			})
	}
}
